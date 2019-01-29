package mapset

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/OneOfOne/xxhash"
	"github.com/gofunky/hashstructure"
	"hash"
	"strings"
)

var (
	hasher      hash.Hash64
	hashOptions *hashstructure.HashOptions
)

func init() {
	hasher = xxhash.New64()
	hashOptions = &hashstructure.HashOptions{
		Hasher: hasher,
	}
}

type threadUnsafeSet struct {
	hash   uint64
	anyMap map[uint64]interface{}
}

// OrderedPair represents a 2-tuple of values.
type OrderedPair struct {
	First  interface{}
	Second interface{}
}

func newThreadUnsafeSet() threadUnsafeSet {
	return threadUnsafeSet{
		anyMap: make(map[uint64]interface{}),
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
		h := hashFor(val)
		set.addWithHash(val, h)
	}
}

func (set *threadUnsafeSet) addWithHash(val interface{}, h uint64) {
	if _, ok := set.anyMap[h]; !ok {
		set.hash = set.hash ^ h
		set.anyMap[h] = val
	}
}

func (set *threadUnsafeSet) Contains(i ...interface{}) bool {
	argLength := len(i)
	cardinality := set.Cardinality()
	if argLength > cardinality {
		return false
	}
	if argLength == cardinality {
		var inputHash uint64 = 0
		for _, val := range i {
			h := hashFor(val)
			inputHash = inputHash ^ h
		}
		return inputHash == set.hash
	}
	for _, val := range i {
		h := hashFor(val)
		if !set.containsHash(h) {
			return false
		}
	}
	return true
}

func (set *threadUnsafeSet) containsHash(h uint64) bool {
	if _, ok := set.anyMap[h]; !ok {
		return false
	}
	return true
}

func (set *threadUnsafeSet) IsSubset(other Set) bool {
	if set.Cardinality() > other.Cardinality() {
		return false
	}
	if set.Equal(other) {
		return true
	}
	otherCore := other.CoreSet()
	for h := range set.anyMap {
		if !otherCore.containsHash(h) {
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
	unionedSet := set.Clone().CoreSet()
	for h, elem := range other.CoreSet().anyMap {
		unionedSet.addWithHash(elem, h)
	}
	return &unionedSet
}

func (set *threadUnsafeSet) Intersect(other Set) Set {
	intersection := newThreadUnsafeSet()
	otherCore := other.CoreSet()
	// loop over smaller set
	if set.Cardinality() < other.Cardinality() {
		for h, elem := range set.anyMap {
			if otherCore.containsHash(h) {
				intersection.addWithHash(elem, h)
			}
		}
	} else {
		for h, elem := range otherCore.anyMap {
			if set.containsHash(h) {
				intersection.addWithHash(elem, h)
			}
		}
	}
	return &intersection
}

func (set *threadUnsafeSet) Difference(other Set) Set {
	difference := newThreadUnsafeSet()
	otherCore := other.CoreSet()
	for h, elem := range set.anyMap {
		if !otherCore.containsHash(h) {
			difference.addWithHash(elem, h)
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
		h := hashFor(val)
		set.hash = set.hash ^ h
		delete(set.anyMap, h)
	}
}

func (set *threadUnsafeSet) Cardinality() int {
	return len(set.anyMap)
}

func (set *threadUnsafeSet) Each(cb func(interface{}) bool) {
	for _, elem := range set.anyMap {
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
		close(ch)
	}()

	return ch
}

func (set *threadUnsafeSet) Iterator() *Iterator {
	iterator, ch, stopCh := newIterator()

	go func() {
	L:
		for _, elem := range set.anyMap {
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
	for key, elem := range set.anyMap {
		nextAny[key] = elem
	}
	return &threadUnsafeSet{
		anyMap: nextAny,
		hash:   set.hash,
	}
}

func (set *threadUnsafeSet) String() string {
	items := make([]string, 0, set.Cardinality())

	for _, elem := range set.anyMap {
		items = append(items, fmt.Sprintf("%v", elem))
	}
	return fmt.Sprintf("Set{%s}", strings.Join(items, ", "))
}

// String outputs a 2-tuple in the form "(A, B)".
func (pair OrderedPair) String() string {
	return fmt.Sprintf("(%v, %v)", pair.First, pair.Second)
}

func (set *threadUnsafeSet) Pop() interface{} {
	for _, item := range set.anyMap {
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
	powSet := NewUnsafeSet(makeSafe(nullset)).CoreSet()
	var interSlice = make([]Set, 0)
	interSlice = append(interSlice, nullset)

	for _, es := range set.anyMap {
		for _, is := range interSlice {
			newSubset := NewUnsafeSet(es).Union(is)
			interSlice = append(interSlice, newSubset)
			powSet.addWithHash(makeSafe(newSubset), newSubset.Hash())
		}
	}

	return &powSet
}

func (set *threadUnsafeSet) CartesianProduct(other Set) Set {
	o := other.CoreSet()
	cartProduct := NewUnsafeSet()

	for _, i := range set.anyMap {
		for _, j := range o.anyMap {
			elem := OrderedPair{First: i, Second: j}
			cartProduct.Add(elem)
		}
	}

	return cartProduct
}

func (set *threadUnsafeSet) ToSlice() []interface{} {
	keys := make([]interface{}, set.Cardinality())
	var i = 0
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
	return &threadSafeSet{threadUnsafeSet: *set}
}

func (set *threadUnsafeSet) Hash() uint64 {
	return set.hash
}

func (set *threadUnsafeSet) MarshalJSON() ([]byte, error) {
	items := make([]string, 0, set.Cardinality())

	for _, elem := range set.anyMap {
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
