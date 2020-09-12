package kvpearl

import (
	"encoding/json"
	"fmt"
	"regexp"
	"testing"
	"time"
)

func (mptr *MapByToResolve) add(key string, tags ...string) *MapByToResolve {
	keyreg := regexp.MustCompile(fmt.Sprintf("^%s$", key))
	if len(key) == 0 {
		keyreg = regexp.MustCompile("^.*$")
	}
	mytags := Tags{}
	for i := range tags {
		mytags[tags[i]] = 1
	}
	val := fmt.Sprintf("Val:%d:%s", len(*mptr), key)
	kvp := KVParsed{
		Key:       &key,
		KeyRegex:  keyreg,
		ToResolve: nil,
		Val:       &val,
		Tags:      mytags,
	}
	(*mptr)[key] = [](*KVParsed){&kvp}
	return mptr
}

func (kvps *KVPearls) add(key string, tags ...string) *KVPearls {
	order := 45
	values := createValues(&order)
	mytags := Tags{}
	for i := range tags {
		mytags[tags[i]] = 1
	}
	(*kvps.orderRef)++
	values.getOrAddValue(&Value{
		Value: fmt.Sprintf("%d:%s", *kvps.orderRef, key),
		Tags:  mytags,
		// order: *kvps.orderRef,
	})
	keys := keys{}
	keys[key] = &Key{
		Key:    key,
		Values: values,
	}
	kvps.kvps = append(kvps.kvps, &KVPearl{
		Keys:    keys,
		order:   kvps.order,
		Created: time.Now(),
	})
	return kvps
}

func TestEmptyMatchTag(t *testing.T) {
	ma := CreateKVPearls()
	mptr := MapByToResolve{}
	mptr.add("key")
	mkvp := ma.Match(mptr)[""]
	if len(mkvp) != 0 {
		t.Error("should be true")
	}
}

func TestEmptyMatchEmptyTags(t *testing.T) {
	ma := CreateKVPearls().add("test1").add("test2")

	m0 := MapByToResolve{}
	if mkvp := ma.Match(m0)[""]; len(mkvp) != 2 {
		t.Error("should be true", len(mkvp))
	}
	m1 := MapByToResolve{}
	m1.add("test1")
	if mkvp := ma.Match(m1)[""]; len(mkvp) != 1 {
		t.Error("should be true")
	}
	m2 := MapByToResolve{}
	m2.add("test2")
	if mkvp := ma.Match(m2)[""]; len(mkvp) != 1 {
		t.Error("should be true")
	}
}

func TestMergeSorted(t *testing.T) {
	kvps := CreateKVPearls()
	kvps.Add().Set(SetArg{
		Key:        "test1",
		Unresolved: &FuncsAndParam{},
		Val:        "vtest1",
		Tags:       []string{"Test1"},
	})
	kvps.Add().Set(SetArg{
		Key:        "test1",
		Unresolved: &FuncsAndParam{},
		Val:        "vtest2",
		Tags:       []string{"Test1"},
	})
	kvps.Add().Set(SetArg{
		Key:        "test1",
		Unresolved: &FuncsAndParam{},
		Val:        "vtest1",
		Tags:       []string{"Test2"},
	})
	kvps.Add().Set(SetArg{
		Key:        "test1",
		Unresolved: &FuncsAndParam{},
		Val:        "vtest3",
		Tags:       []string{"Test3"},
	}).Set(SetArg{Key: "test0", Val: "vtest0"})
	so := kvps.Merge()
	if len(so) != 2 {
		t.Error("should be length 3")
	}
	if !(so[0].Key == "test0" && len(so[0].Values) == 1 && so[0].Values[0].Value == "vtest0") {
		t.Error("is not s[0]")
	}
	if !(so[1].Key == "test1" && len(so[1].Values) == 3 &&
		so[1].Values[0].Value == "vtest2" && so[1].Values[0].Tags[0] == "Test1" &&
		so[1].Values[1].Value == "vtest1" && so[1].Values[1].Tags[0] == "Test1" && so[1].Values[1].Tags[1] == "Test2" &&
		so[1].Values[2].Value == "vtest3" && so[1].Values[2].Tags[0] == "Test3") {
		js, _ := json.MarshalIndent(so[1], "", "  ")
		t.Errorf("is not s[1]:%s", string(js))
	}

}

func TestEmptyMatchTags(t *testing.T) {
	ma := CreateKVPearls().add("test1", "t1T1", "t1T2").add("test2", "t2T1", "t2T2")

	mx := MapByToResolve{}
	mx.add("xx")
	if len(ma.Match(mx)[""]) != 0 {
		t.Error("should be true")
	}
	m1 := MapByToResolve{}
	m1.add("test1")
	if len(ma.Match(m1)[""]) != 1 {
		t.Error("should be true", len(ma.Match(m1)))
	}
	m2 := MapByToResolve{}
	m2.add("test2")
	if len(ma.Match(m2)[""]) != 1 {
		t.Error("should be true")
	}

	unMatch1 := MapByToResolve{}
	unMatch1.add("test1", "unMatch")
	if len(ma.Match(unMatch1)[""]) != 0 {
		t.Error("should be true")
	}
	unMatch2 := MapByToResolve{}
	unMatch2.add("test2", "unMatch")
	if len(ma.Match(unMatch2)[""]) != 0 {
		t.Error("should be true")
	}

	matcht11T1 := MapByToResolve{}
	matcht11T1.add("test1", "t1T1")
	if len(ma.Match(matcht11T1)[""]) != 1 {
		t.Error("should be true")
	}
	matcht21T1 := MapByToResolve{}
	matcht21T1.add("test2", "t1T1")
	if len(ma.Match(matcht21T1)[""]) != 0 {
		t.Error("should be true")
	}

	matcht11T2 := MapByToResolve{}
	matcht11T2.add("test1", "t1T2")
	if len(ma.Match(matcht11T2)[""]) != 1 {
		t.Error("should be true")
	}
	matcht21T2 := MapByToResolve{}
	matcht21T2.add("test2", "t1T2")
	if len(ma.Match(matcht21T2)[""]) != 0 {
		t.Error("should be true")
	}

}

func TestNoTags(t *testing.T) {
	ma := CreateKVPearls().add("test1").add("test2")
	xxx := MapByToResolve{}
	xxx.add("xxx")
	if len(ma.Match(xxx)[""]) != 0 {
		t.Error("Should fail")
	}
	test1 := MapByToResolve{}
	test1.add("test1")
	if len(ma.Match(test1)[""]) != 1 {
		t.Error("Should match")
	}
	test2 := MapByToResolve{}
	test2.add("test2")
	if len(ma.Match(test2)[""]) != 1 {
		t.Error("Should match")
	}
}
