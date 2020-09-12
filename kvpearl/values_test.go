package kvpearl

import (
	"testing"
)

func TestOrdered(t *testing.T) {
	for i := 0; i < 100; i++ {
		values := createValues(&i)
		values.getOrAdd("first")
		values.getOrAdd("second")
		values.getOrAdd("first")
		values.getOrAdd("first", "1")
		values.getOrAdd("second", "1")
		values.getOrAdd("four", "1")
		values.getOrAdd("first", "2")
		values.getOrAdd("second", "2")
		values.getOrAdd("five", "1")
		sorted := *values.Ordered()
		// js, _ := json.MarshalIndent(sorted, "", "  ")
		// t.Error(string(js))
		if sorted.Len() != 4 {
			t.Error("should be 4")
		}
		ref := sorted[0].order
		for i := range sorted {
			val := sorted[i].order
			if ref > val {
				t.Error("sort error")
			}
			ref = val
		}
	}
}

// func TestRevordered(t *testing.T) {
// 	for i := 0; i < 1; i++ {
// 		values := createValues()
// 		values.getOrAdd("first", "1")
// 		values.getOrAdd("second", "1")
// 		values.getOrAdd("four", "1")
// 		values.getOrAdd("first", "2")
// 		values.getOrAdd("second", "2")
// 		values.getOrAdd("five", "1")
// 		sorted := values.RevOrdered()
// 		if sorted.Len() != 4 {
// 			t.Error("should be 4")
// 		}
// 		js, _ := json.MarshalIndent(sorted, "", "  ")
// 		t.Error(string(js))
// 		// if sorted[1].Value != "first" {
// 		// 	t.Error("should be first")
// 		// }
// 		// if sorted[0].Value != "second" {
// 		// 	t.Error("should be second")
// 		// }
// 	}
// }
