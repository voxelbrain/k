package k

import (
	"encoding/json"
	"math/rand"
	"time"
)

// StringArray is an alias for an array of strings
// which has common operations defined as methods.
type StringArray []string

// Contains returns true if the receiver array contains
// an element equivalent to needle.
func (s StringArray) Contains(needle string) bool {
	for _, elem := range s {
		if elem == needle {
			return true
		}
	}
	return false
}

// String returns the JSON representation of the array.
func (s StringArray) String() string {
	a, _ := json.Marshal(s)
	return string(a)
}

var (
	alphaNum = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz1234567890"
	r        = rand.New(rand.NewSource(time.Now().UnixNano()))
)

// RandomString generates a random alpha-numeric string
// of length n
func RandomString(n int) string {
	s := make([]byte, 0, n)
	for i := 0; i < n; i++ {
		s = append(s, alphaNum[r.Int63n(int64(len(alphaNum)))])
	}
	return string(s)
}
