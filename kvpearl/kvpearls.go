package kvpearl

import (
	"sort"
	"strings"
)

// KVPearls is a array of KVPearl
type KVPearls []*KVPearl

// Len is part of sort.Interface.
func (kvps *KVPearls) Len() int {
	return len(*kvps)
}

// Swap is part of sort.Interface.
func (kvps *KVPearls) Swap(i, j int) {
	(*kvps)[i], (*kvps)[j] = (*kvps)[j], (*kvps)[i]
}

// Less is part of sort.Interface. It is implemented by calling the "by" closure in the sorter.
func (kvps *KVPearls) Less(i, j int) bool {
	return (*kvps)[i].Created.UnixNano() < (*kvps)[j].Created.UnixNano()
}

// AsJSON Converts KVPearls to JSONKVPearl
func (kvps *KVPearls) AsJSON() []*JSONKVPearl {
	out := make([]*JSONKVPearl, len(*kvps))
	obt := kvps.orderByTime()
	for i := range *obt {
		out[i] = (*obt)[i].AsJSON()
	}
	return out
}

func (kvps *KVPearls) orderByTime() *KVPearls {
	sort.Sort(kvps)
	return kvps
}

// MapByToResolve KVParsed ordered by ToResolve
type MapByToResolve map[string][]*KVParsed

// Add a KVParsed to MapByToResolve
func (mbr *MapByToResolve) Add(sa *KVParsed) {
	if sa == nil {
		return
	}
	toResolve := ""
	if sa.ToResolve != nil {
		toResolve = *sa.ToResolve
	}
	kvps, found := (*mbr)[toResolve]
	if !found {
		kvps = make([]*KVParsed, 0)
	}
	(*mbr)[toResolve] = append(kvps, sa)
}

// ByKeyValues the shim to sort the keys by the keyvalue
type ByKeyValues struct {
	Key  string
	Vals Values
}

// ArrayByKeyValues an array of ByKeyValues
type ArrayByKeyValues []*ByKeyValues

// JSONByKeyValues Respresentation of ByKeyValues
type JSONByKeyValues struct {
	Key  string
	Vals JSONValues
}

// ArrayOfJSONByKeyValues an array JSONByKeyValues
type ArrayOfJSONByKeyValues []JSONByKeyValues

func (s *ArrayOfJSONByKeyValues) Len() int {
	return len(*s)
}

// Swap is part of sort.Interface.
func (s *ArrayOfJSONByKeyValues) Swap(i, j int) {
	(*s)[i], (*s)[j] = (*s)[j], (*s)[i]
}

// Less is part of sort.Interface. It is implemented by calling the "by" closure in the sorter.
func (s *ArrayOfJSONByKeyValues) Less(i, j int) bool {
	return strings.Compare((*s)[i].Key, (*s)[j].Key) < 0
}

// MapByKeyValues map a "key/unresolved" to ArrayByKeyValues
type MapByKeyValues map[string]ArrayByKeyValues

// ToJSON return a Array of the json respresentation of ArrayByKeyValues
func (a *ArrayByKeyValues) ToJSON() ArrayOfJSONByKeyValues {
	ret := make(ArrayOfJSONByKeyValues, len(*a))
	for i := range *a {
		ret[i] = JSONByKeyValues{
			Key:  (*a)[i].Key,
			Vals: (*a)[i].Vals.RevOrdered().asJson(),
		}
	}
	sort.Sort(&ret)
	return ret
}

func (mbkv *MapByKeyValues) add(key *Key, val *Value) {
	unresolved := ""
	if val.Unresolved != nil {
		unresolved = *val.Unresolved
	}
	bkvs, found := (*mbkv)[unresolved]
	if !found {
		bkvs = make([]*ByKeyValues, 0)
	}
	var bkv *ByKeyValues
	for ibkv := range bkvs {
		tbkv := bkvs[ibkv]
		if tbkv.Key == key.Key {
			bkv = tbkv
			break
		}
	}
	if bkv == nil {
		bkv = &ByKeyValues{
			Key:  key.Key,
			Vals: createValues(),
		}
		bkvs = append(bkvs, bkv)
	}
	bkv.Vals.getOrAddValue(&Value{
		Value:      val.Value,
		Unresolved: &unresolved,
		Tags:       val.Tags,
		order:      val.order,
	})
	(*mbkv)[unresolved] = bkvs
}

// Match KVPearls against the MapByToResolve
func (kvps *KVPearls) Match(toResolves MapByToResolve) MapByKeyValues {
	ret := MapByKeyValues{}
	orderedKvps := kvps.orderByTime()
	// for toResolve := range toResolves {
	// 	kvparseds := toResolves[toResolve]
	// 	ret[toResolve] = make([]ByKeyValue, len(kvparseds))
	// 	for ikvparsed := range kvparseds {
	// 		kvparsed := kvparseds[ikvparsed]
	// 		ret.oneMatch(orderedKvps, kvparsed)
	// 	}
	// }
	// if len(toResolves) == 0 {
	// fmt.Printf("------- 0000\n")
	// tmp := make([]ByKeyValue, 0)
	for o := range *orderedKvps {
		// fmt.Printf("------- 0000:%d\n", o)
		oval := (*orderedKvps)[o]
		for i := range oval.Keys {
			// fmt.Printf("------- 0000:%d:%s\n", o, i)
			key := oval.Keys[i]
			rev := key.Values.RevOrdered()
			for k := range *rev {
				// fmt.Printf("------- 0000:%d:%s:%d\n", o, i, k)
				val := (*rev)[k]
				unresolved := "" // taken from match
				if len(toResolves) == 0 {
					ret.add(key, &Value{
						Value:      val.Value,
						Unresolved: &unresolved,
						Tags:       val.Tags,
						order:      val.order,
					})
				}
				for ires := range toResolves {
					for ikvps := range toResolves[ires] {
						kvp := toResolves[ires][ikvps]
						matchKVP, match := kvp.Match(key, val)
						if match {
							if matchKVP.ToResolve != nil {
								unresolved = *matchKVP.ToResolve
							}
							ret.add(key, &Value{
								Value:      val.Value,
								Unresolved: &unresolved,
								Tags:       val.Tags,
								order:      val.order,
							})
						}
					}
				}
			}
		}
	}
	return ret
}
