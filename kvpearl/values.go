package kvpearl

import "sort"

type Values struct {
	order  int64
	values map[string]*Value
}

func createValues() Values {
	return Values{
		order:  0,
		values: map[string]*Value{},
	}
}

type revOrderedValues []*Value

// Len is part of sort.Interface.
func (s *revOrderedValues) Len() int {
	return len(*s)
}

// Swap is part of sort.Interface.
func (s *revOrderedValues) Swap(i, j int) {
	(*s)[i], (*s)[j] = (*s)[j], (*s)[i]
}

// Less is part of sort.Interface. It is implemented by calling the "by" closure in the sorter.
func (s *revOrderedValues) Less(i, j int) bool {
	return (*s)[i].order > (*s)[j].order
}

func (values *Values) revOrdered() *revOrderedValues {
	ret := make(revOrderedValues, len(values.values))
	retIdx := 0
	for i := range values.values {
		ret[retIdx] = values.values[i]
		retIdx++
	}
	sort.Sort(&ret)
	return &ret
}

type OrderedValues []*Value

// Len is part of sort.Interface.
func (s *OrderedValues) Len() int {
	return len(*s)
}

// Swap is part of sort.Interface.
func (s *OrderedValues) Swap(i, j int) {
	(*s)[i], (*s)[j] = (*s)[j], (*s)[i]
}

// Less is part of sort.Interface. It is implemented by calling the "by" closure in the sorter.
func (s *OrderedValues) Less(i, j int) bool {
	return (*s)[i].order < (*s)[j].order
}

func (values *Values) ordered() *OrderedValues {
	ret := make(OrderedValues, len(values.values))
	retIdx := 0
	for i := range values.values {
		ret[retIdx] = values.values[i]
		retIdx++
	}
	sort.Sort(&ret)
	return &ret
}

func (values *Values) len() int {
	return len(values.values)
}

func (values *Values) get(val string) *Value {
	value, found := values.values[val]
	if found {
		return value
	}
	return nil
}

func (values *Values) getOrAdd(val string) (*Value, bool) {
	value, found := values.values[val]
	if !found {
		values.order++
		value = &Value{
			Value:      val,
			Unresolved: nil,
			Tags:       Tags{},
			order:      values.order,
		}
		values.values[val] = value
	}
	return value, found
}

type JsonValues []*JsonValue

func (ordered *revOrderedValues) asJson() JsonValues {
	ret := make(JsonValues, len(*ordered))
	for i := range *ordered {
		ret[i] = (*ordered)[i].asJson()
	}
	return ret
}
