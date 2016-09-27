package util

import (
	"testing"
)

func TestContainsInt(t *testing.T) {
	if ContainsInt([]int{1, 2, 3}, 4) {
		t.Fatal("search should have failed")
	}
	if !ContainsInt([]int{1, 2, 3}, 1) {
		t.Fatal("search failed")
	}
}

func TestFilterInt(t *testing.T) {
	integers := FilterInt([]int{1, 2}, func(i int) bool { return i > 1 })
	if ContainsInt(integers, 1) {
		t.Fatal("value failing test erroneously included")
	}
	if !ContainsInt(integers, 2) {
		t.Fatal("value missing from slice")
	}
}
