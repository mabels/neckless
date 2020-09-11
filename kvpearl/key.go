package kvpearl

import (
	"strings"
)

// JSONKey JSON Respresation of the Key
type JSONKey struct {
	Key    string
	Values JSONValues
}

// Key the KeyValue Structure
type Key struct {
	Key    string
	Values Values
}

// KeySorter KeyValue sorted by Key
type KeySorter []*JSONKey

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

func (key *Key) setValue(unresolved *FuncsAndParam, val string, tags map[string]int) *Value {
	// mergedValue := map[string](*Value){}
	// for i := range key.Values {
	// 	value, found := mergedValue[key.Values[i].Value]
	// 	if !found {
	// 		value = key.Values[i]
	// 		mergedValue[key.Values[i].Value] = value
	// 	}
	// 	value.Tags = uniqStrings(value.Tags)
	// }
	// value, found := mergedValue[val]
	value, _ := key.Values.getOrAdd(val)
	value.Unresolved = unresolved
	// if value.Unresolved != nil {
	// fmt.Println("SetUnresolved:", *value.Unresolved, val, tags)
	// }

	clearTags := []string{}
	for i := range tags {
		trimmed := strings.TrimSpace(i)
		if len(trimmed) != 0 {
			clearTags = append(clearTags, trimmed)
		}
	}
	value.Tags.add(clearTags...)

	// sort.Sort(&ValueSorter{
	// values: key.Values,
	// })
	// fmt.Println("YYYY=>", key.Key, len(key.Values), key.Values, tags, uniqStrings(append(value.Tags, tags...)))
	// value.Tags = uniqStrings(append(value.Tags, tags...))
	return value
}

func (key *Key) asJSON() *JSONKey {
	return &JSONKey{
		Key:    key.Key,
		Values: key.Values.RevOrdered().asJson(),
	}
}
