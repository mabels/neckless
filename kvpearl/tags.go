package kvpearl

import (
	"sort"
	"strings"
)

// Tags uniq Tags list
type Tags map[string]int

func (tags *Tags) add(toAdd ...string) {
	for i := range toAdd {
		trimmed := strings.TrimSpace(toAdd[i])
		if len(trimmed) != 0 {
			(*tags)[trimmed] = len(*tags) + 1
		}
	}
}

func (tags *Tags) toArray() []string {
	ret := make([]string, len(*tags))
	retIdx := 0
	for i := range *tags {
		ret[retIdx] = i
		retIdx++
	}
	return ret
}

func (tags *Tags) sorted() []string {
	ret := tags.toArray()
	sort.Strings(ret)
	return ret
}

type tag struct {
	tag   string
	order int
}

type tagsArray []tag

// Len is part of sort.Interface.
func (s *tagsArray) Len() int {
	return len(*s)
}

// Swap is part of sort.Interface.
func (s *tagsArray) Swap(i, j int) {
	(*s)[i], (*s)[j] = (*s)[j], (*s)[i]
}

// Less is part of sort.Interface. It is implemented by calling the "by" closure in the sorter.
func (s *tagsArray) Less(i, j int) bool {
	return (*s)[i].order < (*s)[j].order
}

func (tags *Tags) byOrder() []string {
	toTags := make(tagsArray, len(*tags))
	toTagsIdx := 0
	for i := range *tags {
		order := (*tags)[i]
		toTags[toTagsIdx] = tag{
			tag:   i,
			order: order,
		}
		toTagsIdx++
	}
	sort.Sort(&toTags)
	ret := make([]string, len(*tags))
	for i := range toTags {
		ret[i] = toTags[i].tag
	}
	return ret
}

func tags2Map(tags []string) Tags {
	ret := Tags{}
	order := 0
	for i := range tags {
		tag := strings.TrimSpace(tags[i])
		if len(tag) > 0 {
			ret[tag] = order
			order++
		}
	}
	return ret
}

func tagstring2Map(strTags string) Tags {
	stripped := strings.TrimSpace(strTags)
	if len(stripped) > 0 {
		cs := strings.Split(stripped, ",")
		return tags2Map(cs)
	}
	return map[string]int{}
}
