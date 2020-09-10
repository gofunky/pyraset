package mapset

import "sync"

type threadSafeSet struct {
	threadUnsafeSet
	sync.RWMutex
}

func (o SetOptions) newThreadSafeSet() threadSafeSet {
	return threadSafeSet{threadUnsafeSet: o.newThreadUnsafeSet()}
}

func (set *threadSafeSet) Add(i ...interface{}) {
	set.Lock()
	defer set.Unlock()
	set.threadUnsafeSet.Add(i...)
}

func (set *threadSafeSet) Contains(i ...interface{}) bool {
	set.RLock()
	defer set.RUnlock()
	return set.threadUnsafeSet.Contains(i...)
}

func (set *threadSafeSet) IsSubset(other Set) bool {
	o := other.ThreadSafe()

	set.RLock()
	defer set.RUnlock()
	o.RLock()
	defer o.RUnlock()

	return set.threadUnsafeSet.IsSubset(&o.threadUnsafeSet)
}

func (set *threadSafeSet) IsProperSubset(other Set) bool {
	o := other.ThreadSafe()

	set.RLock()
	defer set.RUnlock()
	o.RLock()
	defer o.RUnlock()

	return set.threadUnsafeSet.IsProperSubset(&o.threadUnsafeSet)
}

func (set *threadSafeSet) IsSuperset(other Set) bool {
	return other.IsSubset(set)
}

func (set *threadSafeSet) IsProperSuperset(other Set) bool {
	return other.IsProperSubset(set)
}

func (set *threadSafeSet) Union(other Set) Set {
	o := other.ThreadSafe()

	set.RLock()
	defer set.RUnlock()
	o.RLock()
	defer o.RUnlock()

	return set.threadUnsafeSet.Union(&o.threadUnsafeSet).ThreadSafe()
}

func (set *threadSafeSet) Intersect(other Set) Set {
	o := other.ThreadSafe()

	set.RLock()
	defer set.RUnlock()
	o.RLock()
	defer o.RUnlock()

	return set.threadUnsafeSet.Intersect(&o.threadUnsafeSet).ThreadSafe()
}

func (set *threadSafeSet) Difference(other Set) Set {
	o := other.ThreadSafe()

	set.RLock()
	defer set.RUnlock()
	o.RLock()
	defer o.RUnlock()

	return set.threadUnsafeSet.Difference(&o.threadUnsafeSet).ThreadSafe()
}

func (set *threadSafeSet) SymmetricDifference(other Set) Set {
	o := other.ThreadSafe()

	set.RLock()
	defer set.RUnlock()
	o.RLock()
	defer o.RUnlock()

	return set.threadUnsafeSet.SymmetricDifference(&o.threadUnsafeSet).ThreadSafe()
}

func (set *threadSafeSet) Clear() {
	set.Lock()
	defer set.Unlock()
	set.threadUnsafeSet.Clear()
}

func (set *threadSafeSet) Remove(i ...interface{}) {
	set.Lock()
	defer set.Unlock()
	set.threadUnsafeSet.Remove(i...)
}

func (set *threadSafeSet) Cardinality() int {
	set.RLock()
	defer set.RUnlock()
	return set.threadUnsafeSet.Cardinality()
}

func (set *threadSafeSet) Empty() bool {
	set.RLock()
	defer set.RUnlock()
	return set.threadUnsafeSet.Empty()
}

func (set *threadSafeSet) Each(cb func(interface{}) bool) {
	set.RLock()
	defer set.RUnlock()
	for _, elem := range set.threadUnsafeSet.anyMap {
		if cb(elem) {
			return
		}
	}
}

func (set *threadSafeSet) Iter() <-chan interface{} {
	ch := make(chan interface{})
	go func() {
		set.RLock()
		defer set.RUnlock()
		defer close(ch)

		for _, elem := range set.threadUnsafeSet.anyMap {
			ch <- elem
		}
	}()

	return ch
}

func (set *threadSafeSet) Iterator() *Iterator {
	iterator, ch, stopCh := newIterator()

	go func() {
		set.RLock()
		defer set.RUnlock()
		defer close(ch)
	L:
		for _, elem := range set.threadUnsafeSet.anyMap {
			select {
			case <-stopCh:
				break L
			case ch <- elem:
			}
		}
	}()

	return iterator
}

func (set *threadSafeSet) Equal(other Set) bool {
	o := other.ThreadSafe()

	set.RLock()
	defer set.RUnlock()
	o.RLock()
	defer o.RUnlock()

	return set.threadUnsafeSet.Equal(&o.threadUnsafeSet)
}

func (set *threadSafeSet) Clone() Set {
	set.RLock()
	defer set.RUnlock()

	return set.threadUnsafeSet.Clone().ThreadSafe()
}

func (set *threadSafeSet) String() string {
	set.RLock()
	defer set.RUnlock()
	return set.threadUnsafeSet.String()
}

func (set *threadSafeSet) Hash() uint64 {
	return set.threadUnsafeSet.Hash()
}

func (set *threadSafeSet) UpdateHash() int {
	set.Lock()
	defer set.Unlock()
	return set.threadUnsafeSet.UpdateHash()
}

func (set *threadSafeSet) PowerSet() Set {
	set.RLock()
	defer set.RUnlock()
	return set.threadUnsafeSet.PowerSet().ThreadSafe()
}

func (set *threadSafeSet) Pop() interface{} {
	set.Lock()
	defer set.Unlock()
	return set.threadUnsafeSet.Pop()
}

func (set *threadSafeSet) CartesianProduct(other Set) Set {
	o := other.ThreadSafe()

	set.RLock()
	defer set.RUnlock()
	o.RLock()
	defer o.RUnlock()

	return set.threadUnsafeSet.CartesianProduct(&o.threadUnsafeSet).ThreadSafe()
}

func (set *threadSafeSet) ToSlice() []interface{} {
	set.RLock()
	defer set.RUnlock()
	return set.threadUnsafeSet.ToSlice()
}

func (set *threadSafeSet) threadUnsafeSetSet() threadUnsafeSet {
	return set.threadUnsafeSet
}

func (set *threadSafeSet) ThreadSafe() *threadSafeSet {
	return set
}

func (set *threadSafeSet) MarshalJSON() ([]byte, error) {
	set.RLock()
	defer set.RUnlock()
	return set.threadUnsafeSet.MarshalJSON()
}

func (set *threadSafeSet) UnmarshalJSON(p []byte) error {
	set.Lock()
	defer set.Unlock()
	return set.threadUnsafeSet.UnmarshalJSON(p)
}
