package kvpearl

import (
	"sort"
)

// Values structure with a order
type Values struct {
	order  int64
	Values map[string]*Value
}

func createValues() Values {
	return Values{
		order:  0,
		Values: map[string]*Value{},
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

// RevOrdered give a list Values revervse by order
func (values *Values) RevOrdered() *revOrderedValues {
	ret := make(revOrderedValues, len(values.Values))
	retIdx := 0
	for i := range values.Values {
		ret[retIdx] = values.Values[i]
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

func (values *Values) Ordered() *OrderedValues {
	ret := make(OrderedValues, len(values.Values))
	retIdx := 0
	for i := range values.Values {
		ret[retIdx] = values.Values[i]
		retIdx++
	}
	sort.Sort(&ret)
	// for i := range ret {
	// 	fmt.Printf("%s:%d\n", ret[i].Value, ret[i].order)
	// }
	return &ret
}

func (values *Values) len() int {
	return len(values.Values)
}

func (values *Values) get(val string) *Value {
	value, found := values.Values[val]
	if found {
		return value
	}
	return nil
}

func (values *Values) getOrAddValue(val *Value) (*Value, bool) {
	value, found := values.Values[val.Value]
	if !found {
		values.order++
		value = &Value{
			Value:      val.Value,
			Unresolved: val.Unresolved,
			Tags:       val.Tags,
			order:      values.order,
		}
		values.Values[val.Value] = value
	}
	return value, found
}

func (values *Values) getOrAdd(val string) (*Value, bool) {
	return values.getOrAddValue(&Value{
		Value:      val,
		Unresolved: nil,
		Tags:       Tags{},
	})
}

type JsonValues []*JSONValue

func (ordered *revOrderedValues) asJson() JsonValues {
	ret := make(JsonValues, len(*ordered))
	for i := range *ordered {
		ret[i] = (*ordered)[i].asJSON()
	}
	return ret
}
