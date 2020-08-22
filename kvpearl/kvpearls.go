package kvpearl

import (
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

func (mas *MergeArgs) Match(key string, tags Tags) bool {
	if len(*mas) == 0 {
		return true
	} else {
		ma, found := (*mas)[key]
		if !found || len(ma.Tags) == 0 {
			return true
		}
		for tag := range ma.Tags {
			_, found = tags[tag]
			if found {
				return true
			}
		}
		return false
	}
}

func (kvps *KVPearls) Merge(keys MergeArgs) *KVPearl {
	kvp := Create()
	for i := range *kvps.orderByTime() {
		for j := range (*kvps)[i].Keys {
			key := (*kvps)[i].Keys[j]
			vals := key.Values.ordered()
			for k := range *vals {
				value := (*vals)[k]
				if keys.Match(key.Key, value.Tags) {
					kvp.Set(SetArg{
						Key:        key.Key,
						Unresolved: value.Unresolved,
						Val:        value.Value,
						Tags:       value.Tags.sorted(),
					})
				}
			}
		}
	}
	kvp.Created = time.Now()
	return kvp
}
