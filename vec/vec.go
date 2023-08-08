package vec

import (
	"golang.org/x/exp/constraints"
	"golang.org/x/exp/slices"
)

type Vec[T constraints.Ordered] []T

func New[T constraints.Ordered]() *Vec[T] {
	return &Vec[T]{}
}

func (v *Vec[T]) Len() int {
	return len(*v)
}

// Push appends an element to the end.
func (v *Vec[T]) Push(val ...T) *Vec[T] {
	if v.Len()+len(val) > v.Cap() {
		v.Grow(len(val))
	}
	(*v) = append(*v, val...)
	return v
}

// Pop removes the last element and returns it.
func (v *Vec[T]) Pop() (e T, ok bool) {
	return v.Remove(v.Len() - 1)
}

func (v *Vec[T]) Remove(i int) (e T, ok bool) {
	if i >= v.Len() || i < 0 {
		return
	}
	val := (*v)[i]
	if i == v.Len()-1 {
		(*v) = (*v)[:i]
	} else {
		copy((*v)[i:], (*v)[i+1:])
		(*v) = (*v)[:v.Len()-1]
	}
	return val, true
}

func (v *Vec[T]) Get(i int) (e T, ok bool) {
	if i >= v.Len() || i < 0 {
		return
	}
	return (*v)[i], true
}

func (v *Vec[T]) Set(i int, val T) (ok bool) {
	if i >= v.Len() || i < 0 {
		return
	}
	(*v)[i] = val
	return true
}

func (v *Vec[T]) Cap() int {
	return cap(*v)
}

func (v *Vec[T]) IsEmpty() bool {
	return v.Len() == 0
}

// Clear the vector and shrink 1/4 the capacity.
func (v *Vec[T]) Clear() *Vec[T] {
	(*v) = (*v)[:0]
	return v
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
		if cbk((*v)[i]) {
			ret = append(ret, i)
		}
	}
	return
}

func (v *Vec[T]) First() (e T, ok bool) {
	if v.Len() == 0 {
		return
	}
	return (*v)[0], true
}

func (v *Vec[T]) Last() (e T, ok bool) {
	if v.Len() == 0 {
		return
	}
	return (*v)[v.Len()-1], true
}

func (v *Vec[T]) Range(cbk func(i int, val T) (Continue bool)) {
	for i := 0; i < v.Len(); i++ {
		if !cbk(i, (*v)[i]) {
			return
		}
	}
}

// SliceCut returns a slice of the vector.
func (v *Vec[T]) SliceCut(i, j int) []T {
	if i < 0 || j > v.Len() || i > j {
		return nil
	}
	return (*v)[i:j]
}

func (v *Vec[T]) Slice() []T {
	return (*v)
}

// Cut removes the elements in the interval [i, j).
func (v *Vec[T]) Cut(i, j int) (e []T) {
	if i < 0 || j > v.Len() || i > j {
		return
	}
	val := (*v)[i:j]
	(*v) = append((*v)[:i], (*v)[j:]...)
	return val
}

// SplitOff splits the vector into two at the given index.
//
// if i < 0 or i out of range, return false.
//
// data dont overwrite each other.
func (v *Vec[T]) SplitOff(i int) *Vec[T] {
	return New[T]().Push(v.Cut(i, v.Len())...)
}

func (v *Vec[T]) Swap(i, j int) {
	if i < 0 || i >= v.Len() || j < 0 || j >= v.Len() {
		return
	}
	(*v)[i], (*v)[j] = (*v)[j], (*v)[i]
}

// ConpareAndRemove removes the element at index i if it is equal to val.
func (v *Vec[T]) ConpareAndRemove(i int, val T) (e T, ok bool) {
	if i < 0 || i >= v.Len() {
		return
	}
	if (*v)[i] == val {
		return v.Remove(i)
	}
	return
}

func (v *Vec[T]) Less(i, j int) bool {
	return (*v)[i] < (*v)[j]
}

// After Grow(n), at least n elements can be appended to the slice without another allocation.
func (v *Vec[T]) Grow(n int) *Vec[T] {
	(*v) = slices.Grow(*v, n)
	return v
}

// BinarySearch searches for target in a sorted slice and returns the position where target is found,
// or the position where target would appear in the sort order;
//
// it also returns a bool saying whether the target is really found in the slice.
//
// The slice must be sorted in increasing order.
func (v *Vec[T]) BinarySearch(target T) (int, bool) {
	return slices.BinarySearch(*v, target)
}

func (v *Vec[T]) BinarySearchFunc(target any, cmp func(T, any) int) (int, bool) {
	return slices.BinarySearchFunc(*v, target, cmp)
}

// Clip removes unused capacity
func (v *Vec[T]) Clip() *Vec[T] {
	(*v) = slices.Clip(*v)
	return v
}

func (v *Vec[T]) Clone() *Vec[T] {
	return New[T]().Push(v.Slice()...)
}

func (v *Vec[T]) CloneToSlice() []T {
	return slices.Clone(*v)
}

func (v *Vec[T]) Compact() *Vec[T] {
	(*v) = slices.Compact(*v)
	return v
}

func (v *Vec[T]) CompactFunc(f func(T, T) bool) {
	(*v) = slices.CompactFunc(*v, f)
}

func (v *Vec[T]) Compare(other Vec[T]) int {
	return slices.Compare(*v, other)
}

func (v *Vec[T]) CompareFunc(other Vec[T], cmp func(T, T) int) int {
	return slices.CompareFunc(*v, other, cmp)
}

func (v *Vec[T]) Index(val T) int {
	return slices.Index(*v, val)
}

func (v *Vec[T]) Contains(val T) bool {
	return v.Index(val) >= 0
}

func (v *Vec[T]) Delete(i, j int) {
	(*v) = slices.Delete(*v, i, j)
}

func (v *Vec[T]) Equal(other Vec[T]) bool {
	return slices.Equal(*v, other)
}

func (v *Vec[T]) EqualSlice(other []T) bool {
	return slices.Equal(*v, other)
}

// Insert into constraints.Ordered position.
//
// if i < 0 or i out of range, panic.
func (v *Vec[T]) Insert(i int, val ...T) {
	(*v) = slices.Insert(*v, i, val...)
}

func (v *Vec[T]) IsSorted() bool {
	return slices.IsSorted(*v)
}

func (v *Vec[T]) Max() T {
	return slices.Max(*v)
}

func (v *Vec[T]) Min() T {
	return slices.Min(*v)
}

func (v *Vec[T]) Replace(i int, j int, val ...T) {
	(*v) = slices.Replace(*v, i, j, val...)
}

// Reverse reverses the elements of the vec in place.
func (v *Vec[T]) Reverse() *Vec[T] {
	slices.Reverse(*v)
	return v
}

// Sort sorts a vec of any ordered type in ascending order.
//
// When sorting floating-point numbers, NaNs are ordered before other values.
func (v *Vec[T]) Sort() *Vec[T] {
	slices.Sort(*v)
	return v
}

// SortStableFunc sorts the vec x while keeping the original order of equal elements,
// using cmp to compare elements in the same way as SortFunc.
func (v *Vec[T]) SortStableFunc(cmp func(a T, b T) int) {
	slices.SortStableFunc(*v, cmp)
}
