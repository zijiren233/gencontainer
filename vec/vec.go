package vec

import (
	"sort"

	"github.com/zijiren233/gencontainer/restrictions"
)

var _ sort.Interface = (*Vec[int])(nil)

type Vec[T restrictions.Ordered] struct {
	data []T
}

type VecConf[T restrictions.Ordered] func(*Vec[T])

// Preallocates memory for the vector.
func WithCap[T restrictions.Ordered](c int) VecConf[T] {
	return func(v *Vec[T]) {
		if v.data == nil || c > v.Cap() {
			v.Resize(c)
		}
	}
}

func WithValues[T restrictions.Ordered](val ...T) VecConf[T] {
	return func(v *Vec[T]) {
		if v.data == nil {
			v.Resize(len(val))
		}
		v.data = append(v.data, val...)
	}
}

func New[T restrictions.Ordered](conf ...VecConf[T]) *Vec[T] {
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

// Insert into restrictions.Ordered position.
//
// if i < 0 or i out of range, panic.
func (v *Vec[T]) Insert(i int, val ...T) {
	if i > v.Len() || i < 0 {
		panic("index out of range")
	}
	v.data = append(v.data[:i], append(val, v.data[i:]...)...)
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

func (v *Vec[T]) FindFirst(val T) (index int, contain bool) {
	v.Range(func(i int, v T) bool {
		if v == val {
			index, contain = i, true
			return false
		}
		return true
	})
	return
}

func (v *Vec[T]) FindLast(val T) (i int, ok bool) {
	for i = v.Len() - 1; i >= 0; i-- {
		if v.data[i] == val {
			return i, true
		}
	}
	return
}

func (v *Vec[T]) FindAll(val T) (ret []int) {
	v.Range(func(i int, v T) bool {
		if v == val {
			ret = append(ret, i)
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

func (v *Vec[T]) Contain(val T) (contain bool) {
	_, contain = v.FindFirst(val)
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

// Interval returns a interval slice.
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
// If size < len, truncate and if len < cap/4, shrink.
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
		// truncate
		if size < c/4 {
			// shrink
			new := make([]T, size, c/4)
			copy(new, v.data)
			v.data = new
		} else {
			// not shrink
			v.data = v.data[:size]
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
		v.Remove(i)
	}
	return
}

func (v *Vec[T]) Less(i, j int) bool {
	return v.data[i] < v.data[j]
}
