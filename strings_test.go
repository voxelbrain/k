package k

import (
	"testing"
)

func TestStringArray_Contains(t *testing.T) {
	var uut StringArray
	uut = StringArray([]string{"a", "b", "c"})
	if !uut.Contains("a") {
		t.Fatalf("Array %s does contain \"a\", returned false", uut)
	}

	uut = StringArray([]string{"a", "b", "c"})
	if uut.Contains("x") {
		t.Fatalf("Array %s does not contain \"x\", returned true", uut)
	}
}
