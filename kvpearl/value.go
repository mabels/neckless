package kvpearl

// Value structure
type Value struct {
	Value      string
	Unresolved *FuncsAndParam
	Tags       Tags
	order      int64
	// Tags       map[string]struct{}
}

// JSONValue the Value resprestation
type JSONValue struct {
	Value      string
	Unresolved *FuncsAndParam `json:"Unresolved,omitempty"`
	// Order      time.Time // `json:"-"`
	Tags []string
}

func (val *Value) asJSON() *JSONValue {
	return &JSONValue{
		Value:      val.Value,
		Unresolved: val.Unresolved,
		Tags:       val.Tags.sorted(),
	}
}
