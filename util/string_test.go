package util

import (
	"testing"
)

func TestContainsString(t *testing.T) {
	if ContainsString([]string{"i like eggs"}, "bacon", false) {
		t.Fatal("search string erroneously found")
	}
	if ContainsString([]string{"i like Bacon"}, "bacon", true) {
		t.Fatal("case-sensitive search should have failed")
	}
	if !ContainsString([]string{"i like Bacon"}, "bacon", false) {
		t.Fatal("case-insensitive search should have succeeded")
	}
}
