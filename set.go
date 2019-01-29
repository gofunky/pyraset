// mapset implements a simple and generic set collection. Items stored within it are unordered and unique. It supports
// typical set operations: membership testing, intersection, union, difference, symmetric difference and cloning.
//
// Package mapset provides two implementations of the Set interface. The default implementation is safe for concurrent
// access, but a non-thread-safe implementation is also provided for programs that can benefit from the slight
// complexity improvement and that can enforce mutual exclusion through other means.
package mapset

import "github.com/gofunky/hashstructure"

// Set is the primary interface provided by the mapset package.  It represents an unordered set of data and a
// large number of operations that can be applied to that set.
type Set interface {
	hashstructure.Hashable

	// Add the given elements to this set.
	Add(i ...interface{})

	// Cardinality determines the number of elements in the set.
	Cardinality() int

	// Clear removes all elements from the set, leaving the empty set.
	Clear()

	// Clone produces a clone of the set using the same implementation, duplicating all keys.
	Clone() Set

	// Contains determines whether the given items are all in the set.
	Contains(i ...interface{}) bool

	// Difference determines the difference between this set and the given set. The returned set will contain
	// all elements of this set that are not also elements of other.
	//
	// Note that the argument to Difference  must be of the same type as the receiver of the method.
	// Otherwise, Difference will panic.
	Difference(other Set) Set

	// Equal determines if two sets are equal to each other. If they have the same cardinality and contain the same
	// elements, they are considered equal. The order in which the elements were added is irrelevant.
	//
	// Note that the argument to Equal must be  of the same type as the receiver of the  method.
	// Otherwise, Equal will panic.
	Equal(other Set) bool

	// Returns a new set containing only the elements that exist only in both sets.
	//
	// Note that the argument to Intersect must be of the same type as the receiver of the method.
	// Otherwise, Intersect will panic.
	Intersect(other Set) Set

	// IsProperSubset determines if every element in this set is in the other set but the two sets are not equal.
	//
	// Note that the argument to IsProperSubset must be of the same type as the receiver of the method.
	// Otherwise, IsProperSubset will panic.
	IsProperSubset(other Set) bool

	// IsProperSuperset determines if every element in the other set is in this set but the two sets are notequal.
	//
	// Note that the argument to IsSuperset must be of the same type as the receiver of the method.
	// Otherwise, IsSuperset will panic.
	IsProperSuperset(other Set) bool

	// IsSubset determines if every element in this set is in the other set.
	//
	// Note that the argument to IsSubset  must be of the same type as the receiver of the method.
	// Otherwise, IsSubset will panic.
	IsSubset(other Set) bool

	// IsSuperset determines if every element in the other set is in this set.
	//
	// Note that the argument to IsSuperset  must be of the same type as the receiver of the method.
	// Otherwise, IsSuperset will panic.
	IsSuperset(other Set) bool

	// Each iterates over elements and executes the passed func against each element.
	// If passed func returns true, stop iteration eagerly.
	Each(func(interface{}) bool)

	// Iter returns a channel of elements that you can range over.
	Iter() <-chan interface{}

	// Iterator that you can use to range over the set.
	Iterator() *Iterator

	// Remove the given elements from this set.
	Remove(i ...interface{})

	// String provides a convenient string representation of the current state of the set.
	String() string

	// SymmetricDifference provides a new set with all elements which are  in either this set or the other set
	// but not in both.
	//
	// Note that the argument to SymmetricDifference must be of the same type as the receiver of the method.
	// Otherwise, SymmetricDifference will panic.
	SymmetricDifference(other Set) Set

	// Union provides a new set with all elements in this set and the given set.
	//
	// Note that the argument to Union must be of the same type as the receiver of the method.
	// Otherwise, IsSuperset will panic.
	Union(other Set) Set

	// Pop removes and returns an arbitrary item from the set.
	Pop() interface{}

	// PowerSet builds all subsets of a given set (Power Set).
	// By default, the sub sets derive their thread safety from the parent. Use threadSafe to override the default.
	PowerSet(threadSafe ...bool) Set

	// CartesianProduct builds the Cartesian Product of this set and the given set.
	CartesianProduct(other Set) Set

	// ToSlice converts the members of the set as to a slice.
	ToSlice() []interface{}

	// CoreSet provides the non-thread-safe core set.
	CoreSet() threadUnsafeSet

	// ThreadSafe provides the thread-safe core set.
	ThreadSafe() *threadSafeSet

	// MarshalJSON creates a JSON array from the set, it marshals all elements
	MarshalJSON() ([]byte, error)

	// UnmarshalJSON recreates a set from a JSON array, it only decodes primitive types.
	// Numbers are decoded as json.Number.
	UnmarshalJSON(p []byte) error
}

// NewSet creates a set that contains the given elements.
// Operations on the resulting set are thread-safe.
func NewSet(elements ...interface{}) Set {
	set := newThreadSafeSet()
	set.Add(elements...)
	return &set
}

// NewUnsafeSet creates a set that contains the given elements.
// Operations on the resulting set are not thread-safe.
func NewUnsafeSet(elements ...interface{}) Set {
	set := newThreadUnsafeSet()
	set.Add(elements...)
	return &set
}
