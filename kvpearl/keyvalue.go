package kvpearl

import "strings"

// KeyValue pair with just one Value
type KeyValue struct {
	Key   string
	Value string
}

// KeyValues is a array of a KeyValue
type KeyValues []KeyValue

// Len is part of sort.Interface.
func (s *KeyValues) Len() int {
	return len(*s)
}

// Swap is part of sort.Interface.
func (s *KeyValues) Swap(i, j int) {
	(*s)[i], (*s)[j] = (*s)[j], (*s)[i]
}

// Less is part of sort.Interface. It is implemented by calling the "by" closure in the sorter.
func (s *KeyValues) Less(i, j int) bool {
	return strings.Compare((*s)[i].Key, (*s)[j].Key) < 0
}
