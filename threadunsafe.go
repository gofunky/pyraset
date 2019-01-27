package mapset

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gofunky/hashstructure"
	"github.com/minio/highwayhash"
	"math/rand"
	"strings"
)

var (
	hashKey     = make([]byte, 32)
	hashOptions *hashstructure.HashOptions
)

func init() {
	_, err := rand.Read(hashKey)
	if err != nil {
		panic(err)
	}
	hasher, err := highwayhash.New64(hashKey)
	if err != nil {
		panic(err)
	}
	hashOptions = &hashstructure.HashOptions{
		Hasher: hasher,
	}
}

type threadUnsafeSet struct {
	hash   uint64
	anyMap map[uint64]interface{}
	subMap map[uint64]Set
}

// OrderedPair represents a 2-tuple of values.
type OrderedPair struct {
	First  interface{}
	Second interface{}
}

func newThreadUnsafeSet() threadUnsafeSet {
	return threadUnsafeSet{
		anyMap: make(map[uint64]interface{}),
		subMap: make(map[uint64]Set),
	}
}

func (pair *OrderedPair) Equal(other OrderedPair) bool {
	if pair.First == other.First &&
		pair.Second == other.Second {
		return true
	}
	return false
}

func (set *threadUnsafeSet) Add(i ...interface{}) {
	for _, val := range i {
		switch t := val.(type) {
		case Set:
			h := t.Hash()
			if _, ok := set.subMap[h]; !ok {
				set.hash = set.hash ^ h
				set.subMap[h] = t
			}
		default:
			h := hashFor(t)
			if _, ok := set.anyMap[h]; !ok {
				set.hash = set.hash ^ h
				set.anyMap[h] = t
			}
		}
	}
}

func (set *threadUnsafeSet) Contains(i ...interface{}) bool {
	for _, val := range i {
		switch t := val.(type) {
		case Set:
			if _, ok := set.subMap[t.Hash()]; !ok {
				return false
			}
		default:
			h := hashFor(t)
			if _, ok := set.anyMap[h]; !ok {
				return false
			}
		}
	}
	return true
}

func (set *threadUnsafeSet) IsSubset(other Set) bool {
	if set.Cardinality() > other.Cardinality() {
		return false
	}
	for _, elem := range set.subMap {
		if !other.Contains(elem) {
			return false
		}
	}
	for _, elem := range set.anyMap {
		if !other.Contains(elem) {
			return false
		}
	}
	return true
}

func (set *threadUnsafeSet) IsProperSubset(other Set) bool {
	return set.IsSubset(other) && !set.Equal(other)
}

func (set *threadUnsafeSet) IsSuperset(other Set) bool {
	return other.IsSubset(set)
}

func (set *threadUnsafeSet) IsProperSuperset(other Set) bool {
	return set.IsSuperset(other) && !set.Equal(other)
}

func (set *threadUnsafeSet) Union(other Set) Set {
	unionedSet := set.Clone()
	for _, elem := range other.CoreSet().subMap {
		unionedSet.Add(elem)
	}
	for _, elem := range other.CoreSet().anyMap {
		unionedSet.Add(elem)
	}
	return unionedSet
}

func (set *threadUnsafeSet) Intersect(other Set) Set {
	intersection := newThreadUnsafeSet()
	// loop over smaller set
	if set.Cardinality() < other.Cardinality() {
		for _, elem := range set.anyMap {
			if other.Contains(elem) {
				intersection.Add(elem)
			}
		}
		for _, elem := range set.subMap {
			if other.Contains(elem) {
				intersection.Add(elem)
			}
		}
	} else {
		for _, elem := range other.CoreSet().anyMap {
			if set.Contains(elem) {
				intersection.Add(elem)
			}
		}
		for _, elem := range other.CoreSet().subMap {
			if set.Contains(elem) {
				intersection.Add(elem)
			}
		}
	}
	return &intersection
}

func (set *threadUnsafeSet) Difference(other Set) Set {
	difference := newThreadUnsafeSet()
	for _, elem := range set.subMap {
		if !other.Contains(elem) {
			difference.Add(elem)
		}
	}
	for _, elem := range set.anyMap {
		if !other.Contains(elem) {
			difference.Add(elem)
		}
	}
	return &difference
}

func (set *threadUnsafeSet) SymmetricDifference(other Set) Set {
	aDiff := set.Difference(other)
	bDiff := other.Difference(set)
	return aDiff.Union(bDiff)
}

func (set *threadUnsafeSet) Clear() {
	*set = newThreadUnsafeSet()
}

func (set *threadUnsafeSet) Remove(i ...interface{}) {
	for _, val := range i {
		switch t := val.(type) {
		case Set:
			h := t.Hash()
			set.hash = set.hash ^ h
			delete(set.subMap, h)
		default:
			h := hashFor(t)
			set.hash = set.hash ^ h
			delete(set.anyMap, h)
		}
	}
}

func (set *threadUnsafeSet) Cardinality() int {
	return len(set.anyMap) + len(set.subMap)
}

func (set *threadUnsafeSet) Each(cb func(interface{}) bool) {
	for elem := range set.Iter() {
		if cb(elem) {
			break
		}
	}
}

func (set *threadUnsafeSet) Iter() <-chan interface{} {
	ch := make(chan interface{})
	go func() {
		for _, elem := range set.anyMap {
			ch <- elem
		}
		for _, elem := range set.subMap {
			ch <- elem
		}
		close(ch)
	}()

	return ch
}

func (set *threadUnsafeSet) Iterator() *Iterator {
	iterator, ch, stopCh := newIterator()

	go func() {
	L:
		for elem := range set.Iter() {
			select {
			case <-stopCh:
				break L
			case ch <- elem:
			}
		}
		close(ch)
	}()

	return iterator
}

func (set *threadUnsafeSet) Equal(other Set) bool {
	if set.Cardinality() != other.Cardinality() {
		return false
	}
	return set.Hash() == other.Hash()
}

func (set *threadUnsafeSet) Clone() Set {
	nextAny := make(map[uint64]interface{})
	nextSub := make(map[uint64]Set)
	for key, elem := range set.anyMap {
		nextAny[key] = elem
	}
	for key, elem := range set.subMap {
		nextSub[key] = elem
	}
	return &threadUnsafeSet{
		anyMap: nextAny,
		subMap: nextSub,
		hash:   set.hash,
	}
}

func (set *threadUnsafeSet) String() string {
	items := make([]string, 0, set.Cardinality())

	for elem := range set.Iter() {
		items = append(items, fmt.Sprintf("%v", elem))
	}
	return fmt.Sprintf("Set{%s}", strings.Join(items, ", "))
}

// String outputs a 2-tuple in the form "(A, B)".
func (pair OrderedPair) String() string {
	return fmt.Sprintf("(%v, %v)", pair.First, pair.Second)
}

func (set *threadUnsafeSet) Pop() interface{} {
	for item := range set.Iter() {
		set.Remove(item)
		return item
	}
	return nil
}

func (set *threadUnsafeSet) PowerSet(threadSafe ...bool) Set {
	var makeSafe = func(s Set) Set {
		return s.ThreadSafe()
	}
	if len(threadSafe) == 0 || !threadSafe[0] {
		makeSafe = func(s Set) Set {
			return s
		}
	}

	nullset := NewUnsafeSet()
	powSet := NewUnsafeSet(makeSafe(nullset))
	var interSlice = make([]Set, 0)
	interSlice = append(interSlice, nullset)

	for es := range set.Iter() {
		for _, is := range interSlice {
			newSubset := NewUnsafeSet(es).Union(is)
			interSlice = append(interSlice, newSubset)
			powSet.Add(makeSafe(newSubset))
		}
	}

	return powSet
}

func (set *threadUnsafeSet) CartesianProduct(other Set) Set {
	o := other.CoreSet()
	cartProduct := NewUnsafeSet()

	for i := range set.Iter() {
		for j := range o.Iter() {
			elem := OrderedPair{First: i, Second: j}
			cartProduct.Add(elem)
		}
	}

	return cartProduct
}

func (set *threadUnsafeSet) ToSlice() []interface{} {
	keys := make([]interface{}, set.Cardinality())
	var i = 0
	for _, elem := range set.subMap {
		keys[i] = elem
		i++
	}
	for _, elem := range set.anyMap {
		keys[i] = elem
		i++
	}

	return keys
}

func (set *threadUnsafeSet) CoreSet() threadUnsafeSet {
	return *set
}

func (set *threadUnsafeSet) ThreadSafe() *threadSafeSet {
	return &threadSafeSet{s: *set}
}

func (set *threadUnsafeSet) Hash() uint64 {
	return set.hash
}

func (set *threadUnsafeSet) MarshalJSON() ([]byte, error) {
	items := make([]string, 0, set.Cardinality())

	for elem := range set.Iter() {
		b, err := json.Marshal(elem)
		if err != nil {
			return nil, err
		}

		items = append(items, string(b))
	}

	return []byte(fmt.Sprintf("[%s]", strings.Join(items, ","))), nil
}

func (set *threadUnsafeSet) UnmarshalJSON(b []byte) error {
	var i []interface{}

	d := json.NewDecoder(bytes.NewReader(b))
	d.UseNumber()
	err := d.Decode(&i)
	if err != nil {
		return err
	}

	for _, v := range i {
		switch t := v.(type) {
		case []interface{}, map[string]interface{}:
			continue
		default:
			set.Add(t)
		}
	}

	return nil
}

func hashFor(i interface{}) uint64 {
	h, err := hashstructure.Hash(i, hashOptions)
	if err != nil {
		panic(err)
	}
	return h
}
