package kvpearl

import (
	"encoding/json"
	"sort"
	"strings"

	"neckless.adviser.com/key"
	"neckless.adviser.com/pearl"
)

type Value struct {
	Value string
	Tags  []string
}
type Key struct {
	Key    string
	Values [](*Value)
}

type KVPearl struct {
	// Tags []string
	Keys map[string](*Key)
}

func uniqStrings(strs []string) []string {
	set := map[string](struct{}){}
	for i := range strs {
		set[strs[i]] = struct{}{}
	}
	ret := make([]string, len(set))
	// fmt.Printf("uniq-1:%d\n", len(set), len(ret))
	idx := 0
	for i := range set {
		ret[idx] = i
		idx++
		// fmt.Printf("i:%s:%d\n", i, len(ret))
	}
	// fmt.Println("xxxx:", ret)
	sort.Strings(ret)
	return ret
}

func Create( /* tags ...string */ ) *KVPearl {
	return &KVPearl{
		// Tags: uniqStrings(tags),
		Keys: map[string](*Key){},
	}
}

func ValueBy(p1, p2 *Value) bool {
	return strings.Compare(p1.Value, p2.Value) < 0
}

type ValueSorter struct {
	values [](*Value)
	by     func(p1, p2 *Value) bool // Closure used in the Less method.
}

// Len is part of sort.Interface.
func (s *ValueSorter) Len() int {
	return len(s.values)
}

// Swap is part of sort.Interface.
func (s *ValueSorter) Swap(i, j int) {
	s.values[i], s.values[j] = s.values[j], s.values[i]
}

// Less is part of sort.Interface. It is implemented by calling the "by" closure in the sorter.
func (s *ValueSorter) Less(i, j int) bool {
	return s.by(s.values[i], s.values[j])
}

func (key *Key) setValue(val string, tags []string) *Value {
	mergedValue := map[string](*Value){}
	for i := range key.Values {
		value, found := mergedValue[key.Values[i].Value]
		if !found {
			value = key.Values[i]
			mergedValue[key.Values[i].Value] = value
		}
		value.Tags = uniqStrings(value.Tags)
	}
	value, found := mergedValue[val]
	if !found {
		value = &Value{
			Value: val,
			Tags:  []string{},
		}
		mergedValue[val] = value
	}
	value.Tags = uniqStrings(append(value.Tags, tags...))

	// fmt.Println("XXXX", len(key.Values), len(mergedValue), mergedValue)
	if len(key.Values) < len(mergedValue) {
		for i := len(key.Values); i < len(mergedValue); i++ {
			key.Values = append(key.Values, &Value{
				Value: "",
				Tags:  []string{},
			})
			// fmt.Println("APPEND-To-Values", len(key.Values))
		}
	}
	// key.Values = make([]Value, len(mergedValue))
	idx := 0
	for i := range mergedValue {
		key.Values[idx] = mergedValue[i]
		idx++
	}

	sort.Sort(&ValueSorter{
		values: key.Values,
		by:     ValueBy,
	})
	// fmt.Println("YYYY=>", key.Key, len(key.Values), key.Values, tags, uniqStrings(append(value.Tags, tags...)))
	// value.Tags = uniqStrings(append(value.Tags, tags...))
	return value
}

func (kvp *KVPearl) Set(keyVal string, val string, tags ...string) *KVPearl {
	key, found := kvp.Keys[keyVal]
	if !found {
		key = &Key{
			Key:    keyVal,
			Values: []*Value{},
		}
		kvp.Keys[keyVal] = key
		// fmt.Println("Set=>", kvp, key, keyVal, val, tags)
	} else {
		// fmt.Println("Found-Set=>", kvp, key, keyVal, val, tags)
	}
	key.setValue(val, tags)
	// fmt.Println("Post-Set-Values", kvp, key, len(key.Values))
	return kvp
}

func FromJson(jsStr []byte) (*KVPearl, error) {
	kvp := Create()
	err := json.Unmarshal(jsStr, &kvp)
	if err != nil {
		return nil, err
	}
	// kvp.Tags = uniqStrings(kvp.Tags)
	for i := range kvp.Keys {
		key := kvp.Keys[i]
		for j := range key.Values {
			kvp.Set(i, key.Values[j].Value, key.Values[j].Tags...)
		}
	}
	return kvp, nil
}

func (kvp *KVPearl) AsJson() *KVPearl {
	return kvp
}

func (kvp *KVPearl) ClosePearl(owners *pearl.PearlOwner) (*pearl.Pearl, error) {
	jsonStr, err := json.Marshal(kvp.AsJson())
	if err != nil {
		return nil, err
	}
	return pearl.Close(&pearl.CloseRequestPearl{
		Type:    "KVPearl",
		Payload: jsonStr,
		Owners:  *owners,
	})
}

func OpenPearl(pk *key.PrivateKey, prl *pearl.Pearl) (*KVPearl, error) {
	op, err := pearl.Open(pk, prl)
	if err != nil {
		return nil, err
	}
	return FromJson(op.Payload)
}

type KVPearls []KVPearl

type KeyValue struct {
	Key   string
	Value string
}

func KeyValueBy(p1, p2 *KeyValue) bool {
	return strings.Compare(p1.Key, p2.Key) < 0
}

type KeyValueSorter struct {
	values []KeyValue
	by     func(p1, p2 *KeyValue) bool // Closure used in the Less method.
}

// Len is part of sort.Interface.
func (s *KeyValueSorter) Len() int {
	return len(s.values)
}

// Swap is part of sort.Interface.
func (s *KeyValueSorter) Swap(i, j int) {
	s.values[i], s.values[j] = s.values[j], s.values[i]
}

// Less is part of sort.Interface. It is implemented by calling the "by" closure in the sorter.
func (s *KeyValueSorter) Less(i, j int) bool {
	return s.by(&s.values[i], &s.values[j])
}

func findTag(t1, t2 []string) bool {
	if len(t2) == 0 {
		return true
	}
	if len(t1) == 0 {
		return false
	}
	var ref []string
	var toSet []string
	if len(t1) > len(t2) {
		ref = t2
		toSet = t1
	} else {
		ref = t1
		toSet = t2
	}
	set := map[string](struct{}){}
	for i := range toSet {
		set[toSet[i]] = struct{}{}
	}
	for i := range ref {
		_, found := set[ref[i]]
		if found {
			// fmt.Println("FOUND:", t1, t2, set)
			return true
		}
	}
	// fmt.Println("!FOUND:", t1, t2, set)
	return false
}

func matchTag(kvp KVPearl, tags []string) []KeyValue {
	// if len(kvp.Tags) != 0 {
	// 	if !findTag(kvp.Tags, tags) {
	// 		return []KeyValue{}
	// 	}
	// }
	mapRet := map[string]KeyValue{}
	for i := range kvp.Keys {
		key := kvp.Keys[i]
		for j := range key.Values {
			if findTag(key.Values[j].Tags, tags) {
				mapRet[key.Key] = KeyValue{
					Key:   key.Key,
					Value: key.Values[j].Value,
				}
			}
		}

	}
	ret := make([]KeyValue, len(mapRet))
	idx := 0
	for i := range mapRet {
		ret[idx] = mapRet[i]
		idx++
	}
	return ret
}

func (kvps *KVPearls) Ls(tags ...string) []KeyValue {
	// ret := KVPearls{}
	byKey := map[string]KeyValue{}
	for i := range *kvps {
		kvpearl := (*kvps)[i]
		kvs := matchTag(kvpearl, tags)
		for j := range kvs {
			kv := kvs[j]
			byKey[kv.Key] = KeyValue{
				Value: kv.Value,
				Key:   kv.Key,
			}
		}
	}
	out := make([]KeyValue, len(byKey))
	idx := 0
	for i := range byKey {
		out[idx] = byKey[i]
		idx++
	}
	// fmt.Println("In=>", out)
	sort.Sort(&KeyValueSorter{
		values: out,
		by:     KeyValueBy,
	})
	// fmt.Println("Out=>", out)
	return out
}

func (kvps *KVPearls) Get(name string, tags ...string) *KeyValue {
	out := kvps.Ls(tags...)
	for i := range out {
		if strings.Compare(out[i].Key, name) == 0 {
			return &out[i]
		}
	}
	return nil
}
