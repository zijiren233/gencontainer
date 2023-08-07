package vec

import (
	"golang.org/x/exp/constraints"
	"golang.org/x/exp/slices"
)

type Vec[T constraints.Ordered] struct {
	data []T
}

type VecConf[T constraints.Ordered] func(*Vec[T])

// Preallocates memory for the vector.
func WithCap[T constraints.Ordered](c int) VecConf[T] {
	return func(v *Vec[T]) {
		if v.data == nil || c > v.Cap() {
			v.Resize(c)
		}
	}
}

func WithValues[T constraints.Ordered](val ...T) VecConf[T] {
	return func(v *Vec[T]) {
		if v.data == nil {
			v.data = val
		} else {
			v.data = append(v.data, val...)
		}
	}
}

func New[T constraints.Ordered](conf ...VecConf[T]) *Vec[T] {
	vec := &Vec[T]{}
	for _, c := range conf {
		c(vec)
	}
	if vec.data == nil {
		vec.data = make([]T, 0)
	}
	return vec
}

func (v *Vec[T]) Len() int {
	return len(v.data)
}

// Push appends an element to the end.
func (v *Vec[T]) Push(val ...T) {
	v.data = append(v.data, val...)
}

// Pop removes the last element and returns it.
func (v *Vec[T]) Pop() (e T, ok bool) {
	if v.Len() == 0 {
		return
	}
	val := v.data[v.Len()-1]
	v.data = v.data[:v.Len()-1]
	return val, true
}

func (v *Vec[T]) Remove(i int) (e T, ok bool) {
	if i >= v.Len() || i < 0 {
		return
	}
	val := v.data[i]
	if i == v.Len()-1 {
		v.data = v.data[:i]
	} else {
		copy(v.data[i:], v.data[i+1:])
		v.data = v.data[:v.Len()-1]
	}
	return val, true
}

func (v *Vec[T]) Get(i int) (e T, ok bool) {
	if i >= v.Len() || i < 0 {
		return
	}
	return v.data[i], true
}

func (v *Vec[T]) Set(i int, val T) (ok bool) {
	if i >= v.Len() || i < 0 {
		return
	}
	v.data[i] = val
	return true
}

func (v *Vec[T]) Cap() int {
	return cap(v.data)
}

func (v *Vec[T]) IsEmpty() bool {
	return v.Len() == 0
}

// Clear the vector and shrink 1/4 the capacity.
func (v *Vec[T]) Clear() {
	v.Resize(0)
}

func (v *Vec[T]) FindAll(val T) (indexes []int) {
	v.Range(func(i int, v T) bool {
		if v == val {
			indexes = append(indexes, i)
		}
		return true
	})
	return
}

func (v *Vec[T]) Find(cbk func(v T) (matched bool)) (ret []int) {
	for i := 0; i < v.Len(); i++ {
		if cbk(v.data[i]) {
			ret = append(ret, i)
		}
	}
	return
}

func (v *Vec[T]) First() (e T, ok bool) {
	if v.Len() == 0 {
		return
	}
	return v.data[0], true
}

func (v *Vec[T]) Last() (e T, ok bool) {
	if v.Len() == 0 {
		return
	}
	return v.data[v.Len()-1], true
}

func (v *Vec[T]) Range(cbk func(i int, val T) (Continue bool)) {
	for i := 0; i < v.Len(); i++ {
		if !cbk(i, v.data[i]) {
			return
		}
	}
}

// Interval returns a interval vec.
func (v *Vec[T]) Interval(start, end int) (e []T, ok bool) {
	if start < 0 || end > v.Len() || start > end {
		return
	}
	return v.data[start:end], true
}

// Slice returns a slice of the vector.
func (v *Vec[T]) Slice() []T {
	return v.data[:v.Len()]
}

// If size < 0, do nothing.
//
// If size > len, realloc, increase cap, but not change len.
//
// If size < len, truncate and if len < cap/8, shrink.
func (v *Vec[T]) Resize(size int) {
	if size < 0 {
		return
	}
	if v.data == nil {
		v.data = make([]T, 0, size)
		return
	}
	l := v.Len()
	c := v.Cap()
	if size < l {
		v.data = v.data[:size]
		if size < c/8 {
			// shrink
			v.Clip()
		}
	} else if size > c {
		// realloc
		new := make([]T, l, size)
		copy(new, v.data)
		v.data = new
	}
}

// SplitOff splits the vector into two at the given index.
//
// if i < 0 or i out of range, return false.
//
// data dont overwrite each other.
func (v *Vec[T]) SplitOff(i int) (*Vec[T], bool) {
	if i < 0 || i > v.Len() {
		return nil, false
	}
	v2 := New[T](WithValues(v.data[i:]...))
	v.Resize(i)
	return v2, true
}

func (v *Vec[T]) Swap(i, j int) {
	if i < 0 || i >= v.Len() || j < 0 || j >= v.Len() {
		return
	}
	v.data[i], v.data[j] = v.data[j], v.data[i]
}

// ConpareAndRemove removes the element at index i if it is equal to val.
func (v *Vec[T]) ConpareAndRemove(i int, val T) (e T, ok bool) {
	if i < 0 || i >= v.Len() {
		return
	}
	if v.data[i] == val {
		return v.Remove(i)
	}
	return
}

func (v *Vec[T]) Less(i, j int) bool {
	return v.data[i] < v.data[j]
}

func (v *Vec[T]) Grow(n int) {
	v.data = slices.Grow(v.data, n)
}

// BinarySearch searches for target in a sorted slice and returns the position where target is found,
// or the position where target would appear in the sort order;
//
// it also returns a bool saying whether the target is really found in the slice.
//
// The slice must be sorted in increasing order.
func (v *Vec[T]) BinarySearch(target T) (int, bool) {
	return slices.BinarySearch(v.data, target)
}

func (v *Vec[T]) BinarySearchFunc(target any, cmp func(T, any) int) (int, bool) {
	return slices.BinarySearchFunc(v.data, target, cmp)
}

// Clip removes unused capacity
func (v *Vec[T]) Clip() {
	v.data = slices.Clip(v.data)
}

func (v *Vec[T]) Clone() *Vec[T] {
	return New[T](WithValues(slices.Clone(v.data)...))
}

func (v *Vec[T]) Compact() {
	v.data = slices.Compact(v.data)
}

func (v *Vec[T]) CompactFunc(f func(T, T) bool) {
	v.data = slices.CompactFunc(v.data, f)
}

func (v *Vec[T]) Compare(other *Vec[T]) int {
	return slices.Compare(v.data, other.data)
}

func (v *Vec[T]) CompareFunc(other *Vec[T], cmp func(T, T) int) int {
	return slices.CompareFunc(v.data, other.data, cmp)
}

func (v *Vec[T]) Index(val T) int {
	return slices.Index(v.data, val)
}

func (v *Vec[T]) Contains(val T) bool {
	return v.Index(val) >= 0
}

func (v *Vec[T]) Delete(i, j int) {
	v.data = slices.Delete(v.data, i, j)
}

func (v *Vec[T]) Equal(other *Vec[T]) bool {
	return slices.Equal(v.data, other.data)
}

// Insert into constraints.Ordered position.
//
// if i < 0 or i out of range, panic.
func (v *Vec[T]) Insert(i int, val ...T) {
	v.data = slices.Insert(v.data, i, val...)
}

func (v *Vec[T]) IsSorted() bool {
	return slices.IsSorted(v.data)
}

func (v *Vec[T]) Max() T {
	return slices.Max(v.data)
}

func (v *Vec[T]) Min() T {
	return slices.Min(v.data)
}

func (v *Vec[T]) Replace(i int, j int, val ...T) {
	v.data = slices.Replace(v.data, i, j, val...)
}

// Reverse reverses the elements of the vec in place.
func (v *Vec[T]) Reverse() {
	slices.Reverse(v.data)
}

// Sort sorts a vec of any ordered type in ascending order.
//
// When sorting floating-point numbers, NaNs are ordered before other values.
func (v *Vec[T]) Sort() {
	slices.Sort(v.data)
}

// SortFunc sorts the vec x in ascending order as determined by the cmp function. This sort is not guaranteed to be stable. cmp(a, b) should return a negative number when a < b, a positive number when a > b and zero when a == b.
//
// SortFunc requires that cmp is a strict weak ordering. See https://en.wikipedia.org/wiki/Weak_ordering#Strict_weak_orderings.
func (v *Vec[T]) SortFunc(cmp func(a T, b T) int) {
	slices.SortFunc(v.data, cmp)
}

// SortStableFunc sorts the vec x while keeping the original order of equal elements,
// using cmp to compare elements in the same way as SortFunc.
func (v *Vec[T]) SortStableFunc(cmp func(a T, b T) int) {
	slices.SortStableFunc(v.data, cmp)
}
