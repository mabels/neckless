package kvpearl

import (
	"sort"
	"strings"
	"time"
)

// KVPearls is a array of KVPearl
type KVPearls struct {
	kvps     []*KVPearl
	orderRef *int
	order    int
}

// CreateKVPearls a chain of Pearls
func CreateKVPearls(orderps ...*int) *KVPearls {
	order := 0x0abcdef
	orderp := &order
	if len(orderps) != 0 {
		orderp = orderps[0]
	}
	(*orderp)++
	return &KVPearls{
		orderRef: orderp,
		order:    *orderp,
		kvps:     []*KVPearl{},
	}
}

// Len is part of sort.Interface.
func (kvps *KVPearls) Len() int {
	return len(kvps.kvps)
}

// Swap is part of sort.Interface.
func (kvps *KVPearls) Swap(i, j int) {
	(kvps.kvps)[i], (kvps.kvps)[j] = (kvps.kvps)[j], (kvps.kvps)[i]
}

// Less is part of sort.Interface. It is implemented by calling the "by" closure in the sorter.
func (kvps *KVPearls) Less(i, j int) bool {
	if kvps.kvps[i].Created.UnixNano() == kvps.kvps[j].Created.UnixNano() {
		return kvps.kvps[i].order < kvps.kvps[j].order
	}
	return (kvps.kvps)[i].Created.UnixNano() < kvps.kvps[j].Created.UnixNano()
}

// Add a KVPearl
func (kvps *KVPearls) Add(kvp ...*KVPearl) *KVPearl {
	if len(kvp) == 0 {
		(*kvps.orderRef)++
		ret := &KVPearl{
			// Tags: uniqStrings(tags),
			Keys:     keys{},
			orderRef: kvps.orderRef,
			order:    *kvps.orderRef,
			Created:  time.Now(),
		}
		kvps.kvps = append(kvps.kvps, ret)
		return ret
	}
	for i := range kvp {
		my := KVPearl{
			Keys:    kvp[i].Keys,
			order:   kvps.order,
			Created: kvp[i].Created,
			Pearl:   kvp[i].Pearl,
		}
		kvps.kvps = append(kvps.kvps, &my)

	}
	return kvp[len(kvp)-1]
}

// AsJSON Converts KVPearls to JSONKVPearl
func (kvps *KVPearls) AsJSON() []*JSONKVPearl {
	out := make([]*JSONKVPearl, len(kvps.kvps))
	obt := kvps.orderByTime()
	for i := range obt.kvps {
		out[i] = (obt.kvps)[i].AsJSON()
	}
	return out
}

func (kvps *KVPearls) orderByTime() *KVPearls {
	sort.Sort(kvps)
	return kvps
}

// Merge the KVPearls
func (kvps *KVPearls) Merge() sortedKeys {
	ret := keys{}
	orderedKvps := kvps.orderByTime()
	for okvp := range orderedKvps.kvps {
		kvp := orderedKvps.kvps[okvp]
		for keyStr := range kvp.Keys {
			keyp, found := ret[keyStr]
			if !found {
				keyp = &Key{
					Key:    keyStr,
					Values: createValues(kvps.orderRef),
				}
				ret[keyStr] = keyp
			}
			orderVal := kvp.Keys[keyStr].Values.Ordered()
			for o := range *orderVal {
				val := (*orderVal)[o]
				// fmt.Println(val.Value)
				keyp.Values.getOrAddValue(val)
			}
		}
	}
	return ret.Sorted()
}

// MapByToResolve KVParsed ordered by ToResolve
type MapByToResolve map[string][]*KVParsed

// Add a KVParsed to MapByToResolve
func (mbr *MapByToResolve) Add(sa *KVParsed) {
	if sa == nil {
		return
	}
	toResolve := FuncsAndParam{
		Param: "",
		Funcs: []string{},
	}
	if sa.ToResolve != nil {
		toResolve = *sa.ToResolve
	}
	kvps, found := (*mbr)[toResolve.Param]
	if !found {
		kvps = make([]*KVParsed, 0)
	}
	(*mbr)[toResolve.Param] = append(kvps, sa)
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
			Vals: (*a)[i].Vals.Ordered().asJson(),
		}
	}
	sort.Sort(&ret)
	return ret
}

func (mbkv *MapByKeyValues) add(key string, val *Value, orderRef *int) {
	unresolved := FuncsAndParam{
		Param: "",
		Funcs: []string{},
	}
	if val.Unresolved != nil {
		unresolved = *val.Unresolved
	}
	bkvs, found := (*mbkv)[unresolved.Param]
	if !found {
		bkvs = make([]*ByKeyValues, 0)
	}
	var bkv *ByKeyValues
	for ibkv := range bkvs {
		tbkv := bkvs[ibkv]
		if tbkv.Key == key {
			bkv = tbkv
			break
		}
	}
	if bkv == nil {
		bkv = &ByKeyValues{
			Key:  key,
			Vals: createValues(orderRef),
		}
		bkvs = append(bkvs, bkv)
	}
	bkv.Vals.getOrAddValue(&Value{
		Value:      val.Value,
		Unresolved: &unresolved,
		Tags:       val.Tags,
	})
	(*mbkv)[unresolved.Param] = bkvs
}

// Match KVPearls against the MapByToResolve
func (kvps *KVPearls) Match(toResolves MapByToResolve) MapByKeyValues {
	ret := MapByKeyValues{}
	kvp := kvps.Merge()

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
	for o := range kvp {
		// fmt.Printf("------- 0000:%d\n", o)
		oval := kvp[o]
		for i := range oval.Values {
			// fmt.Printf("------- 0000:%d:%s\n", o, i)
			// key := oval.Values[i]
			// rev := key.Values.Ordered()
			// fmt.Printf("------- 0000:%d:%s:%d\n", o, i, k)
			val := oval.Values[i]
			var unresolved *FuncsAndParam
			if len(toResolves) == 0 {
				ret.add(oval.Key, &Value{
					Value:      val.Value,
					Unresolved: unresolved,
					Tags:       val.Tags.asTags(),
					// order:      *kvps.orderRef,
				}, kvps.orderRef)
			}
			for ires := range toResolves {
				for ikvps := range toResolves[ires] {
					kvp := toResolves[ires][ikvps]
					matchKVP, match := kvp.Match(oval.Key, val)
					if match {
						if matchKVP.ToResolve != nil {
							unresolved = matchKVP.ToResolve
						}
						ret.add(oval.Key, &Value{
							Value:      val.Value,
							Unresolved: unresolved,
							Tags:       val.Tags.asTags(),
							// order:      *kvps.orderRef,
						}, kvps.orderRef)
					}
				}
			}
		}
	}
	return ret
}
