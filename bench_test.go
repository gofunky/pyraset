package mapset

import (
	"math/rand"
	"testing"
)

func nrand(n int) []int {
	i := make([]int,n)
	for ind := range i  { i[ind] =rand.Int()
	}
	return i
}

func toInterfaces(i []int) []interface{} {
	ifs := make([]interface{}, len(i))
	for ind, v := range i {
		ifs[ind] = v
	}
	return ifs
}

func benchAdd(b *testing.B, s Set) {
	nums := nrand(b.N)
	b.ResetTimer()
	for _, v := range nums {
		s.Add(v)
	}
}

func BenchmarkAddSafe(b *testing.B) {
	benchAdd(b, NewSet())
}

func BenchmarkAddUnsafe(b *testing.B) {
	benchAdd(b, NewUnsafeSet())
}

func benchRemove(b *testing.B, s Set) {
	nums := nrand(b.N)
	for _, v := range nums {
		s.Add(v)
	}

	b.ResetTimer()
	for _, v := range nums {
		s.Remove(v)
	}
}

func BenchmarkRemoveSafe(b *testing.B) {
	benchRemove(b, NewSet())
}

func BenchmarkRemoveUnsafe(b *testing.B) {
	benchRemove(b, NewUnsafeSet())
}

func benchCardinality(b *testing.B, s Set) {
	for i := 0; i < b.N; i++ {
		s.Cardinality()
	}
}

func BenchmarkCardinalitySafe(b *testing.B) {
	benchCardinality(b, NewSet())
}

func BenchmarkCardinalityUnsafe(b *testing.B) {
	benchCardinality(b, NewUnsafeSet())
}

func benchClear(b *testing.B, s Set) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Clear()
	}
}

func BenchmarkClearSafe(b *testing.B) {
	benchClear(b, NewSet())
}

func BenchmarkClearUnsafe(b *testing.B) {
	benchClear(b, NewUnsafeSet())
}

func benchPop(b *testing.B, s Set) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Pop()
	}
}

func BenchmarkPopSafe(b *testing.B) {
	benchPop(b, NewSet())
}

func BenchmarkPopUnsafe(b *testing.B) {
	benchPop(b, NewUnsafeSet())
}

func benchClone(b *testing.B, n int, s Set) {
	nums := toInterfaces(nrand(n))
	for _, v := range nums {
		s.Add(v)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Clone()
	}
}

func BenchmarkClone1Safe(b *testing.B) {
	benchClone(b, 1, NewSet())
}

func BenchmarkClone1Unsafe(b *testing.B) {
	benchClone(b, 1, NewUnsafeSet())
}

func BenchmarkClone10Safe(b *testing.B) {
	benchClone(b, 10, NewSet())
}

func BenchmarkClone10Unsafe(b *testing.B) {
	benchClone(b, 10, NewUnsafeSet())
}

func BenchmarkClone100Safe(b *testing.B) {
	benchClone(b, 100, NewSet())
}

func BenchmarkClone100Unsafe(b *testing.B) {
	benchClone(b, 100, NewUnsafeSet())
}

func benchContains(b *testing.B, n int, s Set) {
	nums := toInterfaces(nrand(n))
	for _, v := range nums {
		s.Add(v)
	}

	nums[n-1] = -1 // Definitely not in s

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Contains(nums...)
	}
}

func BenchmarkContains1Safe(b *testing.B) {
	benchContains(b, 1, NewSet())
}

func BenchmarkContains1Unsafe(b *testing.B) {
	benchContains(b, 1, NewUnsafeSet())
}

func BenchmarkContains10Safe(b *testing.B) {
	benchContains(b, 10, NewSet())
}

func BenchmarkContains10Unsafe(b *testing.B) {
	benchContains(b, 10, NewUnsafeSet())
}

func BenchmarkContains100Safe(b *testing.B) {
	benchContains(b, 100, NewSet())
}

func BenchmarkContains100Unsafe(b *testing.B) {
	benchContains(b, 100, NewUnsafeSet())
}

func benchEqual(b *testing.B, n int, s, t Set) {
	nums := nrand(n)
	for _, v := range nums {
		s.Add(v)
		t.Add(v)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Equal(t)
	}
}

func BenchmarkEqual1Safe(b *testing.B) {
	benchEqual(b, 1, NewSet(), NewSet())
}

func BenchmarkEqual1Unsafe(b *testing.B) {
	benchEqual(b, 1, NewUnsafeSet(), NewUnsafeSet())
}

func BenchmarkEqual10Safe(b *testing.B) {
	benchEqual(b, 10, NewSet(), NewSet())
}

func BenchmarkEqual10Unsafe(b *testing.B) {
	benchEqual(b, 10, NewUnsafeSet(), NewUnsafeSet())
}

func BenchmarkEqual100Safe(b *testing.B) {
	benchEqual(b, 100, NewSet(), NewSet())
}

func BenchmarkEqual100Unsafe(b *testing.B) {
	benchEqual(b, 100, NewUnsafeSet(), NewUnsafeSet())
}

func benchDifference(b *testing.B, n int, s, t Set) {
	nums := nrand(n)
	for _, v := range nums {
		s.Add(v)
	}
	for _, v := range nums[:n/2] {
		t.Add(v)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Difference(t)
	}
}

func benchIsSubset(b *testing.B, n int, s, t Set) {
	nums := nrand(n)
	for _, v := range nums {
		s.Add(v)
		t.Add(v)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.IsSubset(t)
	}
}

func BenchmarkIsSubset1Safe(b *testing.B) {
	benchIsSubset(b, 1, NewSet(), NewSet())
}

func BenchmarkIsSubset1Unsafe(b *testing.B) {
	benchIsSubset(b, 1, NewUnsafeSet(), NewUnsafeSet())
}

func BenchmarkIsSubset10Safe(b *testing.B) {
	benchIsSubset(b, 10, NewSet(), NewSet())
}

func BenchmarkIsSubset10Unsafe(b *testing.B) {
	benchIsSubset(b, 10, NewUnsafeSet(), NewUnsafeSet())
}

func BenchmarkIsSubset100Safe(b *testing.B) {
	benchIsSubset(b, 100, NewSet(), NewSet())
}

func BenchmarkIsSubset100Unsafe(b *testing.B) {
	benchIsSubset(b, 100, NewUnsafeSet(), NewUnsafeSet())
}

func benchIsSuperset(b *testing.B, n int, s, t Set) {
	nums := nrand(n)
	for _, v := range nums {
		s.Add(v)
		t.Add(v)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.IsSuperset(t)
	}
}

func BenchmarkIsSuperset1Safe(b *testing.B) {
	benchIsSuperset(b, 1, NewSet(), NewSet())
}

func BenchmarkIsSuperset1Unsafe(b *testing.B) {
	benchIsSuperset(b, 1, NewUnsafeSet(), NewUnsafeSet())
}

func BenchmarkIsSuperset10Safe(b *testing.B) {
	benchIsSuperset(b, 10, NewSet(), NewSet())
}

func BenchmarkIsSuperset10Unsafe(b *testing.B) {
	benchIsSuperset(b, 10, NewUnsafeSet(), NewUnsafeSet())
}

func BenchmarkIsSuperset100Safe(b *testing.B) {
	benchIsSuperset(b, 100, NewSet(), NewSet())
}

func BenchmarkIsSuperset100Unsafe(b *testing.B) {
	benchIsSuperset(b, 100, NewUnsafeSet(), NewUnsafeSet())
}

func benchIsProperSubset(b *testing.B, n int, s, t Set) {
	nums := nrand(n)
	for _, v := range nums {
		s.Add(v)
		t.Add(v)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.IsProperSubset(t)
	}
}

func BenchmarkIsProperSubset1Safe(b *testing.B) {
	benchIsProperSubset(b, 1, NewSet(), NewSet())
}

func BenchmarkIsProperSubset1Unsafe(b *testing.B) {
	benchIsProperSubset(b, 1, NewUnsafeSet(), NewUnsafeSet())
}

func BenchmarkIsProperSubset10Safe(b *testing.B) {
	benchIsProperSubset(b, 10, NewSet(), NewSet())
}

func BenchmarkIsProperSubset10Unsafe(b *testing.B) {
	benchIsProperSubset(b, 10, NewUnsafeSet(), NewUnsafeSet())
}

func BenchmarkIsProperSubset100Safe(b *testing.B) {
	benchIsProperSubset(b, 100, NewSet(), NewSet())
}

func BenchmarkIsProperSubset100Unsafe(b *testing.B) {
	benchIsProperSubset(b, 100, NewUnsafeSet(), NewUnsafeSet())
}

func benchIsProperSuperset(b *testing.B, n int, s, t Set) {
	nums := nrand(n)
	for _, v := range nums {
		s.Add(v)
		t.Add(v)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.IsProperSuperset(t)
	}
}

func BenchmarkIsProperSuperset1Safe(b *testing.B) {
	benchIsProperSuperset(b, 1, NewSet(), NewSet())
}

func BenchmarkIsProperSuperset1Unsafe(b *testing.B) {
	benchIsProperSuperset(b, 1, NewUnsafeSet(), NewUnsafeSet())
}

func BenchmarkIsProperSuperset10Safe(b *testing.B) {
	benchIsProperSuperset(b, 10, NewSet(), NewSet())
}

func BenchmarkIsProperSuperset10Unsafe(b *testing.B) {
	benchIsProperSuperset(b, 10, NewUnsafeSet(), NewUnsafeSet())
}

func BenchmarkIsProperSuperset100Safe(b *testing.B) {
	benchIsProperSuperset(b, 100, NewSet(), NewSet())
}

func BenchmarkIsProperSuperset100Unsafe(b *testing.B) {
	benchIsProperSuperset(b, 100, NewUnsafeSet(), NewUnsafeSet())
}

func BenchmarkDifference1Safe(b *testing.B) {
	benchDifference(b, 1, NewSet(), NewSet())
}

func BenchmarkDifference1Unsafe(b *testing.B) {
	benchDifference(b, 1, NewUnsafeSet(), NewUnsafeSet())
}

func BenchmarkDifference10Safe(b *testing.B) {
	benchDifference(b, 10, NewSet(), NewSet())
}

func BenchmarkDifference10Unsafe(b *testing.B) {
	benchDifference(b, 10, NewUnsafeSet(), NewUnsafeSet())
}

func BenchmarkDifference100Safe(b *testing.B) {
	benchDifference(b, 100, NewSet(), NewSet())
}

func BenchmarkDifference100Unsafe(b *testing.B) {
	benchDifference(b, 100, NewUnsafeSet(), NewUnsafeSet())
}

func benchIntersect(b *testing.B, n int, s, t Set) {
	nums := nrand(int(float64(n) * float64(1.5)))
	for _, v := range nums[:n] {
		s.Add(v)
	}
	for _, v := range nums[n/2:] {
		t.Add(v)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Intersect(t)
	}
}

func BenchmarkIntersect1Safe(b *testing.B) {
	benchIntersect(b, 1, NewSet(), NewSet())
}

func BenchmarkIntersect1Unsafe(b *testing.B) {
	benchIntersect(b, 1, NewUnsafeSet(), NewUnsafeSet())
}

func BenchmarkIntersect10Safe(b *testing.B) {
	benchIntersect(b, 10, NewSet(), NewSet())
}

func BenchmarkIntersect10Unsafe(b *testing.B) {
	benchIntersect(b, 10, NewUnsafeSet(), NewUnsafeSet())
}

func BenchmarkIntersect100Safe(b *testing.B) {
	benchIntersect(b, 100, NewSet(), NewSet())
}

func BenchmarkIntersect100Unsafe(b *testing.B) {
	benchIntersect(b, 100, NewUnsafeSet(), NewUnsafeSet())
}

func benchPower(b *testing.B, n int, s, t Set) {
	nums := nrand(int(float64(n) * float64(1.5)))
	for _, v := range nums[:n] {
		s.Add(v)
	}
	for _, v := range nums[n/2:] {
		t.Add(v)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.PowerSet()
	}
}

func BenchmarkPower1Safe(b *testing.B) {
	benchPower(b, 1, NewSet(), NewSet())
}

func BenchmarkPower1Unsafe(b *testing.B) {
	benchPower(b, 1, NewUnsafeSet(), NewUnsafeSet())
}

func BenchmarkPower10Safe(b *testing.B) {
	benchPower(b, 10, NewSet(), NewSet())
}

func BenchmarkPower10Unsafe(b *testing.B) {
	benchPower(b, 10, NewUnsafeSet(), NewUnsafeSet())
}

func BenchmarkPower20Safe(b *testing.B) {
	benchPower(b, 20, NewSet(), NewSet())
}

func BenchmarkPower20Unsafe(b *testing.B) {
	benchPower(b, 20, NewUnsafeSet(), NewUnsafeSet())
}

func benchCartesianProduct(b *testing.B, n int, s, t Set) {
	nums := nrand(int(float64(n) * float64(1.5)))
	for _, v := range nums[:n] {
		s.Add(v)
	}
	for _, v := range nums[n/2:] {
		t.Add(v)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.CartesianProduct(t)
	}
}

func BenchmarkCartesianProduct1Safe(b *testing.B) {
	benchCartesianProduct(b, 1, NewSet(), NewSet())
}

func BenchmarkCartesianProduct1Unsafe(b *testing.B) {
	benchCartesianProduct(b, 1, NewUnsafeSet(), NewUnsafeSet())
}

func BenchmarkCartesianProduct10Safe(b *testing.B) {
	benchCartesianProduct(b, 10, NewSet(), NewSet())
}

func BenchmarkCartesianProduct10Unsafe(b *testing.B) {
	benchCartesianProduct(b, 10, NewUnsafeSet(), NewUnsafeSet())
}

func BenchmarkCartesianProduct100Safe(b *testing.B) {
	benchCartesianProduct(b, 100, NewSet(), NewSet())
}

func BenchmarkCartesianProduct100Unsafe(b *testing.B) {
	benchCartesianProduct(b, 100, NewUnsafeSet(), NewUnsafeSet())
}

func benchSymmetricDifference(b *testing.B, n int, s, t Set) {
	nums := nrand(int(float64(n) * float64(1.5)))
	for _, v := range nums[:n] {
		s.Add(v)
	}
	for _, v := range nums[n/2:] {
		t.Add(v)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.SymmetricDifference(t)
	}
}

func BenchmarkSymmetricDifference1Safe(b *testing.B) {
	benchSymmetricDifference(b, 1, NewSet(), NewSet())
}

func BenchmarkSymmetricDifference1Unsafe(b *testing.B) {
	benchSymmetricDifference(b, 1, NewUnsafeSet(), NewUnsafeSet())
}

func BenchmarkSymmetricDifference10Safe(b *testing.B) {
	benchSymmetricDifference(b, 10, NewSet(), NewSet())
}

func BenchmarkSymmetricDifference10Unsafe(b *testing.B) {
	benchSymmetricDifference(b, 10, NewUnsafeSet(), NewUnsafeSet())
}

func BenchmarkSymmetricDifference100Safe(b *testing.B) {
	benchSymmetricDifference(b, 100, NewSet(), NewSet())
}

func BenchmarkSymmetricDifference100Unsafe(b *testing.B) {
	benchSymmetricDifference(b, 100, NewUnsafeSet(), NewUnsafeSet())
}

func benchUnion(b *testing.B, n int, s, t Set) {
	nums := nrand(n)
	for _, v := range nums[:n/2] {
		s.Add(v)
	}
	for _, v := range nums[n/2:] {
		t.Add(v)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Union(t)
	}
}

func BenchmarkUnion1Safe(b *testing.B) {
	benchUnion(b, 1, NewSet(), NewSet())
}

func BenchmarkUnion1Unsafe(b *testing.B) {
	benchUnion(b, 1, NewUnsafeSet(), NewUnsafeSet())
}

func BenchmarkUnion10Safe(b *testing.B) {
	benchUnion(b, 10, NewSet(), NewSet())
}

func BenchmarkUnion10Unsafe(b *testing.B) {
	benchUnion(b, 10, NewUnsafeSet(), NewUnsafeSet())
}

func BenchmarkUnion100Safe(b *testing.B) {
	benchUnion(b, 100, NewSet(), NewSet())
}

func BenchmarkUnion100Unsafe(b *testing.B) {
	benchUnion(b, 100, NewUnsafeSet(), NewUnsafeSet())
}

func benchEach(b *testing.B, n int, s Set) {
	nums := nrand(n)
	for _, v := range nums {
		s.Add(v)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Each(func(elem interface{}) bool {
			return false
		})
	}
}

func BenchmarkEach1Safe(b *testing.B) {
	benchEach(b, 1, NewSet())
}

func BenchmarkEach1Unsafe(b *testing.B) {
	benchEach(b, 1, NewUnsafeSet())
}

func BenchmarkEach10Safe(b *testing.B) {
	benchEach(b, 10, NewSet())
}

func BenchmarkEach10Unsafe(b *testing.B) {
	benchEach(b, 10, NewUnsafeSet())
}

func BenchmarkEach100Safe(b *testing.B) {
	benchEach(b, 100, NewSet())
}

func BenchmarkEach100Unsafe(b *testing.B) {
	benchEach(b, 100, NewUnsafeSet())
}

func benchIter(b *testing.B, n int, s Set) {
	nums := nrand(n)
	for _, v := range nums {
		s.Add(v)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c := s.Iter()
		for range c {

		}
	}
}

func BenchmarkIter1Safe(b *testing.B) {
	benchIter(b, 1, NewSet())
}

func BenchmarkIter1Unsafe(b *testing.B) {
	benchIter(b, 1, NewUnsafeSet())
}

func BenchmarkIter10Safe(b *testing.B) {
	benchIter(b, 10, NewSet())
}

func BenchmarkIter10Unsafe(b *testing.B) {
	benchIter(b, 10, NewUnsafeSet())
}

func BenchmarkIter100Safe(b *testing.B) {
	benchIter(b, 100, NewSet())
}

func BenchmarkIter100Unsafe(b *testing.B) {
	benchIter(b, 100, NewUnsafeSet())
}

func benchIterator(b *testing.B, n int, s Set) {
	nums := nrand(n)
	for _, v := range nums {
		s.Add(v)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c := s.Iterator().C
		for range c {

		}
	}
}

func BenchmarkIterator1Safe(b *testing.B) {
	benchIterator(b, 1, NewSet())
}

func BenchmarkIterator1Unsafe(b *testing.B) {
	benchIterator(b, 1, NewUnsafeSet())
}

func BenchmarkIterator10Safe(b *testing.B) {
	benchIterator(b, 10, NewSet())
}

func BenchmarkIterator10Unsafe(b *testing.B) {
	benchIterator(b, 10, NewUnsafeSet())
}

func BenchmarkIterator100Safe(b *testing.B) {
	benchIterator(b, 100, NewSet())
}

func BenchmarkIterator100Unsafe(b *testing.B) {
	benchIterator(b, 100, NewUnsafeSet())
}

func benchString(b *testing.B, n int, s Set) {
	nums := nrand(n)
	for _, v := range nums {
		s.Add(v)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = s.String()
	}
}

func BenchmarkString1Safe(b *testing.B) {
	benchString(b, 1, NewSet())
}

func BenchmarkString1Unsafe(b *testing.B) {
	benchString(b, 1, NewUnsafeSet())
}

func BenchmarkString10Safe(b *testing.B) {
	benchString(b, 10, NewSet())
}

func BenchmarkString10Unsafe(b *testing.B) {
	benchString(b, 10, NewUnsafeSet())
}

func BenchmarkString100Safe(b *testing.B) {
	benchString(b, 100, NewSet())
}

func BenchmarkString100Unsafe(b *testing.B) {
	benchString(b, 100, NewUnsafeSet())
}

func benchToSlice(b *testing.B, s Set) {
	nums := nrand(b.N)
	for _, v := range nums {
		s.Add(v)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.ToSlice()
	}
}

func BenchmarkToSliceSafe(b *testing.B) {
	benchToSlice(b, NewSet())
}

func BenchmarkToSliceUnsafe(b *testing.B) {
	benchToSlice(b, NewUnsafeSet())
}
