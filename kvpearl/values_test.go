package kvpearl

import "testing"

func TestOrdered(t *testing.T) {
	values := createValues()
	values.getOrAdd("first")
	values.getOrAdd("second")
	values.getOrAdd("first")
	sorted := *values.Ordered()
	if sorted.Len() != 2 {
		t.Error("should be 2")
	}
	if sorted[0].Value != "first" {
		t.Error("should be first")
	}
	if sorted[1].Value != "second" {
		t.Error("should be second")
	}
}

func TestRevordered(t *testing.T) {
	values := createValues()
	values.getOrAdd("first")
	values.getOrAdd("second")
	values.getOrAdd("first")
	sorted := *values.RevOrdered()
	if sorted.Len() != 2 {
		t.Error("should be 2")
	}
	if sorted[1].Value != "first" {
		t.Error("should be first")
	}
	if sorted[0].Value != "second" {
		t.Error("should be second")
	}
}
