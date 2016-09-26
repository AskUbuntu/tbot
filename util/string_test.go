package util

import (
	"testing"
)

func TestContainsAny(t *testing.T) {
	if ContainsAny("bacon", []string{"i like eggs"}, false) {
		t.Fatal("search string erroneously found")
	}
	if ContainsAny("bacon", []string{"i like Bacon"}, true) {
		t.Fatal("case-sensitive search should have failed")
	}
	if !ContainsAny("bacon", []string{"i like Bacon"}, false) {
		t.Fatal("case-insensitive search should have succeeded")
	}
}
