package kvpearl

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"strings"
	"testing"
	"time"

	"neckless.adviser.com/key"
	"neckless.adviser.com/pearl"
)

var TestTime int64 = time.Now().UnixNano()

func testTime() time.Time {
	TestTime++
	n10e9 := int64(math.Pow10(9))
	return time.Unix(TestTime/n10e9, TestTime%n10e9)
}

// func TestUniqStrings(t *testing.T) {
// 	my := uniqStrings([]string{"hello", "doof", "hello", "zwei", "doof"})
// 	ref := []string{"doof", "hello", "zwei"}
// 	if len(my) != len(ref) {
// 		t.Error("no expected len", len(my), len(ref))
// 	}
// 	for i := range ref {
// 		if strings.Compare(my[i], ref[i]) != 0 {
// 			t.Error("need ", my[i], ref[i])
// 		}
// 	}
// }

// func TestValueSorter(t *testing.T) {
// 	for i := 0; i < 100; i++ {
// 		values := [](*Value){
// 			&Value{
// 				Value: "Morgen",
// 				Order: time.Unix(10000, 0),
// 			},
// 			&Value{
// 				Value: "Heute",
// 				Order: time.Unix(5000, 0),
// 			},
// 			&Value{
// 				Value: "Gestern",
// 				Order: time.Unix(2000, 0),
// 			},
// 		}
// 		sort.Sort(&ValueSorter{
// 			values: values,
// 		})
// 		if len(values) != 3 {
// 			t.Error("unexpected value len")
// 		}
// 		if strings.Compare(values[2].Value, "Gestern") != 0 {
// 			t.Error("unsorted", values)
// 		}
// 		if strings.Compare(values[1].Value, "Heute") != 0 {
// 			t.Error("unsorted", values)
// 		}
// 		if strings.Compare(values[0].Value, "Morgen") != 0 {
// 			t.Error("unsorted", values)
// 		}
// 	}
// }

func TestCreate(t *testing.T) {
	c := CreateKVPearls().Add( /* "hello", "doof", "hello", "zwei", "doof" */ )
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
	if k0.Values.len() != 1 {
		t.Error("no expected", k0.Values)
	}
	if strings.Compare(k0.Values.get("V0").Value, "V0") != 0 {
		t.Error("no expected", k0.Values)
	}
	if len(k0.Values.get("V0").Tags) != 0 {
		t.Error("no expected", k0.Values.get("V0").Tags)
	}
	// K1
	k1 := c.Keys["K1"]
	if strings.Compare(k1.Key, "K1") != 0 {
		t.Error("no expected", k1.Key)
	}
	if k1.Values.len() != 2 {
		t.Error("no expected", k1.Values)
	}
	if strings.Compare(k1.Values.get("V1").Value, "V1") != 0 {
		t.Error("no expected V1:", k1.Values.get("V1").Value)
	}
	if strings.Compare(k1.Values.get("V2").Value, "V2") != 0 {
		t.Error("no expected V2:", k1.Values.get("V2").Value)
	}
	if len(k1.Values.get("V2").Tags) != 4 {
		t.Error("no expected 0:", k1.Key, k1.Values.get("V2").Value, k1.Values.get("V2").Tags)
	}
	// t.Error(k1.Values.len(), k1.Values)
	if len(k1.Values.get("V1").Tags) != 3 {
		t.Error("no expected 1:", k1.Key, k1.Values.get("V1").Value, k1.Values.get("V2").Tags)
	}
	// func (kvp *KVPearl) Set(keyVal string, val string, tags ...string) Key {
	k2 := c.Keys["K2"]
	if strings.Compare(k2.Key, "K2") != 0 {
		t.Error("no expected", k2.Key)
	}
	if k2.Values.len() != 1 {
		t.Error("no expected Values", k2.Values)
	}
	if strings.Compare(k2.Values.get("V2").Value, "V2") != 0 {
		t.Error("no expected V2", k2.Values)
	}
	if len(k2.Values.get("V2").Tags) != 1 {
		t.Error("no expected Tags", k2)
	}
}

func TestSet(t *testing.T) {
	o := CreateKVPearls()
	for i := 0; i < 100; i++ {
		c := o.Add( /*"hello", "doof", "hello", "zwei", "doof"*/ )
		// if len(c.Tags) != 3 {
		// 	t.Error("no expected", len(c.Tags), c.Tags)
		// }
		c.Set(SetArg{Key: "K0", Val: "V0"})
		c.Set(SetArg{Key: "K1", Val: "V1", Tags: []string{"V11", "V12"}})
		c.Set(SetArg{Key: "K1", Val: "V1", Tags: []string{"V11", "V12", "V13"}})
		c.Set(SetArg{Key: "K1", Val: "V2", Tags: []string{"V21", "V22"}})
		c.Set(SetArg{Key: "K1", Val: "V2", Tags: []string{"V23"}})
		c.Set(SetArg{Key: "K1", Val: "V2", Tags: []string{"V24", "V23"}})
		c.Set(SetArg{Key: "K2", Val: "V2", Tags: []string{"T2"}})
		compareKVPearl(t, c)
	}
	// t.Error("xx")
}

func TestJson(t *testing.T) {
	c := CreateKVPearls().Add( /* "hello", "doof", "hello", "zwei", "doof" */ )
	// if len(c.Tags) != 3 {
	// 	t.Error("no expected", len(c.Tags), c.Tags)
	// }
	c.Set(SetArg{Key: "K0", Val: "V0"})
	c.Set(SetArg{Key: "K1", Val: "V1", Tags: []string{"V11", "V12"}})
	c.Set(SetArg{Key: "K1", Val: "V1", Tags: []string{"V11", "V12", "V13"}})
	c.Set(SetArg{Key: "K1", Val: "V2", Tags: []string{"V21", "V22"}})
	c.Set(SetArg{Key: "K1", Val: "V2", Tags: []string{"V23"}})
	c.Set(SetArg{Key: "K1", Val: "V2", Tags: []string{"V24", "V23"}})
	c.Set(SetArg{Key: "K2", Val: "V2", Tags: []string{"T2"}})
	compareKVPearl(t, c)
	var jo *JSONKVPearl
	var prev *JSONKVPearl = nil
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
	c := CreateKVPearls().Add( /*"hello", "doof", "hello", "zwei", "doof"*/ )
	// if len(c.Tags) != 3 {
	// 	t.Error("no expected", len(c.Tags), c.Tags)
	// }
	c.Set(SetArg{Key: "K0", Val: "V0"})
	c.Set(SetArg{Key: "K1", Val: "V1", Tags: []string{"V11", "V12"}})
	c.Set(SetArg{Key: "K1", Val: "V1", Tags: []string{"V11", "V12", "V13"}})
	c.Set(SetArg{Key: "K1", Val: "V2", Tags: []string{"V21", "V22"}})
	c.Set(SetArg{Key: "K1", Val: "V2", Tags: []string{"V23"}})
	c.Set(SetArg{Key: "K1", Val: "V2", Tags: []string{"V24", "V23"}})
	c.Set(SetArg{Key: "K2", Val: "V2", Tags: []string{"T2"}})
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

// func TestLs(t *testing.T) {
// 	o := CreateKVPearls()
// 	o.Add().Set(SetArg{Key: "AB", Val: "BA"}).Set(SetArg{Key: "ZA", Val: "AZ"})
// 	o.Add().Set(SetArg{Key: "K1", Val: "V1"}).Set(SetArg{Key: "K1", Val: "V4"}).Set(SetArg{Key: "K2", Val: "VV"}).Set(SetArg{Key: "K3", Val: "Geheim"})
// 	o.Add().Set(SetArg{Key: "K1", Val: "1V1"}).Set(SetArg{Key: "K1", Val: "V4"}).Set(SetArg{Key: "K2", Val: "1VV"}).Set(SetArg{Key: "K3", Val: "SehrGeheim"})
// 	o.Add().Set(SetArg{Key: "K6", Val: "1V1"}).Set(SetArg{Key: "K7", Val: "V4"}).Set(SetArg{Key: "K8", Val: "1VV"}).Set(SetArg{Key: "K1", Val: "SehrGeheim"})
// 	o.Add()
// 	o.Add().Set(SetArg{Key: "YU", Val: "UY"}).Set(SetArg{Key: "AA", Val: "AA"})
// 	u := o.Ls()
// 	if len(u) != 10 {
// 		t.Error("not the right len", u)
// 	}
// 	if strings.Compare(u[0].Key, "AA") != 0 {
// 		t.Error("not the right key", u)
// 	}
// 	if strings.Compare(u[len(u)-1].Key, "ZA") != 0 {
// 		t.Error("not the right key", u)
// 	}
// 	if strings.Compare(u[2].Key, "K1") != 0 || strings.Compare(u[2].Value, "SehrGeheim") != 0 {
// 		t.Error("not the right key", u[2])
// 	}
// 	if strings.Compare(u[3].Key, "K2") != 0 || strings.Compare(u[3].Value, "1VV") != 0 {
// 		t.Error("not the right key", u[2])
// 	}
// 	if strings.Compare(u[7].Key, "K8") != 0 || strings.Compare(u[7].Value, "1VV") != 0 {
// 		t.Error("not the right key", u[2])
// 	}
// 	// t.Error(u)

// }

// func TestGet(t *testing.T) {
// 	o := CreateKVPearls()
// 	o.Add().Set(SetArg{Key: "K1", Val: "V1"}).Set(SetArg{Key: "K1", Val: "V4"}).Set(SetArg{Key: "K2", Val: "VV"}).Set(SetArg{Key: "K3", Val: "Geheim"})
// 	o.Add().Set(SetArg{Key: "K1", Val: "1V1"}).Set(SetArg{Key: "K1", Val: "V4"}).Set(SetArg{Key: "K2", Val: "1VV"}).Set(SetArg{Key: "K3", Val: "SehrGeheim"})
// 	o.Add().Set(SetArg{Key: "K6", Val: "1V1"}).Set(SetArg{Key: "K7", Val: "V4"}).Set(SetArg{Key: "K8", Val: "1VV"}).Set(SetArg{Key: "K1", Val: "SehrGeheim"})
// 	ret := o.Get("XX")
// 	if ret != nil {
// 		t.Error("should be nil")
// 	}
// 	ret = o.Get("K1")
// 	// t.Error(ret)
// 	if strings.Compare("SehrGeheim", ret.Value) != 0 {
// 		t.Error("should be", ret)
// 	}
// }

// func TestLsTags(t *testing.T) {
// 	o := CreateKVPearls()
// 	o.Add().Set(SetArg{Key: "K1", Val: "V1", Tags: []string{"PROD"}}).Set(SetArg{Key: "K1", Val: "V4"}).Set(SetArg{Key: "K2", Val: "VV"}).Set(SetArg{Key: "K3", Val: "Geheim", Tags: []string{"PROD"}})
// 	o.Add().Set(SetArg{Key: "K1", Val: "1V1", Tags: []string{"TEST"}}).Set(SetArg{Key: "K1", Val: "V4", Tags: []string{"WURST", "TEST"}}).Set(SetArg{Key: "K2", Val: "1VV"}).Set(SetArg{Key: "K3", Val: "SehrGeheim"})
// 	o.Add().Set(SetArg{Key: "K6", Val: "1V1"}).Set(SetArg{Key: "K7", Val: "V4"}).Set(SetArg{Key: "K8", Val: "1VV", Tags: []string{"PROD", "DEV"}}).Set(SetArg{Key: "K1", Val: "SehrGeheim"})
// 	u := o.Ls("PROD")
// 	if len(u) != 3 {
// 		t.Error("unknown length", u)
// 	}
// 	if strings.Compare("K1", u[0].Key) != 0 || strings.Compare("V1", u[0].Value) != 0 {
// 		t.Error("len == 1", u)
// 	}
// 	if strings.Compare("K8", u[2].Key) != 0 || strings.Compare("1VV", u[2].Value) != 0 {
// 		t.Error("len == 1", u)
// 	}

// 	u = o.Ls("TEST")
// 	if len(u) != 1 {
// 		t.Error("len == 1", u)
// 	}
// 	// jo, _ := json.MarshalIndent(u, "", "  ")
// 	// t.Error(string(jo))
// 	if strings.Compare("K1", u[0].Key) != 0 || strings.Compare("1V1", u[0].Value) != 0 {
// 		t.Error("len == 1", u[0].Value)
// 	}
// }

// func TestGetTags(t *testing.T) {
// 	o := CreateKVPearls()
// 	o.Add().Set(SetArg{Key: "K1", Val: "V1", Tags: []string{"PROD"}}).Set(SetArg{Key: "K1", Val: "V4"}).Set(SetArg{Key: "K2", Val: "VV"}).Set(SetArg{Key: "K3", Val: "Geheim", Tags: []string{"PROD"}})
// 	o.Add().Set(SetArg{Key: "K1", Val: "1V1", Tags: []string{"TEST"}}).Set(SetArg{Key: "K1", Val: "V4", Tags: []string{"WURST", "TEST"}}).Set(SetArg{Key: "K2", Val: "1VV"}).Set(SetArg{Key: "K3", Val: "SehrGeheim"})
// 	o.Add().Set(SetArg{Key: "K6", Val: "1V1"}).Set(SetArg{Key: "K7", Val: "V4"}).Set(SetArg{Key: "K8", Val: "1VV"}).Set(SetArg{Key: "K1", Val: "SehrGeheim"})
// 	ret := o.Get("XX")
// 	if ret != nil {
// 		t.Error("should be nil")
// 	}
// 	ret = o.Get("K1")
// 	if strings.Compare("SehrGeheim", ret.Value) != 0 {
// 		t.Error("should be", ret)
// 	}
// 	ret = o.Get("K1", "PROD")
// 	if strings.Compare("V1", ret.Value) != 0 {
// 		t.Error("should be", ret)
// 	}
// 	ret = o.Get("K1", "TEST")
// 	if strings.Compare("V4", ret.Value) != 0 {
// 		t.Error("should be", ret)
// 	}
// 	ret = o.Get("K1", "PROD", "TEST")
// 	if strings.Compare("V4", ret.Value) != 0 {
// 		t.Error("should be", ret)
// 	}
// }

func TestParse(t *testing.T) {
	kvps := CreateKVPearls()
	for i := 0; i < 100; i++ {
		_, err := Parse("")
		if err == nil {
			t.Error("no error not allowed")
		}
		_, err = Parse("mmm")
		if err == nil {
			t.Error("no error not allowed")
		}
		kvp := kvps.Add()
		p, err := Parse("mmm=")
		if err != nil {
			t.Error("no error not allowed")
		}
		sa, err := p.ToSetArgs()
		if err != nil {
			t.Error("no error not allowed")
		}
		kvp.Set(*sa)

		p, err = Parse("mmm=ooo")
		if err != nil {
			t.Error("no error not allowed", err)
		}
		sa, err = p.ToSetArgs()
		if err != nil {
			t.Error("no error not allowed")
		}
		kvp.Set(*sa)
		if strings.Compare(kvp.Keys.get("mmm").Key, "mmm") != 0 {
			t.Error("no error not allowed")
		}
		// ooo and empty
		if kvp.Keys.get("mmm").Values.len() != 2 {
			t.Error("no error not allowed")
		}
		if len(kvp.Keys.get("mmm").Values.get("ooo").Tags) != 0 {
			t.Error("no error not allowed", kvp.Keys.get("mmm").Values.get("ooo"), len(kvp.Keys.get("mmm").Values.get("ooo").Tags))
		}
		if strings.Compare(kvp.Keys.get("mmm").Values.get("ooo").Value, "ooo") != 0 {
			t.Error("no error not allowed")
		}
		p, err = Parse("mmm=ooo[")
		if err != nil {
			t.Error("no error not allowed", err)
		}
		sa, err = p.ToSetArgs()
		if err != nil {
			t.Error("no error not allowed")
		}
		kvp.Set(*sa)
		if strings.Compare(kvp.Keys.get("mmm").Values.get("ooo[").Value, "ooo[") != 0 {
			t.Error("no error not allowed", err, kvp.Keys.get("mmm").Values.get("ooo[").Value)
		}
		p, err = Parse("mmm=ooo[vv")
		if err != nil {
			t.Error("no error not allowed", err)
		}
		sa, err = p.ToSetArgs()
		if err != nil {
			t.Error("no error not allowed")
		}
		kvp.Set(*sa)
		if strings.Compare(kvp.Keys.get("mmm").Values.get("ooo[vv").Value, "ooo[vv") != 0 {
			t.Error("no error not allowed", err, kvp.Keys.get("mmm").Values.get("ooo[vv").Value)
		}

		p, err = Parse("mmm=ooo]vv")
		if err != nil {
			t.Error("no error not allowed", err)
		}
		sa, err = p.ToSetArgs()
		if err != nil {
			t.Error("no error not allowed")
		}
		kvp.Set(*sa)
		if strings.Compare(kvp.Keys.get("mmm").Values.get("ooo]vv").Value, "ooo]vv") != 0 {
			t.Error("no error not allowed", err, kvp.Keys.get("mmm").Values)
		}

		p, err = Parse("mmm=yyy[]")
		if err != nil {
			t.Error("no error not allowed", err)
		}
		sa, err = p.ToSetArgs()
		if err != nil {
			t.Error("no error not allowed")
		}
		kvp.Set(*sa)

		if strings.Compare(kvp.Keys.get("mmm").Values.get("yyy").Value, "yyy") != 0 {
			t.Error("no error not allowed")
		}
		if len(kvp.Keys.get("mmm").Values.get("yyy").Tags) != 0 {
			t.Error("no error not allowed")
		}

		p, err = Parse("mmm=uuu[AA]")
		if err != nil {
			t.Error("no error not allowed", err)
		}
		sa, err = p.ToSetArgs()
		if err != nil {
			t.Error("no error not allowed")
		}
		kvp.Set(*sa)
		if strings.Compare(kvp.Keys.get("mmm").Values.get("uuu").Value, "uuu") != 0 {
			t.Error("no error not allowed")
		}
		if strings.Compare(kvp.Keys.get("mmm").Values.get("uuu").Tags.sorted()[0], "AA") != 0 {
			t.Error("no error not allowed")
		}
		p, err = Parse("mmm=rrr[AA,BB]")
		if err != nil {
			t.Error("no error not allowed", err)
		}
		sa, err = p.ToSetArgs()
		if err != nil {
			t.Error("no error not allowed")
		}
		kvp.Set(*sa)
		if strings.Compare(kvp.Keys.get("mmm").Values.get("rrr").Value, "rrr") != 0 {
			t.Error("no error not allowed")
		}
		if strings.Compare(kvp.Keys.get("mmm").Values.get("rrr").Tags.sorted()[0], "AA") != 0 {
			t.Error("no error not allowed")
		}
		if strings.Compare(kvp.Keys.get("mmm").Values.get("rrr").Tags.sorted()[1], "BB") != 0 {
			t.Error("no error not allowed")
		}
		p, err = Parse("mmm=sss[AA,,BB,]")
		if err != nil {
			t.Error("no error not allowed", err)
		}
		sa, err = p.ToSetArgs()
		if err != nil {
			t.Error("no error not allowed")
		}
		kvp.Set(*sa)
		if strings.Compare(kvp.Keys.get("mmm").Values.get("sss").Value, "sss") != 0 {
			t.Error("no error not allowed")
		}
		if vals := kvp.Keys.get("mmm").Values.get("sss").Tags.sorted(); strings.Compare(vals[0], "") != 0 {
			t.Errorf("no error not allowed:%s:%d:%s", vals[0], len(vals), vals)
		}
		if vals := kvp.Keys.get("mmm").Values.get("sss").Tags.sorted(); strings.Compare(vals[1], "AA") != 0 {
			t.Errorf("no error not allowed:%s:%d:%s", vals[0], len(vals), vals)
		}
		if strings.Compare(kvp.Keys.get("mmm").Values.get("sss").Tags.sorted()[2], "BB") != 0 {
			t.Error("no error not allowed")
		}
		if len(kvp.Keys.get("mmm").Values.get("sss").Tags) != 3 {
			t.Error("no error not allowed", kvp.Keys.get("mmm").Values.get("sss").Tags)
		}
		p, err = Parse("mmm=zzz,AA,,BB,")
		if err != nil {
			t.Error("no error not allowed", err)
		}
		sa, err = p.ToSetArgs()
		if err != nil {
			t.Error("no error not allowed")
		}
		kvp.Set(*sa)
		if strings.Compare(kvp.Keys.get("mmm").Values.get("zzz").Value, "zzz") != 0 {
			t.Error("no error not allowed")
		}
		if strings.Compare(kvp.Keys.get("mmm").Values.get("zzz").Tags.sorted()[0], "") != 0 {
			t.Error("no error not allowed")
		}
		if strings.Compare(kvp.Keys.get("mmm").Values.get("zzz").Tags.sorted()[1], "AA") != 0 {
			t.Error("no error not allowed")
		}
		if strings.Compare(kvp.Keys.get("mmm").Values.get("zzz").Tags.sorted()[2], "BB") != 0 {
			t.Error("no error not allowed")
		}
		if len(kvp.Keys.get("mmm").Values.get("zzz").Tags) != 3 {
			t.Error("no error not allowed", kvp.Keys.get("mmm").Values.get("zzz").Tags)
		}
	}
}

func TestMerge(t *testing.T) {
	for i := 0; i < 100; i++ {
		kvps := CreateKVPearls()
		kvps.Add().Set(SetArg{Key: "K1", Val: "V1", Tags: []string{"PROD"}}).Set(SetArg{Key: "K1", Val: "V4"}).Set(SetArg{Key: "K2", Val: "VV"}).Set(SetArg{Key: "K3", Val: "Geheim", Tags: []string{"PROD"}})
		kvps.Add().Set(SetArg{Key: "K1", Val: "1V1", Tags: []string{"TEST"}}).Set(SetArg{Key: "K1", Val: "V4", Tags: []string{"WURST", "TEST"}}).Set(SetArg{Key: "K2", Val: "1VV"}).Set(SetArg{Key: "K3", Val: "SehrGeheim"})
		kvps.Add().Set(SetArg{Key: "K6", Val: "1V1"}).Set(SetArg{Key: "K7", Val: "V4"}).Set(SetArg{Key: "K8", Val: "1VV"}).Set(SetArg{Key: "K1", Val: "SehrGeheim"})
		kvp := kvps.Match(MapByToResolve{})[""]
		if len(kvp) != 6 {
			t.Error(len(kvp))
		}
		if kvp[0].Key != "K1" {
			t.Errorf("should be K1:%s", kvp[0].Key)
		}
		if kvp[1].Key != "K2" {
			t.Errorf("should be K2:%s", kvp[1].Key)
		}
		if kvp[5].Key != "K8" {
			t.Errorf("should be K8:%s", kvp[5].Key)
		}
		// js, _ := json.MarshalIndent(kvp.ToJSON(), "", "  ")
		// t.Error(string(js))
		if strings.Compare(kvp[0].Vals.get("SehrGeheim").Value, "SehrGeheim") != 0 {
			t.Error("failed order")
		}
		if strings.Compare(kvp[1].Vals.get("1VV").Value, "1VV") != 0 {
			t.Error("failed order")
		}
		// t.Error(string(js))
		// kvp.Keys
		// kvp.FingerPrint
		// kvp.Created
	}

}

func TestResolv(t *testing.T) {
	kvp := CreateKVPearls().Add()
	p, err := Parse("mmm@ooo")
	if err != nil {
		t.Error(err)
	}
	if *p.Key != "mmm" {
		t.Error("should be mmm")
	}
	if *&p.ToResolve.Param != "ooo" {
		t.Error("should be ooo")
	}
	if p.KeyRegex.String() != "^mmm$" {
		t.Error(fmt.Sprintf("should be %s", p.KeyRegex.String()))
	}
	if p.Val != nil {
		t.Error("should be ooo")
	}
	if len(p.Tags) != 0 {
		t.Error("should be ooo")
	}

	p, err = Parse("mmm@aaa")
	if err != nil {
		t.Error(err)
	}
	p, err = p.Resolv(func(string, FuncsAndParam) (*string, error) { m := "ooo"; return &m, nil })
	if err != nil {
		t.Error(err)
	}
	sa, err := p.ToSetArgs()
	if err != nil {
		t.Error(err)
	}
	kvp.Set(*sa)
	if kvp.Keys.get("mmm").Values.get("ooo").Value != "ooo" {
		t.Error("should be ooo")
	}
	if *&kvp.Keys.get("mmm").Values.get("ooo").Unresolved.Param != "aaa" {
		t.Error("should be aaa")
	}
}

func TestResolvWithCommaTags(t *testing.T) {
	kvp := CreateKVPearls().Add()
	p, err := Parse("mmm@ooo,T1,T2")
	if err != nil {
		t.Error(err)
	}
	p, err = p.Resolv(func(string, FuncsAndParam) (*string, error) { m := "ooo"; return &m, nil })
	if err != nil {
		t.Error(err)
	}
	sa, err := p.ToSetArgs()
	if err != nil {
		t.Error(err)
	}
	kvp.Set(*sa)
	if kvp.Keys.get("mmm").Values.get("ooo").Value != "ooo" {
		t.Error("should be ooo")
	}
	if len(kvp.Keys.get("mmm").Values.get("ooo").Tags) != 2 {
		t.Error("should tags len 2")
	}
	if kvp.Keys.get("mmm").Values.get("ooo").Tags.sorted()[0] != "T1" {
		t.Error("should tags T1")
	}
	if kvp.Keys.get("mmm").Values.get("ooo").Tags.sorted()[1] != "T2" {
		t.Error("should tags T2")
	}
	p, err = Parse("mmm@aaa,T3,T4")
	if err != nil {
		t.Error(err)
	}
	p, err = p.Resolv(func(string, FuncsAndParam) (*string, error) { m := "ooo"; return &m, nil })
	if err != nil {
		t.Error(err)
	}
	sa, err = p.ToSetArgs()
	if err != nil {
		t.Error(err)
	}
	kvp.Set(*sa)
	if kvp.Keys.get("mmm").Values.get("ooo").Value != "ooo" {
		t.Error("should be ooo")
	}
	if *&kvp.Keys.get("mmm").Values.get("ooo").Unresolved.Param != "aaa" {
		t.Error("should be aaa")
	}
	if len(kvp.Keys.get("mmm").Values.get("ooo").Tags) != 4 {
		t.Error("should tags len 4")
	}
	if kvp.Keys.get("mmm").Values.get("ooo").Tags.sorted()[0] != "T1" {
		t.Error("should tags T1")
	}
	if kvp.Keys.get("mmm").Values.get("ooo").Tags.sorted()[1] != "T2" {
		t.Error("should tags T2")
	}
	if kvp.Keys.get("mmm").Values.get("ooo").Tags.sorted()[2] != "T3" {
		t.Error("should tags T3")
	}
	if kvp.Keys.get("mmm").Values.get("ooo").Tags.sorted()[3] != "T4" {
		t.Error("should tags T4")
	}
	p, err = Parse("mmm@aaa,T3,T4")
	if err != nil {
		t.Error(err)
	}
	p, err = p.Resolv(func(string, FuncsAndParam) (*string, error) { m := "rrr"; return &m, nil })
	if err != nil {
		t.Error(err)
	}
	sa, err = p.ToSetArgs()
	if err != nil {
		t.Error(err)
	}
	kvp.Set(*sa)
	if kvp.Keys.get("mmm").Values.get("rrr").Value != "rrr" {
		t.Error("should be ooo")
	}
	if *&kvp.Keys.get("mmm").Values.get("rrr").Unresolved.Param != "aaa" {
		t.Error("should be aaa")
	}
	if len(kvp.Keys.get("mmm").Values.get("rrr").Tags) != 2 {
		t.Error("should tags len 2")
	}
	if kvp.Keys.get("mmm").Values.get("rrr").Tags.sorted()[0] != "T3" {
		t.Error("should tags T1")
	}
	if kvp.Keys.get("mmm").Values.get("rrr").Tags.sorted()[1] != "T4" {
		t.Error("should tags T2")
	}
}
func testSA(t *testing.T, pos string, err error, sa *KVParsed, vals ...string) {
	// jssa, _ := json.Marshal(sa)
	// t.Error(string(jssa))
	if err != nil {
		t.Error(pos, "not expected here", err)
	}
	if sa.Key == nil {
		t.Error(pos, "should be mmm")
	}
	if *sa.Key != "mmm" {
		t.Error(pos, "should be mmm")
	}
	if sa.KeyRegex.String() != "^mmm$" {
		t.Error(pos, "should be ^mmm$")
	}
	if sa.ToResolve == nil {
		t.Error(pos, "should be unresolved nil")
	}
	if len(vals) == 0 {
		if sa.Val != nil {
			t.Error(pos, "should be val")
		}
	} else {
		if *sa.Val != vals[0] {
			t.Error(pos, "should be val")
		}
	}
	// if *sa.Val != "" {
	// t.Error("should be \"\"")
	// }
	if len(sa.Tags) != 0 {
		t.Error(pos, "should be 0:", sa.Tags)
	}
}

func TestOrderPreserve(t *testing.T) {
	for i := 0; i < 100; i++ {
		inKvp := CreateKVPearls().Add()
		inKvp.Set(SetArg{Key: "M1", Val: "first", Tags: []string{"m1first-1"}})
		inKvp.Set(SetArg{Key: "M1", Val: "second", Tags: []string{"m1second"}})
		inKvp.Set(SetArg{Key: "M1", Val: "first", Tags: []string{"m1first-2"}})
		u, _ := json.Marshal(inKvp.AsJSON())
		outKvp, _ := FromJSON(u)
		val := outKvp.Keys.get("M1").Values
		if val.len() != 2 {
			t.Error("need to be 2")
		}
		if val.get("first").order <= val.get("second").order {
			t.Error("order should be right")
		}
		if val.get("second").Tags.sorted()[0] != "m1second" {
			t.Error("second m1second")
		}
		if !(val.get("first").Tags.sorted()[0] == "m1first-1" &&
			val.get("first").Tags.sorted()[1] == "m1first-2") {
			t.Errorf("first m1first-1 m1first-2:%d", i)
		}
	}
}

func TestEmptyParse(t *testing.T) {
	sa, err := Parse("mmm@[]")
	testSA(t, "1", err, sa)
	sa, err = Parse("mmm@[]")
	sa, err = sa.Resolv(func(string, FuncsAndParam) (*string, error) { m := "rrr"; return &m, nil })
	testSA(t, "2", err, sa, "rrr")
	sa, err = Parse("mmm@,")
	testSA(t, "3", err, sa)
	sa, err = Parse("mmm@,")
	sa, err = sa.Resolv(func(string, FuncsAndParam) (*string, error) { m := "rrr"; return &m, nil })
	testSA(t, "4", err, sa, "rrr")
	sa, err = Parse("mmm@[,]")
	if err != nil {
		t.Error(err)
	}
	if !(len(sa.Tags.toArray()) == 1 && len(sa.Tags.toArray()[0]) == 0) {
		t.Error("There should be an empty")
	}
	sa, err = Parse("mmm@,, ,")
	if err != nil {
		t.Error(err)
	}
	if !(len(sa.Tags.toArray()) == 1 && len(sa.Tags.toArray()[0]) == 0) {
		t.Error("There should be an empty")
	}
}

func TestResolvWithBracketsTags(t *testing.T) {
	kvp := CreateKVPearls().Add()
	p, err := Parse("mmm@ZZZ[T1,T2]")
	if err != nil {
		t.Error(err)
	}
	p, err = p.Resolv(func(key string, fparam FuncsAndParam) (*string, error) {
		if key != "mmm" {
			t.Error("should be mmm")
		}
		if fparam.Param != "ZZZ" {
			t.Error("should be ZZZ")
		}
		m := "ooo"
		return &m, nil
	})
	if err != nil {
		t.Error(err)
	}
	sa, err := p.ToSetArgs()
	if err != nil {
		t.Error(err)
	}
	kvp.Set(*sa)
	if kvp.Keys.get("mmm").Values.get("ooo").Value != "ooo" {
		t.Error("should be ooo")
	}
	if len(kvp.Keys.get("mmm").Values.get("ooo").Tags) != 2 {
		t.Error("should tags len 2")
	}
	if kvp.Keys.get("mmm").Values.get("ooo").Tags.sorted()[0] != "T1" {
		t.Error("should tags T1")
	}
	if kvp.Keys.get("mmm").Values.get("ooo").Tags.sorted()[1] != "T2" {
		t.Error("should tags T2")
	}
	p, err = Parse("mmm@aaa[T3,T4]")
	if err != nil {
		t.Error(err)
	}
	p, err = p.Resolv(func(string, FuncsAndParam) (*string, error) { m := "ooo"; return &m, nil })
	if err != nil {
		t.Error(err)
	}
	sa, err = p.ToSetArgs()
	if err != nil {
		t.Error(err)
	}
	kvp.Set(*sa)
	if kvp.Keys.get("mmm").Values.get("ooo").Value != "ooo" {
		t.Error("should be ooo")
	}
	if *&kvp.Keys.get("mmm").Values.get("ooo").Unresolved.Param != "aaa" {
		t.Error("should be aaa")
	}
	if len(kvp.Keys.get("mmm").Values.get("ooo").Tags) != 4 {
		t.Error("should tags len 4")
	}
	if kvp.Keys.get("mmm").Values.get("ooo").Tags.sorted()[0] != "T1" {
		t.Error("should tags T1")
	}
	if kvp.Keys.get("mmm").Values.get("ooo").Tags.sorted()[1] != "T2" {
		t.Error("should tags T2")
	}
	if kvp.Keys.get("mmm").Values.get("ooo").Tags.sorted()[2] != "T3" {
		t.Error("should tags T3")
	}
	if kvp.Keys.get("mmm").Values.get("ooo").Tags.sorted()[3] != "T4" {
		t.Error("should tags T4")
	}

	p, err = Parse("mmm@aaa[T3,T4]")
	if err != nil {
		t.Error(err)
	}
	p, err = p.Resolv(func(string, FuncsAndParam) (*string, error) { m := "rrr"; return &m, nil })
	if err != nil {
		t.Error(err)
	}
	sa, err = p.ToSetArgs()
	if err != nil {
		t.Error(err)
	}
	kvp.Set(*sa)
	if kvp.Keys.get("mmm").Values.get("rrr").Value != "rrr" {
		t.Error("should be ooo")
	}
	if *&kvp.Keys.get("mmm").Values.get("rrr").Unresolved.Param != "aaa" {
		t.Error("should be aaa")
	}
	if len(kvp.Keys.get("mmm").Values.get("rrr").Tags) != 2 {
		t.Error("should tags len 2")
	}
	if kvp.Keys.get("mmm").Values.get("rrr").Tags.sorted()[0] != "T3" {
		t.Error("should tags T1", kvp.Keys.get("mmm").Values.get("rrr").Tags.sorted())
	}
	if kvp.Keys.get("mmm").Values.get("rrr").Tags.sorted()[1] != "T4" {
		t.Error("should tags T2", kvp.Keys.get("mmm").Values.get("rrr").Tags.sorted())
	}
}
