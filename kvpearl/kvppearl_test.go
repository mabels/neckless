package kvpearl

import (
	"encoding/json"
	"sort"
	"strings"
	"testing"

	"neckless.adviser.com/key"
	"neckless.adviser.com/pearl"
)

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
			},
			&Value{
				Value: "Heute",
			},
			&Value{
				Value: "Gestern",
			},
		}
		sort.Sort(&ValueSorter{
			values: values,
			by:     ValueBy,
		})
		if len(values) != 3 {
			t.Error("unexpected value len")
		}
		if strings.Compare(values[0].Value, "Gestern") != 0 {
			t.Error("unsorted", values)
		}
		if strings.Compare(values[1].Value, "Heute") != 0 {
			t.Error("unsorted", values)
		}
		if strings.Compare(values[2].Value, "Morgen") != 0 {
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
	if strings.Compare(k1.Values[0].Value, "V1") != 0 {
		t.Error("no expected V1:", k1.Values[0].Value)
	}
	if strings.Compare(k1.Values[1].Value, "V2") != 0 {
		t.Error("no expected V2:", k1.Values[1].Value)
	}
	if len(k1.Values[0].Tags) != 3 {
		t.Error("no expected 0:", k1.Key, k1.Values[0].Value, k1.Values[0].Tags)
	}
	// t.Error(len(k1.Values), k1.Values)
	if len(k1.Values[1].Tags) != 4 {
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
		c.Set("K0", "V0")
		c.Set("K1", "V1", "V11", "V12")
		c.Set("K1", "V1", "V11", "V12", "V13")
		c.Set("K1", "V2", "V21", "V22")
		c.Set("K1", "V2", "V23")
		c.Set("K1", "V2", "V24", "V23")
		c.Set("K2", "V2", "T2")
		compareKVPearl(t, c)
	}
	// t.Error("xx")
}

func TestJson(t *testing.T) {
	c := Create( /* "hello", "doof", "hello", "zwei", "doof" */ )
	// if len(c.Tags) != 3 {
	// 	t.Error("no expected", len(c.Tags), c.Tags)
	// }
	c.Set("K0", "V0")
	c.Set("K1", "V1", "V11", "V12")
	c.Set("K1", "V1", "V11", "V12", "V13")
	c.Set("K1", "V2", "V21", "V22")
	c.Set("K1", "V2", "V23")
	c.Set("K1", "V2", "V24", "V23")
	c.Set("K2", "V2", "T2")
	compareKVPearl(t, c)
	jsonStr, err := json.Marshal(c.AsJson())
	if err != nil {
		t.Error("should not ", err)
	}
	// t.Error(c, string(jsonStr))
	kvp, err := FromJson(jsonStr)
	if err != nil {
		t.Error("should not ", err)
	}
	compareKVPearl(t, kvp)

	// func FromJson(jsStr []byte) (*KVPearl, error) {
}

func TestPearl(t *testing.T) {
	c := Create( /*"hello", "doof", "hello", "zwei", "doof"*/ )
	// if len(c.Tags) != 3 {
	// 	t.Error("no expected", len(c.Tags), c.Tags)
	// }
	c.Set("K0", "V0")
	c.Set("K1", "V1", "V11", "V12")
	c.Set("K1", "V1", "V11", "V12", "V13")
	c.Set("K1", "V2", "V21", "V22")
	c.Set("K1", "V2", "V23")
	c.Set("K1", "V2", "V24", "V23")
	c.Set("K2", "V2", "T2")
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
		*Create().Set("AB", "BA").Set("ZA", "AZ"),
		*Create().Set("K1", "V1").Set("K1", "V4").Set("K2", "VV").Set("K3", "Geheim"),
		*Create().Set("K1", "1V1").Set("K1", "V4").Set("K2", "1VV").Set("K3", "SehrGeheim"),
		*Create().Set("K6", "1V1").Set("K7", "V4").Set("K8", "1VV").Set("K1", "SehrGeheim"),
		*Create(),
		*Create().Set("YU", "UY").Set("AA", "AA"),
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
		*Create().Set("K1", "V1").Set("K1", "V4").Set("K2", "VV").Set("K3", "Geheim"),
		*Create().Set("K1", "1V1").Set("K1", "V4").Set("K2", "1VV").Set("K3", "SehrGeheim"),
		*Create().Set("K6", "1V1").Set("K7", "V4").Set("K8", "1VV").Set("K1", "SehrGeheim"),
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
		*Create().Set("K1", "V1", "PROD").Set("K1", "V4").Set("K2", "VV").Set("K3", "Geheim", "PROD"),
		*Create().Set("K1", "1V1", "TEST").Set("K1", "V4", "WURST", "TEST").Set("K2", "1VV").Set("K3", "SehrGeheim"),
		*Create().Set("K6", "1V1").Set("K7", "V4").Set("K8", "1VV", "PROD", "DEV").Set("K1", "SehrGeheim"),
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
	if strings.Compare("K1", u[0].Key) != 0 || strings.Compare("V4", u[0].Value) != 0 {
		t.Error("len == 1", u)
	}
}

func TestGetTags(t *testing.T) {
	o := KVPearls{
		*Create().Set("K1", "V1", "PROD").Set("K1", "V4").Set("K2", "VV").Set("K3", "Geheim", "PROD"),
		*Create().Set("K1", "1V1", "TEST").Set("K1", "V4", "WURST", "TEST").Set("K2", "1VV").Set("K3", "SehrGeheim"),
		*Create().Set("K6", "1V1").Set("K7", "V4").Set("K8", "1VV").Set("K1", "SehrGeheim"),
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
