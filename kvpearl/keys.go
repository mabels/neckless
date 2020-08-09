package kvpearl

import (
	"sort"
	"strings"
)

type keys map[string](*Key)

func (ks *keys) get(val string) *Key {
	key, found := (*ks)[val]
	if found {
		return key
	}
	return nil
}

func (ks *keys) getOrAdd(val string) (*Key, bool) {
	key, found := (*ks)[val]
	if !found {
		key = &Key{
			Key:    val,
			Values: createValues(),
		}
		(*ks)[val] = key
	}
	return key, found
}

type sortedKeys [](*Key)

// Len is part of sort.Interface.
func (s *sortedKeys) Len() int {
	return len(*s)
}

// Swap is part of sort.Interface.
func (s *sortedKeys) Swap(i, j int) {
	(*s)[i], (*s)[j] = (*s)[j], (*s)[i]
}

// Less is part of sort.Interface. It is implemented by calling the "by" closure in the sorter.
func (s *sortedKeys) Less(i, j int) bool {
	return strings.Compare((*s)[i].Key, (*s)[j].Key) < 0
}

func (sk *sortedKeys) asJson() []*JsonKey {
	ret := make([]*JsonKey, len(*sk))
	for i := range *sk {
		ret[i] = (*sk)[i].asJson()
	}
	return ret
}

func (sk *sortedKeys) AsStrings() []string {
	ret := make([]string, len(*sk))
	for i := range *sk {
		ret[i] = (*sk)[i].Key
	}
	return ret
}

func (ks *keys) Sorted() *sortedKeys {
	jsKeys := make(sortedKeys, len(*ks))
	keyIdx := 0
	for i := range *ks {
		key := (*ks)[i]
		jsKeys[keyIdx] = key
		keyIdx++
	}
	sort.Sort(&jsKeys)
	return &jsKeys
}
