package util

import (
	"reflect"
	"testing"
)

func TestContainsString(t *testing.T) {
	if ContainsString("i like eggs", []string{"bacon"}, false) {
		t.Fatal("search string erroneously found")
	}
	if ContainsString("i like Bacon", []string{"bacon"}, true) {
		t.Fatal("case-sensitive search should have failed")
	}
	if !ContainsString("i like Bacon", []string{"bacon"}, false) {
		t.Fatal("case-insensitive search should have succeeded")
	}
}

func TestSplitAndTrimString(t *testing.T) {
	if !reflect.DeepEqual(SplitAndTrimString("a,b", ","), []string{"a", "b"}) {
		t.Fatal("items do not match")
	}
}
