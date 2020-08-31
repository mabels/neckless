package kvpearl

import (
	"sort"
	"strings"
)

type KVPearls []*KVPearl

// type jsonKVPearlSorter []*JsonKVPearl

// Len is part of sort.Interface.
func (s *KVPearls) Len() int {
	return len(*s)
}

// Swap is part of sort.Interface.
func (s *KVPearls) Swap(i, j int) {
	(*s)[i], (*s)[j] = (*s)[j], (*s)[i]
}

// Less is part of sort.Interface. It is implemented by calling the "by" closure in the sorter.
func (s *KVPearls) Less(i, j int) bool {
	return (*s)[i].Created.UnixNano() < (*s)[j].Created.UnixNano()
}

func (kvps *KVPearls) AsJSON() []*JsonKVPearl {
	out := make([]*JsonKVPearl, len(*kvps))
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

// type MergeArgs map[string]*KVParsed

type MapByToResolve map[string][]*KVParsed

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

func (kvp *KVParsed) Match(key *Key, val *Value) (*KVParsed, bool) {
	// findKey := false
	// for ikvps := range kvps {
	// kvp := kvps[ikvps]
	if kvp.KeyRegex.MatchString(key.Key) {
		// fmt.Printf("Matched:%s:%d", kvp.Key, len(kvp.Tags))
		if len(kvp.Tags) == 0 {
			return kvp, true
		}
		for tag := range val.Tags {
			// fmt.Printf("%s:%s\n", tag, kvp.Tags)
			_, found := kvp.Tags[tag]
			if found {
				return kvp, true
			}
		}
	}

	// if *kvp.Key != key {
	// 	return nil, false
	// }
	// if len(rtags) == 0 {
	// 	return nil, true
	// }
	// if len(kvp.Tags) == 0 {
	// 	return kvp, true
	// }
	// for tag := range rtags {
	// 	_, found := kvp.Tags[tag]
	// 	if found {
	// 		return kvp, true
	// 	}
	// }
	return nil, false
}

type ByKeyValues struct {
	Key  string
	Vals Values
}

type ArrayByKeyValues []*ByKeyValues

type JsonByKeyValues struct {
	Key  string
	Vals JsonValues
}
type ArrayOfJsonByKeyValues []JsonByKeyValues

func (s *ArrayOfJsonByKeyValues) Len() int {
	return len(*s)
}

// Swap is part of sort.Interface.
func (s *ArrayOfJsonByKeyValues) Swap(i, j int) {
	(*s)[i], (*s)[j] = (*s)[j], (*s)[i]
}

// Less is part of sort.Interface. It is implemented by calling the "by" closure in the sorter.
func (s *ArrayOfJsonByKeyValues) Less(i, j int) bool {
	return strings.Compare((*s)[i].Key, (*s)[j].Key) < 0
}

type MapByKeyValues map[string]ArrayByKeyValues

func (a *ArrayByKeyValues) ToJson() []JsonByKeyValues {
	ret := make(ArrayOfJsonByKeyValues, len(*a))
	for i := range *a {
		ret[i] = JsonByKeyValues{
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
	// 	ret[""] = tmp
	// 	// fmt.Printf("------- done:%d:%d\n", len(ret), len(tmp))
	// }
	return ret
}

// func (mapBy *MapByKeyValues) oneMatch(orderedKvps *KVPearls, kvparsed *KVParsed) {
// 	for k := range *orderedKvps {
// 		kvp := (*orderedKvps)[k]
// 		for i := range kvp.Keys {
// 			key := kvp.Keys[i]
// 			rev := *key.Values.RevOrdered()
// 			for j := range rev {
// 				val := rev[j]
// 				matchedKVP, found := kvparsed.Match(key.Key, val.Tags)
// 				if found {
// 					unresolved := ""
// 					if matchedKVP.ToResolve != nil {
// 						unresolved = *matchedKVP.ToResolve
// 					}
// 					kvs, found := (*mapBy)[unresolved]
// 					if !found {
// 						kvs = make([]ByKeyValues, 0)
// 					}
// 					my := ByKeyValues{
// 						Key: key.Key,
// 						Val: Value{
// 							Value:      val.Value,
// 							Unresolved: &unresolved,
// 							Tags:       val.Tags,
// 							order:      val.order,
// 						},
// 					}
// 					(*mapBy)[unresolved] = append(kvs, my)
// 				}
// 			}
// 		}
// 	}
// }

// func (kvps *KVPearls) Merge(keys MergeArgs) *KVPearl {
// 	kvp := Create()
// 	for i := range *kvps.orderByTime() {
// 		for j := range (*kvps)[i].Keys {
// 			key := (*kvps)[i].Keys[j]
// 			vals := key.Values.RevOrdered()
// 			for k := range *vals {
// 				value := (*vals)[k]
// 				matchKVP, found := keys.Match(key.Key, value.Tags)
// 				if found {
// 					var unresolved *string
// 					if matchKVP != nil {
// 						unresolved = matchKVP.ToResolve
// 					}
// 					kvp.Set(SetArg{
// 						Key:        key.Key,
// 						Unresolved: unresolved,
// 						Val:        value.Value,
// 						Tags:       value.Tags.sorted(),
// 					})
// 				}
// 			}
// 		}
// 	}
// 	kvp.Created = time.Now()
// 	return kvp
// }
