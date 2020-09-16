package kvpearl

// Value structure
type Value struct {
	Value      string
	Unresolved *FuncsAndParam
	order      int
	Tags       Tags
}

// JSONTags is the Sorted Array Representation a tags
type JSONTags []string

// JSONValue the Value resprestation
type JSONValue struct {
	Value      string
	Unresolved *FuncsAndParam `json:"Unresolved,omitempty"`
	// Order      time.Time // `json:"-"`
	Tags JSONTags
}

func (tags *JSONTags) asTags() Tags {
	ret := Tags{}
	for i := range *tags {
		ret[(*tags)[i]] = i
	}
	return ret
}

func (val *Value) asJSON() *JSONValue {
	return &JSONValue{
		Value:      val.Value,
		Unresolved: val.Unresolved,
		Tags:       val.Tags.sorted(),
	}
}
