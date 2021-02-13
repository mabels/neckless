package kvpearl

import (
	"testing"
)

func TestFuncsAndParamSimple(t *testing.T) {
	fap := ParseFuncsAndParams("")
	if fap.Param != "" {
		t.Error("should be ''")
	}
	if len(fap.Funcs) != 0 {
		t.Error("should be 0")
	}
}

func TestFuncsAndParamHelloSimple(t *testing.T) {
	fap := ParseFuncsAndParams("Hello")
	if fap.Param != "Hello" {
		t.Error("should be ''")
	}
	if len(fap.Funcs) != 0 {
		t.Error("should be 0")
	}
}

func TestFuncsAndParamOneAction(t *testing.T) {
	fap := ParseFuncsAndParams("Hello(World)")
	if fap.Param != "World" {
		t.Error("should be ''")
	}
	if !(len(fap.Funcs) == 1 && fap.Funcs[0] == "Hello") {
		t.Error("should be 0")
	}
}

func TestFuncsAndParamNestedAction(t *testing.T) {
	fap := ParseFuncsAndParams("Hello(Global(World))")
	if fap.Param != "World" {
		t.Error("should be ''")
	}
	if !(len(fap.Funcs) == 2 && fap.Funcs[0] == "Hello" && fap.Funcs[1] == "Global") {
		t.Error("should be 0")
	}
}

func TestMatchNoTags(t *testing.T) {
	kvp, err := Parse("XXX=YYY[]")
	if err != nil {
		t.Error(err)
	}
	if len(kvp.Tags) != 0 {
		t.Errorf("Tags should be empty")
	}
	_, matched := kvp.Match("XXX", &JSONValue{
		Tags: JSONTags{},
	})
	if !matched {
		t.Errorf("Should match XXX with no tags")
	}
	_, matched = kvp.Match("XXX", &JSONValue{
		Tags: JSONTags{"   ", ""},
	})
	if !matched {
		t.Errorf("Should match XXX with empty tags")
	}

	_, matched = kvp.Match("XXX", &JSONValue{
		Tags: JSONTags{"PROD"},
	})
	if matched {
		t.Errorf("Matched XXX with tags")
	}

	_, matched = kvp.Match("XXX", &JSONValue{
		Tags: JSONTags{"", "   ", "PROD"},
	})
	if !matched {
		t.Errorf("Matched XXX with tags")
	}
}

func TestMatchEmptyTags(t *testing.T) {
	kvp, err := Parse("XXX=YYY[,,,    ]")
	if err != nil {
		t.Error(err)
	}
	if len(kvp.Tags) != 1 {
		t.Errorf("Tags should be one")
	}
	_, matched := kvp.Match("XXX", &JSONValue{
		Tags: JSONTags{},
	})
	if !matched {
		t.Errorf("Should match XXX with no tags")
	}
	_, matched = kvp.Match("XXX", &JSONValue{
		Tags: JSONTags{"  ", ""},
	})
	if !matched {
		t.Errorf("Should match XXX with no tags")
	}
	_, matched = kvp.Match("XXX", &JSONValue{
		Tags: JSONTags{"PROD"},
	})
	if matched {
		t.Errorf("Matched XXX with tags")
	}

	_, matched = kvp.Match("XXX", &JSONValue{
		Tags: JSONTags{"", "PROD"},
	})
	if !matched {
		t.Errorf("Matched XXX with tags")
	}

	_, matched = kvp.Match("XXX", &JSONValue{
		Tags: JSONTags{"", "WURST"},
	})
	if !matched {
		t.Errorf("Matched XXX with tags")
	}
}

func TestMatchWithTags(t *testing.T) {
	kvp, err := Parse("XXX=YYY[PROD]")
	if err != nil {
		t.Error(err)
	}
	_, matched := kvp.Match("XXX", &JSONValue{
		Tags: JSONTags{},
	})
	if matched {
		t.Errorf("Match XXX with no tags")
	}
	_, matched = kvp.Match("XXX", &JSONValue{
		Tags: JSONTags{"PROD"},
	})
	if !matched {
		t.Errorf("Not Matched XXX with tags")
	}

	_, matched = kvp.Match("XXX", &JSONValue{
		Tags: JSONTags{"PROD", ""},
	})
	if !matched {
		t.Errorf("Not Matched XXX with tags")
	}

	_, matched = kvp.Match("XXX", &JSONValue{
		Tags: JSONTags{"XYY", ""},
	})
	if matched {
		t.Errorf("Not Matched XXX with tags")
	}
}

func TestMatchWithTagsEmpty(t *testing.T) {
	kvp, err := Parse("XXX=YYY[ ,PROD,]")
	if err != nil {
		t.Error(err)
	}
	if len(kvp.Tags) != 2 {
		t.Error("Tags should be len 2")
	}
	// jskvp, _ := json.Marshal(kvp)
	// t.Errorf(">>>>>>%s", jskvp)
	_, matched := kvp.Match("XXX", &JSONValue{
		Tags: JSONTags{},
	})
	if !matched {
		t.Errorf("Match XXX with no tags")
	}
	_, matched = kvp.Match("XXX", &JSONValue{
		Tags: JSONTags{" ", ""},
	})
	if !matched {
		t.Errorf("Match XXX with no tags")
	}
	_, matched = kvp.Match("XXX", &JSONValue{
		Tags: JSONTags{"PROD"},
	})
	if !matched {
		t.Errorf("Not Matched XXX with tags")
	}
	_, matched = kvp.Match("XXX", &JSONValue{
		Tags: JSONTags{"UNK"},
	})
	if matched {
		t.Errorf("Not Matched XXX with UNK")
	}
	_, matched = kvp.Match("XXX", &JSONValue{
		Tags: JSONTags{"UNK", ""},
	})
	if !matched {
		t.Errorf("Not Matched XXX with UNK")
	}
}
