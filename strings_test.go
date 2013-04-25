package k

import (
	"math/rand"
	"regexp"
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

var (
	nonAlphaNum = regexp.MustCompile("^[^a-zA-Z0-9]+$")
)

func TestRandomString(t *testing.T) {
	for i := 0; i < 5; i++ {
		n := int(rand.Int31n(100))
		s := RandomString(n)
		if len(s) != n {
			t.Fatalf("Unexpected length. Expected %d, got %d", n, len(s))
		}

		if nonAlphaNum.MatchString(s) {
			t.Fatalf("String \"%s\" contains non-alpha-numeric characters", s)
		}
	}
}
