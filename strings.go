package k

import (
	"encoding/json"
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
