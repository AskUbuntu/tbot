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

type stringTest struct {
	input  string
	output string
}

func TestTruncate(t *testing.T) {
	values := []stringTest{
		{"12", "12"},
		{"123", "123"},
		{"1234", "12…"},
	}
	for _, v := range values {
		o := Truncate(v.input, 3)
		if o != v.output {
			t.Fatalf("%s != %s", o, v.output)
		}
	}
}
