package mapset

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/OneOfOne/xxhash"
	"github.com/gofunky/hashstructure"
	"strings"
)

type threadUnsafeSet struct {
	options     SetOptions
	hashState   uint64
	anyMap      map[uint64]interface{}
	hashCache   map[interface{}]uint64
	hashOptions *hashstructure.HashOptions
}

func (o SetOptions) newThreadUnsafeSet() threadUnsafeSet {
	hashOptions := &hashstructure.HashOptions{
		Hasher: o.Hasher,
	}
	if hashOptions.Hasher == nil {
		hashOptions.Hasher = xxhash.New64()
	}
	return threadUnsafeSet{
		options:     o,
		anyMap:      make(map[uint64]interface{}),
		hashCache:   make(map[interface{}]uint64),
		hashOptions: hashOptions,
	}
}

func (set *threadUnsafeSet) Add(i ...interface{}) {
	for _, val := range i {
		h := set.hashFor(val)
		set.addWithHash(val, h)
	}
}

func (set *threadUnsafeSet) addWithHash(val interface{}, h uint64) {
	if _, ok := set.anyMap[h]; !ok {
		set.hashState = set.hashState ^ h
		set.anyMap[h] = val
	}
}

func (set *threadUnsafeSet) Contains(i ...interface{}) bool {
	argLength := len(i)
	cardinality := set.Cardinality()
	if argLength > cardinality {
		return false
	}
	setHash := set.Hash()
	if argLength == cardinality {
		var inputHash uint64
		for _, val := range i {
			h := set.hashFor(val)
			inputHash = inputHash ^ h
		}
		return inputHash == setHash
	}
	for _, val := range i {
		h := set.hashFor(val)
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
	intersection := threadUnsafeSet{
		options:     set.options,
		anyMap:      make(map[uint64]interface{}),
		hashCache:   set.hashCache,
		hashOptions: set.hashOptions,
	}
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
	difference := threadUnsafeSet{
		options:     set.options,
		anyMap:      make(map[uint64]interface{}),
		hashCache:   set.hashCache,
		hashOptions: set.hashOptions,
	}
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
	*set = threadUnsafeSet{
		options:     set.options,
		anyMap:      make(map[uint64]interface{}),
		hashCache:   set.hashCache,
		hashOptions: set.hashOptions,
	}
}

func (set *threadUnsafeSet) Remove(i ...interface{}) {
	for _, val := range i {
		h := set.hashFor(val)
		set.removeWithHash(h)
	}
}

func (set *threadUnsafeSet) removeWithHash(hash uint64) {
	set.hashState = set.hashState ^ hash
	delete(set.anyMap, hash)
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
		options:     set.options,
		hashCache:   set.hashCache,
		anyMap:      nextAny,
		hashState:   set.hashState,
		hashOptions: set.hashOptions,
	}
}

func (set *threadUnsafeSet) String() string {
	items := bytes.NewBufferString("Set{")

	for _, elem := range set.anyMap {
		_, err := fmt.Fprintf(items, "%v, ", elem)
		if err != nil {
			panic(err)
		}
	}
	items.Truncate(items.Len() - 2)
	items.WriteString("}")
	return items.String()
}

func (set *threadUnsafeSet) Pop() interface{} {
	for h, item := range set.anyMap {
		set.removeWithHash(h)
		return item
	}
	return nil
}

func (set *threadUnsafeSet) PowerSet() Set {
	var makeSafe = func(s Set) Set {
		return s.ThreadSafe()
	}
	if set.options.Unsafe {
		makeSafe = func(s Set) Set {
			return s
		}
	}

	nullset := set.emptySet()
	powSet := set.emptySet()
	powSet.addWithHash(makeSafe(nullset), nullset.Hash())
	var interSlice = make([]Set, 0)
	interSlice = append(interSlice, nullset)

	for esHash, es := range set.anyMap {
		for _, is := range interSlice {
			newSubset := set.emptySet()
			newSubset.addWithHash(es, esHash)
			nextSubset := newSubset.Union(is)
			interSlice = append(interSlice, nextSubset)
			powSet.addWithHash(makeSafe(nextSubset), nextSubset.Hash())
		}
	}

	return powSet
}

func (set *threadUnsafeSet) CartesianProduct(other Set) Set {
	o := other.CoreSet()
	cartProduct := set.emptySet()

	for _, i := range set.anyMap {
		for _, j := range o.anyMap {
			elem := OrderedPair{First: i, Second: j}
			cartProduct.Add(elem)
		}
	}

	return cartProduct
}

func (set *threadUnsafeSet) ToSlice() (keys []interface{}) {
	keys = make([]interface{}, 0, set.Cardinality())
	for _, elem := range set.anyMap {
		keys = append(keys, elem)
	}
	return keys
}

func (set *threadUnsafeSet) CoreSet() threadUnsafeSet {
	return *set
}

func (set *threadUnsafeSet) ThreadSafe() *threadSafeSet {
	return &threadSafeSet{threadUnsafeSet: *set}
}

func (set threadUnsafeSet) Hash() uint64 {
	return set.hashState
}

func (set *threadUnsafeSet) UpdateHash() (updated int) {
	set.hashCache = make(map[interface{}]uint64)
	for hash, elem := range set.anyMap {
		h := set.hashFor(elem)
		if hash != h {
			set.removeWithHash(hash)
			set.addWithHash(elem, h)
			updated += 1
		}
	}
	return
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

func (set *threadUnsafeSet) emptySet() *threadUnsafeSet {
	return &threadUnsafeSet{
		options:     set.options,
		hashCache:   set.hashCache,
		anyMap:      make(map[uint64]interface{}),
		hashOptions: set.hashOptions,
	}
}

func (set *threadUnsafeSet) hashFor(i interface{}) uint64 {
	if set.options.Cache {
		if el, ok := set.hashCache[i]; ok {
			return el
		}
	}
	h, err := hashstructure.Hash(i, set.hashOptions)
	if err != nil {
		panic(err)
	}
	if set.options.Cache {
		set.hashCache[i] = h
	}
	return h
}
