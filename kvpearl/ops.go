package kvpearl

import (
	"sort"
	"strings"
)

// Ls operation by tags
func (kvps *KVPearls) Ls(tags ...string) KeyValues {
	// ret := KVPearls{}
	byKey := map[string]KeyValue{}
	for i := range *kvps {
		kvpearl := (*kvps)[i]
		kvs := kvpearl.matchTag(tags)
		for j := range kvs {
			kv := kvs[j]
			byKey[kv.Key] = KeyValue{
				Value: kv.Value,
				Key:   kv.Key,
			}
		}
	}
	out := make(KeyValues, len(byKey))
	idx := 0
	for i := range byKey {
		out[idx] = byKey[i]
		idx++
	}
	// fmt.Println("In=>", out)
	sort.Sort(&out)
	// fmt.Println("Out=>", out)
	return out
}

// Get operation
func (kvps *KVPearls) Get(name string, tags ...string) *KeyValue {
	out := kvps.Ls(tags...)
	for i := range out {
		if strings.Compare(out[i].Key, name) == 0 {
			return &out[i]
		}
	}
	return nil
}
