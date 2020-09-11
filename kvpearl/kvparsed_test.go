package kvpearl

import "testing"

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
