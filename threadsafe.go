package mapset

import "sync"

type threadSafeSet struct {
	s threadUnsafeSet
	sync.RWMutex
}

func newThreadSafeSet() threadSafeSet {
	return threadSafeSet{s: newThreadUnsafeSet()}
}

func (set *threadSafeSet) Add(i ...interface{}) {
	set.Lock()
	defer set.Unlock()
	set.s.Add(i...)
}

func (set *threadSafeSet) Contains(i ...interface{}) bool {
	set.RLock()
	defer set.RUnlock()
	return set.s.Contains(i...)
}

func (set *threadSafeSet) IsSubset(other Set) bool {
	o := other.ThreadSafe()

	set.RLock()
	defer set.RUnlock()
	o.RLock()
	defer o.RUnlock()

	return set.s.IsSubset(&o.s)
}

func (set *threadSafeSet) IsProperSubset(other Set) bool {
	o := other.ThreadSafe()

	set.RLock()
	defer set.RUnlock()
	o.RLock()
	defer o.RUnlock()

	return set.s.IsProperSubset(&o.s)
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

	return set.s.Union(&o.s).ThreadSafe()
}

func (set *threadSafeSet) Intersect(other Set) Set {
	o := other.ThreadSafe()

	set.RLock()
	defer set.RUnlock()
	o.RLock()
	defer o.RUnlock()

	return set.s.Intersect(&o.s).ThreadSafe()
}

func (set *threadSafeSet) Difference(other Set) Set {
	o := other.ThreadSafe()

	set.RLock()
	defer set.RUnlock()
	o.RLock()
	defer o.RUnlock()

	return set.s.Difference(&o.s).ThreadSafe()
}

func (set *threadSafeSet) SymmetricDifference(other Set) Set {
	o := other.ThreadSafe()

	set.RLock()
	defer set.RUnlock()
	o.RLock()
	defer o.RUnlock()

	return set.s.SymmetricDifference(&o.s).ThreadSafe()
}

func (set *threadSafeSet) Clear() {
	set.Lock()
	defer set.Unlock()
	set.s = newThreadUnsafeSet()
}

func (set *threadSafeSet) Remove(i ...interface{}) {
	set.Lock()
	defer set.Unlock()
	set.s.Remove(i...)
}

func (set *threadSafeSet) Cardinality() int {
	set.RLock()
	defer set.RUnlock()
	return set.s.Cardinality()
}

func (set *threadSafeSet) Each(cb func(interface{}) bool) {
	set.RLock()
	defer set.RUnlock()
	for elem := range set.s.Iter() {
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

		for _, elem := range set.s.anyMap {
			ch <- elem
		}
		for _, elem := range set.s.subMap {
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
		for elem := range set.s.Iter() {
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

	return set.s.Equal(&o.s)
}

func (set *threadSafeSet) Clone() Set {
	set.RLock()
	defer set.RUnlock()

	return set.s.Clone().ThreadSafe()
}

func (set *threadSafeSet) String() string {
	set.RLock()
	defer set.RUnlock()
	return set.s.String()
}

func (set *threadSafeSet) Hash() uint64 {
	set.RLock()
	defer set.RUnlock()
	return set.s.Hash()
}

func (set *threadSafeSet) PowerSet(threadSafe ...bool) Set {
	set.RLock()
	defer set.RUnlock()
	var safe = true
	if len(threadSafe) > 0 && !threadSafe[0] {
		safe = false
	}
	return set.s.PowerSet(safe).ThreadSafe()
}

func (set *threadSafeSet) Pop() interface{} {
	set.Lock()
	defer set.Unlock()
	return set.s.Pop()
}

func (set *threadSafeSet) CartesianProduct(other Set) Set {
	o := other.ThreadSafe()

	set.RLock()
	defer set.RUnlock()
	o.RLock()
	defer o.RUnlock()

	return set.s.CartesianProduct(&o.s).ThreadSafe()
}

func (set *threadSafeSet) ToSlice() []interface{} {
	set.RLock()
	defer set.RUnlock()
	return set.s.ToSlice()
}

func (set *threadSafeSet) CoreSet() threadUnsafeSet {
	return set.s
}

func (set *threadSafeSet) ThreadSafe() *threadSafeSet {
	return set
}

func (set *threadSafeSet) MarshalJSON() ([]byte, error) {
	set.RLock()
	defer set.RUnlock()
	return set.s.MarshalJSON()
}

func (set *threadSafeSet) UnmarshalJSON(p []byte) error {
	set.Lock()
	defer set.Unlock()
	return set.s.UnmarshalJSON(p)
}
