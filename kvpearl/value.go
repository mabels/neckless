package kvpearl

type Value struct {
	Value      string
	Unresolved *string
	Tags       Tags
	order      int64
	// Tags       map[string]struct{}
}

type JsonValue struct {
	Value      string
	Unresolved *string `json:"Unresolved,omitempty"`
	// Order      time.Time // `json:"-"`
	Tags []string
}

func (val *Value) asJson() *JsonValue {
	return &JsonValue{
		Value:      val.Value,
		Unresolved: val.Unresolved,
		Tags:       val.Tags.sorted(),
	}
}
