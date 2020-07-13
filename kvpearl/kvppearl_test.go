package kvpearl

import (
	"bytes"
	"encoding/json"
	"sort"
	"strings"
	"testing"
	"time"

	"neckless.adviser.com/key"
	"neckless.adviser.com/pearl"
)

var TestTime int64 = 4711

func testTime() time.Time {
	TestTime++
	return time.Unix(TestTime, TestTime)
}

func TestUniqStrings(t *testing.T) {
	my := uniqStrings([]string{"hello", "doof", "hello", "zwei", "doof"})
	ref := []string{"doof", "hello", "zwei"}
	if len(my) != len(ref) {
		t.Error("no expected len", len(my), len(ref))
	}
	for i := range ref {
		if strings.Compare(my[i], ref[i]) != 0 {
			t.Error("need ", my[i], ref[i])
		}
	}
}

func TestValueSorter(t *testing.T) {
	for i := 0; i < 100; i++ {
		values := [](*Value){
			&Value{
				Value: "Morgen",
				Order: time.Unix(10000, 0),
			},
			&Value{
				Value: "Heute",
				Order: time.Unix(5000, 0),
			},
			&Value{
				Value: "Gestern",
				Order: time.Unix(2000, 0),
			},
		}
		sort.Sort(&ValueSorter{
			values: values,
		})
		if len(values) != 3 {
			t.Error("unexpected value len")
		}
		if strings.Compare(values[2].Value, "Gestern") != 0 {
			t.Error("unsorted", values)
		}
		if strings.Compare(values[1].Value, "Heute") != 0 {
			t.Error("unsorted", values)
		}
		if strings.Compare(values[0].Value, "Morgen") != 0 {
			t.Error("unsorted", values)
		}
	}
}

func TestCreate(t *testing.T) {
	c := Create( /* "hello", "doof", "hello", "zwei", "doof" */ )
	// ref := []string{"doof", "hello", "zwei"}
	// for i := range ref {
	// 	if strings.Compare(c.Tags[i], ref[i]) != 0 {
	// 		t.Error("need ", c.Tags[i], ref[i])
	// 	}
	// }
	if len(c.Keys) != 0 {
		t.Error("should be empty", c.Keys)
	}
}

func compareKVPearl(t *testing.T, c *KVPearl) {
	if len(c.Keys) != 3 {
		t.Error("no expected", c.Keys)
	}
	// K0 Test
	k0 := c.Keys["K0"]
	if strings.Compare(k0.Key, "K0") != 0 {
		t.Error("no expected", k0.Key)
	}
	if len(k0.Values) != 1 {
		t.Error("no expected", k0.Values)
	}
	if strings.Compare(k0.Values[0].Value, "V0") != 0 {
		t.Error("no expected", k0.Values)
	}
	if len(k0.Values[0].Tags) != 0 {
		t.Error("no expected", k0.Values[0].Tags)
	}
	// K1
	k1 := c.Keys["K1"]
	if strings.Compare(k1.Key, "K1") != 0 {
		t.Error("no expected", k1.Key)
	}
	if len(k1.Values) != 2 {
		t.Error("no expected", k1.Values)
	}
	if strings.Compare(k1.Values[1].Value, "V1") != 0 {
		t.Error("no expected V1:", k1.Values[0].Value)
	}
	if strings.Compare(k1.Values[0].Value, "V2") != 0 {
		t.Error("no expected V2:", k1.Values[1].Value)
	}
	if len(k1.Values[1].Tags) != 3 {
		t.Error("no expected 0:", k1.Key, k1.Values[0].Value, k1.Values[0].Tags)
	}
	// t.Error(len(k1.Values), k1.Values)
	if len(k1.Values[0].Tags) != 4 {
		t.Error("no expected 1:", k1.Key, k1.Values[1].Value, k1.Values[1].Tags)
	}
	// func (kvp *KVPearl) Set(keyVal string, val string, tags ...string) Key {
	k2 := c.Keys["K2"]
	if strings.Compare(k2.Key, "K2") != 0 {
		t.Error("no expected", k2.Key)
	}
	if len(k2.Values) != 1 {
		t.Error("no expected Values", k2.Values)
	}
	if strings.Compare(k2.Values[0].Value, "V2") != 0 {
		t.Error("no expected V2", k2.Values)
	}
	if len(k2.Values[0].Tags) != 1 {
		t.Error("no expected Tags", k2)
	}
}

func TestSet(t *testing.T) {
	for i := 0; i < 100; i++ {
		c := Create( /*"hello", "doof", "hello", "zwei", "doof"*/ )
		// if len(c.Tags) != 3 {
		// 	t.Error("no expected", len(c.Tags), c.Tags)
		// }
		c.Set(testTime(), "K0", "V0")
		c.Set(testTime(), "K1", "V1", "V11", "V12")
		c.Set(testTime(), "K1", "V1", "V11", "V12", "V13")
		c.Set(testTime(), "K1", "V2", "V21", "V22")
		c.Set(testTime(), "K1", "V2", "V23")
		c.Set(testTime(), "K1", "V2", "V24", "V23")
		c.Set(testTime(), "K2", "V2", "T2")
		compareKVPearl(t, c)
	}
	// t.Error("xx")
}

func TestJson(t *testing.T) {
	c := Create( /* "hello", "doof", "hello", "zwei", "doof" */ )
	// if len(c.Tags) != 3 {
	// 	t.Error("no expected", len(c.Tags), c.Tags)
	// }
	c.Set(testTime(), "K0", "V0")
	c.Set(testTime(), "K1", "V1", "V11", "V12")
	c.Set(testTime(), "K1", "V1", "V11", "V12", "V13")
	c.Set(testTime(), "K1", "V2", "V21", "V22")
	c.Set(testTime(), "K1", "V2", "V23")
	c.Set(testTime(), "K1", "V2", "V24", "V23")
	c.Set(testTime(), "K2", "V2", "T2")
	compareKVPearl(t, c)
	var jo *JsonKVPearl
	var prev *JsonKVPearl = nil
	for i := 0; i < 100; i++ {
		jo = c.AsJSON()
		if len(jo.Keys) != 3 {
			t.Error("we need something to sort")
		}
		if prev != nil {
			jp, _ := json.Marshal(jo.Keys)
			pp, _ := json.Marshal(prev.Keys)
			if bytes.Compare(jp, pp) != 0 {
				t.Error("sorting problem")
			}
		}
		prev = jo
	}

	jsonStr, err := json.Marshal(c.AsJSON())
	if err != nil {
		t.Error("should not ", err)
	}

	// t.Error(c, string(jsonStr))
	kvp, err := FromJSON(jsonStr)
	if err != nil {
		t.Error("should not ", err)
	}
	compareKVPearl(t, kvp)

	// func FromJSON(jsStr []byte) (*KVPearl, error) {
}

func TestPearl(t *testing.T) {
	c := Create( /*"hello", "doof", "hello", "zwei", "doof"*/ )
	// if len(c.Tags) != 3 {
	// 	t.Error("no expected", len(c.Tags), c.Tags)
	// }
	c.Set(testTime(), "K0", "V0")
	c.Set(testTime(), "K1", "V1", "V11", "V12")
	c.Set(testTime(), "K1", "V1", "V11", "V12", "V13")
	c.Set(testTime(), "K1", "V2", "V21", "V22")
	c.Set(testTime(), "K1", "V2", "V23")
	c.Set(testTime(), "K1", "V2", "V24", "V23")
	c.Set(testTime(), "K2", "V2", "T2")
	compareKVPearl(t, c)
	rk, _ := key.CreateRandomKey()
	pk := key.MakePrivateKey(rk)
	prl, err := c.ClosePearl(&pearl.PearlOwner{
		Signer: pk,
		Owners: []*key.PublicKey{pk.Public()},
	})
	if err != nil {
		t.Error("unexpeced error:", err)
	}
	kvp, err := OpenPearl([]*key.PrivateKey{pk}, prl)
	if err != nil {
		t.Error("unexpeced error:", err)
	}
	compareKVPearl(t, kvp)
}

func TestLs(t *testing.T) {
	o := KVPearls{
		*Create().Set(testTime(), "AB", "BA").Set(testTime(), "ZA", "AZ"),
		*Create().Set(testTime(), "K1", "V1").Set(testTime(), "K1", "V4").Set(testTime(), "K2", "VV").Set(testTime(), "K3", "Geheim"),
		*Create().Set(testTime(), "K1", "1V1").Set(testTime(), "K1", "V4").Set(testTime(), "K2", "1VV").Set(testTime(), "K3", "SehrGeheim"),
		*Create().Set(testTime(), "K6", "1V1").Set(testTime(), "K7", "V4").Set(testTime(), "K8", "1VV").Set(testTime(), "K1", "SehrGeheim"),
		*Create(),
		*Create().Set(testTime(), "YU", "UY").Set(testTime(), "AA", "AA"),
	}
	u := o.Ls()
	if len(u) != 10 {
		t.Error("not the right len", u)
	}
	if strings.Compare(u[0].Key, "AA") != 0 {
		t.Error("not the right key", u)
	}
	if strings.Compare(u[len(u)-1].Key, "ZA") != 0 {
		t.Error("not the right key", u)
	}
	if strings.Compare(u[2].Key, "K1") != 0 || strings.Compare(u[2].Value, "SehrGeheim") != 0 {
		t.Error("not the right key", u[2])
	}
	if strings.Compare(u[3].Key, "K2") != 0 || strings.Compare(u[3].Value, "1VV") != 0 {
		t.Error("not the right key", u[2])
	}
	if strings.Compare(u[7].Key, "K8") != 0 || strings.Compare(u[7].Value, "1VV") != 0 {
		t.Error("not the right key", u[2])
	}
	// t.Error(u)

}

func TestGet(t *testing.T) {
	o := KVPearls{
		*Create().Set(testTime(), "K1", "V1").Set(testTime(), "K1", "V4").Set(testTime(), "K2", "VV").Set(testTime(), "K3", "Geheim"),
		*Create().Set(testTime(), "K1", "1V1").Set(testTime(), "K1", "V4").Set(testTime(), "K2", "1VV").Set(testTime(), "K3", "SehrGeheim"),
		*Create().Set(testTime(), "K6", "1V1").Set(testTime(), "K7", "V4").Set(testTime(), "K8", "1VV").Set(testTime(), "K1", "SehrGeheim"),
	}
	ret := o.Get("XX")
	if ret != nil {
		t.Error("should be nil")
	}
	ret = o.Get("K1")
	// t.Error(ret)
	if strings.Compare("SehrGeheim", ret.Value) != 0 {
		t.Error("should be", ret)
	}
}

func TestLsTags(t *testing.T) {
	o := KVPearls{
		*Create().Set(testTime(), "K1", "V1", "PROD").Set(testTime(), "K1", "V4").Set(testTime(), "K2", "VV").Set(testTime(), "K3", "Geheim", "PROD"),
		*Create().Set(testTime(), "K1", "1V1", "TEST").Set(testTime(), "K1", "V4", "WURST", "TEST").Set(testTime(), "K2", "1VV").Set(testTime(), "K3", "SehrGeheim"),
		*Create().Set(testTime(), "K6", "1V1").Set(testTime(), "K7", "V4").Set(testTime(), "K8", "1VV", "PROD", "DEV").Set(testTime(), "K1", "SehrGeheim"),
	}
	u := o.Ls("PROD")
	if len(u) != 3 {
		t.Error("unknown length", u)
	}
	if strings.Compare("K1", u[0].Key) != 0 || strings.Compare("V1", u[0].Value) != 0 {
		t.Error("len == 1", u)
	}
	if strings.Compare("K8", u[2].Key) != 0 || strings.Compare("1VV", u[2].Value) != 0 {
		t.Error("len == 1", u)
	}

	u = o.Ls("TEST")
	if len(u) != 1 {
		t.Error("len == 1", u)
	}
	// jo, _ := json.MarshalIndent(u, "", "  ")
	// t.Error(string(jo))
	if strings.Compare("K1", u[0].Key) != 0 || strings.Compare("V4", u[0].Value) != 0 {
		t.Error("len == 1", u[0].Value)
	}
}

func TestGetTags(t *testing.T) {
	o := KVPearls{
		*Create().Set(testTime(), "K1", "V1", "PROD").Set(testTime(), "K1", "V4").Set(testTime(), "K2", "VV").Set(testTime(), "K3", "Geheim", "PROD"),
		*Create().Set(testTime(), "K1", "1V1", "TEST").Set(testTime(), "K1", "V4", "WURST", "TEST").Set(testTime(), "K2", "1VV").Set(testTime(), "K3", "SehrGeheim"),
		*Create().Set(testTime(), "K6", "1V1").Set(testTime(), "K7", "V4").Set(testTime(), "K8", "1VV").Set(testTime(), "K1", "SehrGeheim"),
	}
	ret := o.Get("XX")
	if ret != nil {
		t.Error("should be nil")
	}
	ret = o.Get("K1")
	if strings.Compare("SehrGeheim", ret.Value) != 0 {
		t.Error("should be", ret)
	}
	ret = o.Get("K1", "PROD")
	if strings.Compare("V1", ret.Value) != 0 {
		t.Error("should be", ret)
	}
	ret = o.Get("K1", "TEST")
	if strings.Compare("V4", ret.Value) != 0 {
		t.Error("should be", ret)
	}
	ret = o.Get("K1", "PROD", "TEST")
	if strings.Compare("V4", ret.Value) != 0 {
		t.Error("should be", ret)
	}
}

func TestParse(t *testing.T) {
	kvp := Create()
	_, err := kvp.Parse("")
	if err == nil {
		t.Error("no error not allowed")
	}
	_, err = kvp.Parse("mmm")
	if err == nil {
		t.Error("no error not allowed")
	}
	_, err = kvp.Parse("mmm=")
	if err == nil {
		t.Error("no error not allowed")
	}
	_, err = kvp.Parse("mmm=ooo")
	if err != nil {
		t.Error("no error not allowed", err)
	}
	if strings.Compare(kvp.Keys["mmm"].Key, "mmm") != 0 {
		t.Error("no error not allowed")
	}
	if len(kvp.Keys["mmm"].Values) != 1 {
		t.Error("no error not allowed")
	}
	if len(kvp.Keys["mmm"].Values[0].Tags) != 0 {
		t.Error("no error not allowed", kvp.Keys["mmm"].Values[0], len(kvp.Keys["mmm"].Values[0].Tags))
	}
	if strings.Compare(kvp.Keys["mmm"].Values[0].Value, "ooo") != 0 {
		t.Error("no error not allowed")
	}
	_, err = kvp.Parse("mmm=ooo[")
	if err == nil {
		t.Error("no error not allowed", err)
	}
	_, err = kvp.Parse("mmm=ooo[vv")
	if err == nil {
		t.Error("no error not allowed", err)
	}

	_, err = kvp.Parse("mmm=ooo[]")
	if err != nil {
		t.Error("no error not allowed", err)
	}
	if strings.Compare(kvp.Keys["mmm"].Values[0].Value, "ooo") != 0 {
		t.Error("no error not allowed")
	}
	if len(kvp.Keys["mmm"].Values[0].Tags) != 0 {
		t.Error("no error not allowed")
	}

	_, err = kvp.Parse("mmm=ooo[AA]")
	if err != nil {
		t.Error("no error not allowed", err)
	}
	if strings.Compare(kvp.Keys["mmm"].Values[0].Value, "ooo") != 0 {
		t.Error("no error not allowed")
	}
	if strings.Compare(kvp.Keys["mmm"].Values[0].Tags[0], "AA") != 0 {
		t.Error("no error not allowed")
	}
	_, err = kvp.Parse("mmm=ooo[AA,BB]")
	if err != nil {
		t.Error("no error not allowed", err)
	}
	if strings.Compare(kvp.Keys["mmm"].Values[0].Value, "ooo") != 0 {
		t.Error("no error not allowed")
	}
	if strings.Compare(kvp.Keys["mmm"].Values[0].Tags[0], "AA") != 0 {
		t.Error("no error not allowed")
	}
	if strings.Compare(kvp.Keys["mmm"].Values[0].Tags[1], "BB") != 0 {
		t.Error("no error not allowed")
	}
	_, err = kvp.Parse("mmm=ooo[AA,,BB,]")
	if err != nil {
		t.Error("no error not allowed", err)
	}
	if strings.Compare(kvp.Keys["mmm"].Values[0].Value, "ooo") != 0 {
		t.Error("no error not allowed")
	}
	if strings.Compare(kvp.Keys["mmm"].Values[0].Tags[0], "AA") != 0 {
		t.Error("no error not allowed")
	}
	if strings.Compare(kvp.Keys["mmm"].Values[0].Tags[1], "BB") != 0 {
		t.Error("no error not allowed")
	}
	if len(kvp.Keys["mmm"].Values[0].Tags) != 2 {
		t.Error("no error not allowed", kvp.Keys["mmm"].Values[0].Tags)
	}
}

func TestMerge(t *testing.T) {
	kvps := []*KVPearl{
		Create().Set(testTime(), "K1", "V1", "PROD").Set(testTime(), "K1", "V4").Set(testTime(), "K2", "VV").Set(testTime(), "K3", "Geheim", "PROD"),
		Create().Set(testTime(), "K1", "1V1", "TEST").Set(testTime(), "K1", "V4", "WURST", "TEST").Set(testTime(), "K2", "1VV").Set(testTime(), "K3", "SehrGeheim"),
		Create().Set(testTime(), "K6", "1V1").Set(testTime(), "K7", "V4").Set(testTime(), "K8", "1VV").Set(testTime(), "K1", "SehrGeheim"),
	}
	kvp := Merge(kvps, []string{}, []string{}).AsJSON()
	// js, _ := json.MarshalIndent(kvp, "", "  ")
	if strings.Compare(kvp.Keys[0].Values[0].Value, "SehrGeheim") != 0 {
		t.Error("failed order")
	}
	if strings.Compare(kvp.Keys[1].Values[0].Value, "1VV") != 0 {
		t.Error("failed order")
	}
	// t.Error(string(js))
	// kvp.Keys
	// kvp.FingerPrint
	// kvp.Created

}
