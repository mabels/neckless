package kvpearl

import (
	"fmt"
	"sort"
	"time"
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

type MergeArgs map[string]*KVParsed

func (kvps *KVPearls) Merge(keys MergeArgs) *KVPearl {
	kvp := Create()
	for i := range *kvps.orderByTime() {
		for j := range (*kvps)[i].Keys {
			key := (*kvps)[i].Keys[j]
			kvparsed, found := keys[key.Key]
			var unresolved *string
			if found {
				unresolved = kvparsed.Unresolved
			}
			if len(keys) == 0 || found {
				vals := key.Values.ordered()
				for k := range *vals {
					value := (*vals)[k]
					if len(value.Tags) == 0 {
						kvp.Set(SetArg{
							Key:        key.Key,
							Unresolved: unresolved,
							Val:        value.Value,
							Tags:       value.Tags.sorted(),
						})
					} else {
						found := false
						for t := range value.Tags {
							_, myFound := keys[t]
							found = found || myFound
						}
						fmt.Println(found, value.Tags, len(keys))
						if found {
							kvp.Set(SetArg{
								Key:        key.Key,
								Unresolved: unresolved,
								Val:        value.Value,
								Tags:       value.Tags.sorted(),
							})
						}
					}
				}
			}
		}
	}
	kvp.Created = time.Now()
	return kvp
}
