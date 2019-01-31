package mapset

import "fmt"

// OrderedPair represents a 2-tuple of values.
type OrderedPair struct {
	First  interface{}
	Second interface{}
}

func (pair *OrderedPair) Equal(other OrderedPair) bool {
	if pair.First == other.First &&
		pair.Second == other.Second {
		return true
	}
	return false
}

// String outputs a 2-tuple in the form "(A, B)".
func (pair OrderedPair) String() string {
	return fmt.Sprintf("(%v, %v)", pair.First, pair.Second)
}
