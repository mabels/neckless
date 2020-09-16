package kvpearl

import (
	"sort"
)

// Values structure with a order
type Values struct {
	orderRef *int
	values   map[string]*Value
}

func createValues(orderRef *int) Values {
	return Values{
		orderRef: orderRef,
		values:   map[string]*Value{},
	}
}

// type revOrderedValues []*Value

// Len is part of sort.Interface.
// func (s *revOrderedValues) Len() int {
// 	return len(*s)
// }

// // Swap is part of sort.Interface.
// func (s *revOrderedValues) Swap(i, j int) {
// 	(*s)[i], (*s)[j] = (*s)[j], (*s)[i]
// }

// // Less is part of sort.Interface. It is implemented by calling the "by" closure in the sorter.
// func (s *revOrderedValues) Less(i, j int) bool {
// 	return (*s)[i].Value > (*s)[j].Value
// }

// // RevOrdered give a list Values revervse by order
// func (values *Values) RevOrdered() *revOrderedValues {
// 	ret := make(revOrderedValues, len(values.Values))
// 	retIdx := 0
// 	for i := range values.Values {
// 		ret[retIdx] = values.Values[i]
// 		retIdx++
// 	}
// 	sort.Sort(&ret)
// 	return &ret
// }

// OrderedValues is the type which allow to handle tags in the right order
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

// Ordered convert the Values to OrderedValues
func (values *Values) Ordered() *OrderedValues {
	ret := make(OrderedValues, len(values.values))
	retIdx := 0
	for i := range values.values {
		ret[retIdx] = values.values[i]
		retIdx++
	}
	sort.Sort(&ret)
	// for i := range ret {
	// 	fmt.Printf("%s:%d\n", ret[i].Value, ret[i].order)
	// }
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

func (values *Values) getOrAddValue(val *Value) (*Value, bool) {
	value, found := values.values[val.Value]
	(*values.orderRef)++
	if !found {
		value = &Value{
			Value:      val.Value,
			Unresolved: val.Unresolved,
			Tags:       val.Tags,
			order:      *(values.orderRef),
		}
		values.values[val.Value] = value
	} else {
		value.Tags.add(val.Tags.toArray()...)
		value.order = *(values.orderRef)
	}
	return value, found
}

func (values *Values) getOrAdd(val string, inTags ...string) (*Value, bool) {
	tags := Tags{}
	for i := range inTags {
		tags[inTags[i]] = i
	}
	return values.getOrAddValue(&Value{
		Value:      val,
		Unresolved: nil,
		Tags:       tags,
	})
}

// JSONValues defines the JSON Respresantation of Values
type JSONValues []*JSONValue

// Value gets the actual defined value
func (jv *JSONValues) Value() *JSONValue {
	return (*jv)[len(*jv)-1]
}

func (s *OrderedValues) asJSON() JSONValues {
	ret := make(JSONValues, len(*s))
	for i := range *s {
		ret[i] = (*s)[i].asJSON()
	}
	return ret
}
