package kvpearl

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strings"
	"time"

	"neckless.adviser.com/key"
	"neckless.adviser.com/pearl"
)

type Value struct {
	Value string
	Order time.Time // `json:"-"`
	Tags  []string
}
type Key struct {
	Key    string
	Values [](*Value)
}

type KeySorter []*Key

// Len is part of sort.Interface.
func (s *KeySorter) Len() int {
	return len(*s)
}

// Swap is part of sort.Interface.
func (s *KeySorter) Swap(i, j int) {
	(*s)[i], (*s)[j] = (*s)[j], (*s)[i]
}

// Less is part of sort.Interface. It is implemented by calling the "by" closure in the sorter.
func (s *KeySorter) Less(i, j int) bool {
	return strings.Compare((*s)[i].Key, (*s)[j].Key) < 0
}

type KVPearl struct {
	// Tags []string
	Keys    map[string](*Key)
	Created time.Time
	Pearl   *pearl.OpenPearl `json:"-"`
}

type JsonKVPearl struct {
	// Tags []string
	Keys        []*Key
	FingerPrint string `json:"FingerPrint,omitempty"`
	Created     time.Time
}

func JsonKVPearlValueBy(p1, p2 *JsonKVPearl) bool {
	return p1.Created.UnixNano() < p2.Created.UnixNano()
}

type JsonKVPearlSorter struct {
	Values [](*JsonKVPearl)
	By     func(p1, p2 *JsonKVPearl) bool // Closure used in the Less method.
}

// Len is part of sort.Interface.
func (s *JsonKVPearlSorter) Len() int {
	return len(s.Values)
}

// Swap is part of sort.Interface.
func (s *JsonKVPearlSorter) Swap(i, j int) {
	s.Values[i], s.Values[j] = s.Values[j], s.Values[i]
}

// Less is part of sort.Interface. It is implemented by calling the "by" closure in the sorter.
func (s *JsonKVPearlSorter) Less(i, j int) bool {
	return s.By(s.Values[i], s.Values[j])
}

func AsJSON(kvps []*KVPearl) []*JsonKVPearl {
	out := make([]*JsonKVPearl, len(kvps))
	for i := range kvps {
		out[i] = kvps[i].AsJSON()
	}
	sort.Sort(&JsonKVPearlSorter{
		Values: out,
		By:     JsonKVPearlValueBy,
	})
	return out
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
		Keys:    map[string](*Key){},
		Created: time.Now(),
	}
}

func (key *Key) setValue(order time.Time, val string, tags []string) *Value {
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
			Order: order,
			Tags:  []string{},
		}
		mergedValue[val] = value
	}
	clearTags := []string{}
	for i := range tags {
		trimmed := strings.TrimSpace(tags[i])
		if len(trimmed) != 0 {
			clearTags = append(clearTags, trimmed)
		}
	}

	value.Tags = uniqStrings(append(value.Tags, clearTags...))

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
	})
	// fmt.Println("YYYY=>", key.Key, len(key.Values), key.Values, tags, uniqStrings(append(value.Tags, tags...)))
	// value.Tags = uniqStrings(append(value.Tags, tags...))
	return value
}

func (kvp *KVPearl) Set(order time.Time, keyVal string, val string, tags ...string) *KVPearl {
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
	key.setValue(order, val, tags)
	// fmt.Println("Post-Set-Values", kvp, key, len(key.Values))
	return kvp
}

func FromJSON(jsStr []byte) (*KVPearl, error) {
	jskvp := JsonKVPearl{}
	err := json.Unmarshal(jsStr, &jskvp)
	if err != nil {
		return nil, err
	}
	kvp := Create()
	kvp.Created = jskvp.Created
	// kvp.Tags = uniqStrings(kvp.Tags)
	for i := range jskvp.Keys {
		key := jskvp.Keys[i]
		for j := range key.Values {
			kvp.Set(kvp.Created, key.Key, key.Values[j].Value, key.Values[j].Tags...)
		}
		kvp.Keys[key.Key] = key
	}
	return kvp, nil
}

type ValueSorter struct {
	values []*Value
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
	return s.values[i].Order.UnixNano() > s.values[j].Order.UnixNano()
}

func (kvp *KVPearl) AsJSON() *JsonKVPearl {
	keys := KeySorter{}
	for i := range kvp.Keys {
		key := kvp.Keys[i]
		sort.Sort(&ValueSorter{values: key.Values})
		keys = append(keys, key)
	}
	sort.Sort(&keys)
	// for i := range keys {
	// 	fmt.Println(keys[i].Key)
	// }
	fpr := ""
	if kvp.Pearl != nil {
		fpr = base64.StdEncoding.EncodeToString(kvp.Pearl.Closed.FingerPrint)
	}
	return &JsonKVPearl{
		Created:     kvp.Created,
		FingerPrint: fpr,
		Keys:        keys,
	}
}

const Type = "KVPearl"

func (kvp *KVPearl) ClosePearl(owners *pearl.PearlOwner) (*pearl.Pearl, error) {
	jsonStr, err := json.Marshal(kvp.AsJSON())
	if err != nil {
		return nil, err
	}
	return pearl.Close(&pearl.CloseRequestPearl{
		Type:    Type,
		Payload: jsonStr,
		Owners:  *owners,
	})
}

func OpenPearl(pks []*key.PrivateKey, prl *pearl.Pearl) (*KVPearl, error) {
	op, err := pearl.Open(pks, prl)
	if err != nil {
		return nil, err
	}
	kvp, err := FromJSON(op.Payload)
	if err != nil {
		return nil, err
	}
	kvp.Pearl = op
	return kvp, err
}

type KVPearls []KVPearl

type KeyValue struct {
	Key   string
	Value string
}

type KeyValueSorter struct {
	values []KeyValue
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
	return strings.Compare(s.values[i].Key, s.values[j].Key) < 0
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
				break
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

func (kv *KVPearl) Parse(arg string) (*KVPearl, error) {
	argRegex := regexp.MustCompile(`^([^=]+)=([^\[]+)(\[([^\]]*)\])*$`)
	split := argRegex.FindStringSubmatch(arg)
	if len(split) != 5 {
		return nil, errors.New(fmt.Sprintf("no matching kv:[%s]", arg))
	}
	// js, _ := json.Marshal(split)
	// fmt.Println(string(js))

	if len(split) > 3 && len(split[4]) > 0 {
		// fmt.Println("split set tag")
		kv.Set(time.Now(), split[1], split[2], strings.Split(split[4], ",")...)
	} else {
		// fmt.Println("split set")
		kv.Set(time.Now(), split[1], split[2])
	}
	return kv, nil
}

func Merge(kvps []*KVPearl, keys []string, tags []string) *KVPearl {
	kvp := Create()
	mapKeys := map[string]struct{}{}
	for i := range keys {
		key := keys[i]
		mapKeys[key] = struct{}{}
	}
	mapTags := map[string]struct{}{}
	for i := range tags {
		key := tags[i]
		mapTags[key] = struct{}{}
	}
	for i := range kvps {
		for j := range kvps[i].Keys {
			key := kvps[i].Keys[j]
			_, found := mapKeys[key.Key]
			if len(keys) == 0 || found {
				for k := range key.Values {
					kvp.Set(key.Values[k].Order, key.Key, key.Values[k].Value, key.Values[k].Tags...)
				}
			}
		}
	}
	kvp.Created = time.Now()
	return kvp
}
