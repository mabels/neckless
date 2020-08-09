package kvpearl

import (
	"sort"
	"strings"
)

type Tags map[string]int

func (tags *Tags) add(toAdd ...string) {
	for i := range toAdd {
		trimmed := strings.TrimSpace(toAdd[i])
		if len(trimmed) != 0 {
			(*tags)[trimmed] = len(*tags) + 1
		}
	}
}

func (tags *Tags) sorted() []string {
	ret := make([]string, len(*tags))
	retIdx := 0
	for i := range *tags {
		ret[retIdx] = i
		retIdx++
	}
	sort.Strings(ret)
	return ret
}

type tag struct {
	tag   string
	order int
}

type tags []tag

// Len is part of sort.Interface.
func (s *tags) Len() int {
	return len(*s)
}

// Swap is part of sort.Interface.
func (s *tags) Swap(i, j int) {
	(*s)[i], (*s)[j] = (*s)[j], (*s)[i]
}

// Less is part of sort.Interface. It is implemented by calling the "by" closure in the sorter.
func (s *tags) Less(i, j int) bool {
	return (*s)[i].order < (*s)[j].order
}

func (tgs *Tags) byOrder() []string {
	toTags := make(tags, len(*tgs))
	toTagsIdx := 0
	for i := range *tgs {
		order := (*tgs)[i]
		toTags[toTagsIdx] = tag{
			tag:   i,
			order: order,
		}
		toTagsIdx++
	}
	sort.Sort(&toTags)
	ret := make([]string, len(*tgs))
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
