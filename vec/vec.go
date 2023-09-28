package vec

import (
	"reflect"

	"github.com/maruel/natural"
	"golang.org/x/exp/slices"
)

type Vec[T any] struct {
	d                 []T
	cmpLess, cmpEqual func(T, T) bool
}

type VecConf[T any] func(v *Vec[T])

func WithCmpLess[T any](cmpLess func(v1, v2 T) bool) VecConf[T] {
	return func(v *Vec[T]) {
		v.cmpLess = cmpLess
	}
}

func WithCmpEqual[T any](cmpEqual func(v1, v2 T) bool) VecConf[T] {
	return func(v *Vec[T]) {
		v.cmpEqual = cmpEqual
	}
}

func New[T any](conf ...VecConf[T]) *Vec[T] {
	v := &Vec[T]{}
	for _, c := range conf {
		c(v)
	}
	return v
}

func (v *Vec[T]) autoDetectSetCmpLess() *Vec[T] {
	if v.cmpLess != nil {
		return v
	}
	reflectType := reflect.TypeOf(v.d).Elem()
	if reflectType.Kind() == reflect.Ptr {
		reflectType = reflectType.Elem()
	}
	switch reflectType.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v.cmpLess = func(a, b T) bool {
			return reflect.ValueOf(a).Int() < reflect.ValueOf(b).Int()
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		v.cmpLess = func(a, b T) bool {
			return reflect.ValueOf(a).Uint() < reflect.ValueOf(b).Uint()
		}
	case reflect.Float32, reflect.Float64:
		v.cmpLess = func(a, b T) bool {
			return reflect.ValueOf(a).Float() < reflect.ValueOf(b).Float()
		}
	case reflect.String:
		v.cmpLess = func(a, b T) bool {
			return natural.Less(reflect.ValueOf(a).String(), reflect.ValueOf(b).String())
		}

	default:
		panic("warning: auto detect type failed, please set cmp less manually or use WithCmpLess when create vec")
	}
	return v
}

func (v *Vec[T]) autoDetectSetCmpEqual() *Vec[T] {
	if v.cmpEqual != nil {
		return v
	}
	reflectType := reflect.TypeOf(v.d).Elem()
	if reflectType.Kind() == reflect.Ptr {
		reflectType = reflectType.Elem()
	}
	switch reflectType.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v.cmpEqual = func(a, b T) bool {
			return reflect.ValueOf(a).Int() == reflect.ValueOf(b).Int()
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		v.cmpEqual = func(a, b T) bool {
			return reflect.ValueOf(a).Uint() == reflect.ValueOf(b).Uint()
		}
	case reflect.Float32, reflect.Float64:
		v.cmpEqual = func(a, b T) bool {
			return reflect.ValueOf(a).Float() == reflect.ValueOf(b).Float()
		}
	case reflect.String:
		v.cmpEqual = func(a, b T) bool {
			return reflect.ValueOf(a).String() == reflect.ValueOf(b).String()
		}
	default:
		panic("warning: auto detect type failed, please set cmp equal manually or use WithCmpEqual when create vec")
	}
	return v
}

func (v *Vec[T]) SetCmpLess(cmpLess func(T, T) bool) *Vec[T] {
	v.cmpLess = cmpLess
	return v
}

func (v *Vec[T]) SetCmpEqual(cmpEqual func(T, T) bool) *Vec[T] {
	v.cmpEqual = cmpEqual
	return v
}

func (v *Vec[T]) Len() int {
	return len(v.d)
}

// Push appends an element to the end.
func (v *Vec[T]) Push(val ...T) *Vec[T] {
	if v.Len()+len(val) > v.Cap() {
		v.Grow(len(val))
	}
	v.d = append(v.d, val...)
	return v
}

// Pop removes the last element and returns it.
func (v *Vec[T]) Pop() (e T, ok bool) {
	return v.Remove(v.Len() - 1)
}

func (v *Vec[T]) Remove(i int) (e T, ok bool) {
	l := v.Len()
	if i >= l || i < 0 {
		return
	}
	e = v.d[i]
	if i == l-1 {
		v.d = v.d[:i]
	} else {
		copy(v.d[i:], v.d[i+1:])
		v.d = v.d[:l-1]
	}
	return e, true
}

func (v *Vec[T]) Get(i int) (e T, ok bool) {
	if i >= v.Len() || i < 0 {
		return
	}
	return v.d[i], true
}

func (v *Vec[T]) Set(i int, val T) (ok bool) {
	if i >= v.Len() || i < 0 {
		return
	}
	v.d[i] = val
	return true
}

func (v *Vec[T]) Cap() int {
	return cap(v.d)
}

func (v *Vec[T]) IsEmpty() bool {
	return v.Len() == 0
}

// Clear the vector and shrink 1/4 the capacity.
func (v *Vec[T]) Clear() *Vec[T] {
	v.d = v.d[:0]
	return v
}

func (v *Vec[T]) Search(target T) (index int) {
	return v.SearchFunc(func(val T) (Matched bool) {
		return v.CmpEqual(target, val)
	})
}

func (v *Vec[T]) SearchFunc(cbk func(v T) (Matched bool)) (index int) {
	return slices.IndexFunc(v.d, func(val T) bool {
		return cbk(val)
	})
}

func (v *Vec[T]) SearchAll(target T) (indexs []int) {
	return v.SearchAllFunc(func(val T) (Matched bool) {
		return v.CmpEqual(target, val)
	})
}

func (v *Vec[T]) SearchAllFunc(cbk func(v T) (Matched bool)) (indexs []int) {
	v.Range(func(i int, val T) (Continue bool) {
		if cbk(val) {
			indexs = append(indexs, i)
		}
		return true
	})
	return
}

// BinarySearch searches for target in a sorted slice and returns the position where target is found,
// or the position where target would appear in the sort order;
//
// it also returns a bool saying whether the target is really found in the slice.
//
// The slice must be sorted in increasing order.
func (v *Vec[T]) BinarySearch(target T) (int, bool) {
	return slices.BinarySearchFunc(v.d, target, func(val, target T) int {
		if v.CmpLess(val, target) {
			return -1
		} else if v.CmpEqual(val, target) {
			return 0
		} else {
			return 1
		}
	})
}

func (v *Vec[T]) BinarySearchFunc(target T, cmp func(val T, target T) int) (int, bool) {
	return slices.BinarySearchFunc(v.d, target, cmp)
}

func (v *Vec[T]) First() (e T, ok bool) {
	if v.Len() == 0 {
		return
	}
	return v.d[0], true
}

func (v *Vec[T]) Last() (e T, ok bool) {
	if v.Len() == 0 {
		return
	}
	return v.d[v.Len()-1], true
}

func (v *Vec[T]) Range(cbk func(i int, val T) (Continue bool)) {
	for i := 0; i < v.Len(); i++ {
		if !cbk(i, v.d[i]) {
			return
		}
	}
}

// SliceCut returns a slice of the vector.
func (v *Vec[T]) SliceCut(i, j int) []T {
	if i < 0 || j > v.Len() || i > j {
		return nil
	}
	return v.d[i:j]
}

func (v *Vec[T]) Slice() []T {
	return v.d
}

// Cut removes the elements in the interval [i, j).
func (v *Vec[T]) Cut(i, j int) (e []T) {
	if i < 0 || j > v.Len() || i > j {
		return
	}
	val := v.d[i:j]
	v.d = append(v.d[:i], v.d[j:]...)
	return val
}

// SplitOff splits the vector into two at the given index.
//
// if i < 0 or i out of range, return false.
//
// data dont overwrite each other.
func (v *Vec[T]) SplitOff(i int) *Vec[T] {
	return New[T](WithCmpLess(v.CmpLess), WithCmpEqual(v.CmpEqual)).Push(v.Cut(i, v.Len())...)
}

func (v *Vec[T]) Swap(i, j int) {
	if i < 0 || i >= v.Len() || j < 0 || j >= v.Len() {
		return
	}
	v.d[i], v.d[j] = v.d[j], v.d[i]
}

// ConpareAndRemove removes the element at index i if it is equal to val.
func (v *Vec[T]) ConpareAndRemove(i int, val T) (e T, ok bool) {
	if i < 0 || i >= v.Len() {
		return
	}
	if v.CmpEqual(v.d[i], val) {
		return v.Remove(i)
	}
	return
}

func (v *Vec[T]) Less(i, j int) bool {
	return v.CmpLess(v.d[i], v.d[j])
}

func (v *Vec[T]) CmpLess(v1, v2 T) bool {
	return v.autoDetectSetCmpLess().cmpLess(v1, v2)
}

func (v *Vec[T]) Equal(i, j int) bool {
	return v.CmpEqual(v.d[i], v.d[j])
}

func (v *Vec[T]) CmpEqual(v1, v2 T) bool {
	return v.autoDetectSetCmpEqual().cmpEqual(v1, v2)
}

// After Grow(n), at least n elements can be appended to the slice without another allocation.
func (v *Vec[T]) Grow(n int) *Vec[T] {
	v.d = slices.Grow(v.d, n)
	return v
}

// Clip removes unused capacity
func (v *Vec[T]) Clip() *Vec[T] {
	v.d = slices.Clip(v.d)
	return v
}

func (v *Vec[T]) Clone() *Vec[T] {
	return New[T](WithCmpLess(v.CmpLess), WithCmpEqual(v.CmpEqual)).Push(v.d...)
}

func (v *Vec[T]) CloneToSlice() []T {
	return slices.Clone(v.d)
}

// Compact removes all identical adjacent elements.
// Only the first element of each group of equal elements is preserved.
// Must be sorted, Called after sorting.
func (v *Vec[T]) Compact() *Vec[T] {
	v.d = slices.CompactFunc(v.d, v.CmpEqual)
	return v
}

// Compact removes all identical adjacent elements.
// Only the first element of each group of equal elements is preserved.
// Must be sorted, Called after sorting.
func (v *Vec[T]) CompactFunc(eq func(T, T) bool) {
	v.d = slices.CompactFunc(v.d, eq)
}

// returns 0 the result is 0 if len(s1) == len(s2), -1 if len(s1) < len(s2),
// and +1 if len(s1) > len(s2).
func (v *Vec[T]) Compare(other Vec[T]) int {
	return v.CompareFunc(other, v.CmpLess, v.CmpEqual)
}

func (v *Vec[T]) CompareFunc(other Vec[T], cmpLess, cmpEqual func(T, T) bool) int {
	return slices.CompareFunc(v.d, other.d, func(v1, v2 T) int {
		if cmpLess(v1, v2) {
			return -1
		} else if cmpEqual(v1, v2) {
			return 0
		} else {
			return 1
		}
	})
}

func (v *Vec[T]) Contains(val T) bool {
	return slices.ContainsFunc(v.d, func(v1 T) bool {
		return v.CmpEqual(v1, val)
	})
}

func (v *Vec[T]) ContainsFunc(cmp func(v T) (matched bool)) bool {
	return slices.ContainsFunc(v.d, cmp)
}

func (v *Vec[T]) Delete(i, j int) {
	v.d = slices.Delete(v.d, i, j)
}

func (v *Vec[T]) EqualVec(other Vec[T]) bool {
	return slices.EqualFunc(v.d, other.d, v.CmpEqual)
}

func (v *Vec[T]) EqualSlice(other []T) bool {
	return slices.EqualFunc(v.d, other, v.CmpEqual)
}

func (v *Vec[T]) EqualVecFunc(other Vec[T], eq func(T, T) bool) bool {
	return slices.EqualFunc(v.d, other.d, eq)
}

func (v *Vec[T]) EqualSliceFunc(other []T, eq func(T, T) bool) bool {
	return slices.EqualFunc(v.d, other, eq)
}

// Insert into constraints.Ordered position.
//
// if i < 0 or i out of range, panic.
func (v *Vec[T]) Insert(i int, val ...T) *Vec[T] {
	v.d = slices.Insert(v.d, i, val...)
	return v
}

func (v *Vec[T]) IsSorted() bool {
	return v.IsSortedFunc(v.CmpLess, v.CmpEqual)
}

func (v *Vec[T]) IsSortedFunc(cmpLess, cmpEqual func(T, T) bool) bool {
	return slices.IsSortedFunc(v.d, func(h, l T) int {
		if cmpLess(h, l) {
			return -1
		} else if cmpEqual(h, l) {
			return 0
		} else {
			return 1
		}
	})
}

func (v *Vec[T]) Max() T {
	return v.MaxFunc(v.CmpLess, v.CmpEqual)
}

func (v *Vec[T]) MaxFunc(cmpLess, cmpEqual func(T, T) bool) T {
	return slices.MaxFunc(v.d, func(val, max T) int {
		if cmpLess(val, max) {
			return -1
		} else if cmpEqual(val, max) {
			return 0
		} else {
			return 1
		}
	})
}

func (v *Vec[T]) Min() T {
	return v.MinFunc(v.CmpLess, v.CmpEqual)
}

func (v *Vec[T]) MinFunc(cmpLess, cmpEqual func(T, T) bool) T {
	return slices.MinFunc(v.d, func(val, min T) int {
		if cmpLess(val, min) {
			return -1
		} else if cmpEqual(val, min) {
			return 0
		} else {
			return 1
		}
	})
}

func (v *Vec[T]) Replace(i int, j int, val ...T) {
	v.d = slices.Replace(v.d, i, j, val...)
}

// Reverse reverses the elements of the vec in place.
func (v *Vec[T]) Reverse() *Vec[T] {
	slices.Reverse(v.d)
	return v
}

// Sort sorts a vec of any ordered type in ascending order.
//
// When sorting floating-point numbers, NaNs are ordered before other values.
func (v *Vec[T]) Sort() *Vec[T] {
	return v.SortFunc(v.CmpLess, v.CmpEqual)
}

func (v *Vec[T]) SortFunc(cmpLess, cmpEqual func(T, T) bool) *Vec[T] {
	slices.SortFunc(v.d, func(a, v T) int {
		if cmpLess(a, v) {
			return -1
		} else if cmpEqual(a, v) {
			return 0
		} else {
			return 1
		}
	})
	return v
}

// SortStableFunc sorts the vec x while keeping the original order of equal elements,
// using cmp to compare elements in the same way as SortFunc.
func (v *Vec[T]) SortStable() *Vec[T] {
	return v.SortStableFunc(v.CmpLess, v.CmpEqual)
}

func (v *Vec[T]) SortStableFunc(cmpLess, cmpEqual func(T, T) bool) *Vec[T] {
	slices.SortStableFunc(v.d, func(a, v T) int {
		if cmpLess(a, v) {
			return -1
		} else if cmpEqual(a, v) {
			return 0
		} else {
			return 1
		}
	})
	return v
}
