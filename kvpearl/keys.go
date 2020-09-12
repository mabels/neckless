package kvpearl

import (
	"sort"
	"strings"
)

type keys map[string](*Key)
type sortedKeys []JSONKey

func (ks *keys) get(val string) *Key {
	key, found := (*ks)[val]
	if found {
		return key
	}
	return nil
}

func (ks *keys) getOrAdd(val string, orderRef *int) (*Key, bool) {
	key, found := (*ks)[val]
	if !found {
		key = &Key{
			Key:    val,
			Values: createValues(orderRef),
		}
		(*ks)[val] = key
	}
	return key, found
}

// Len is part of sort.Interface.
func (sk *sortedKeys) Len() int {
	return len(*sk)
}

// Swap is part of sort.Interface.
func (sk *sortedKeys) Swap(i, j int) {
	(*sk)[i], (*sk)[j] = (*sk)[j], (*sk)[i]
}

// Less is part of sort.Interface. It is implemented by calling the "by" closure in the sorter.
func (sk *sortedKeys) Less(i, j int) bool {
	return strings.Compare((*sk)[i].Key, (*sk)[j].Key) < 0
}

// func (sk *sortedKeys) asJSON() []JSONKey {
// 	ret := make([]JSONKey, len(*sk))
// 	for i := range *sk {
// 		ret[i] = (*sk)[i].asJSON()
// 	}
// 	return ret
// }

func (sk *sortedKeys) AsStrings() []string {
	ret := make([]string, len(*sk))
	for i := range *sk {
		ret[i] = (*sk)[i].Key
	}
	return ret
}

func (ks *keys) Sorted() sortedKeys {
	jsKeys := make(sortedKeys, len(*ks))
	keyIdx := 0
	for i := range *ks {
		key := (*ks)[i]
		jsKeys[keyIdx] = key.asJSON()
		keyIdx++
	}
	sort.Sort(&jsKeys)
	return jsKeys
}
