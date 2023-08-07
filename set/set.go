// hash set
package set

import (
	"github.com/zijiren233/gencontainer/restrictions"
)

type Set[T restrictions.Ordered] struct {
	m map[T]struct{}
}

type SetConf[T restrictions.Ordered] func(*Set[T])

// WithCap sets the capacity of the set.
// Need to be called before WithValues.
func WithCap[T restrictions.Ordered](cap int) SetConf[T] {
	return func(s *Set[T]) {
		if s.m == nil {
			s.m = make(map[T]struct{}, cap)
		}
	}
}

// WithValues sets the values of the set.
// Need to be called after WithCap.
func WithValues[T restrictions.Ordered](val ...T) SetConf[T] {
	return func(s *Set[T]) {
		if s.m == nil {
			s.m = make(map[T]struct{}, len(val))
		}
		for _, v := range val {
			s.m[v] = struct{}{}
		}
	}
}

// Set is a container that stores unique values.
func New[T restrictions.Ordered](conf ...SetConf[T]) *Set[T] {
	set := &Set[T]{}
	for _, c := range conf {
		c(set)
	}
	if set.m == nil {
		set.m = make(map[T]struct{})
	}
	return set
}

// Len returns the length of the set.
func (s *Set[T]) Len() int {
	return len(s.m)
}

func (s *Set[T]) Slice() []T {
	slice := make([]T, 0, s.Len())
	for k := range s.m {
		slice = append(slice, k)
	}
	return slice
}

// DeepClone returns a deep copy of the set.
func (s *Set[T]) Clone() *Set[T] {
	set := &Set[T]{}
	set.m = make(map[T]struct{}, s.Len())
	for k := range s.m {
		set.m[k] = struct{}{}
	}
	return set
}

func (s *Set[T]) Push(val ...T) {
	for _, v := range val {
		s.m[v] = struct{}{}
	}
}

// Slice returns a slice of the set.
// If value is not in the set, return true.
// If value is in the set, return false.
func (s *Set[T]) Insert(val T) bool {
	if s.Contain(val) {
		return false
	}
	s.m[val] = struct{}{}
	return true
}

// If value is not in the set, return true.
// If value is in the set, return false.
func (s *Set[T]) Contain(val T) bool {
	_, ok := s.m[val]
	return ok
}

// Remove removes the value from the set.
// If value is in the set, return true.
// If value is not in the set, return false.
func (s *Set[T]) Remove(val T) bool {
	if s.Contain(val) {
		delete(s.m, val)
		return true
	}
	return false
}

func (s *Set[T]) IsEmpty() bool {
	return s.Len() == 0
}

func (s *Set[T]) Clear() {
	s.m = make(map[T]struct{})
}

// If value is in the set, return true and remove it.
func (s *Set[T]) ContainAndRemove(val T) (ok bool) {
	return s.Remove(val)
}

func (s *Set[T]) Range(f func(val T) (Continue bool)) {
	for k := range s.m {
		if !f(k) {
			break
		}
	}
}

func (s *Set[T]) Equal(set *Set[T]) (b bool) {
	if s.Len() != set.Len() {
		return false
	}
	b = true
	s.Range(func(val T) bool {
		if !set.Contain(val) {
			b = false
			return false
		}
		return true
	})
	return
}

// If s is a subset of set, return true.
func (s *Set[T]) IsSubSet(set *Set[T]) (b bool) {
	if s.Len() > set.Len() {
		return false
	}
	b = true
	s.Range(func(val T) bool {
		if !set.Contain(val) {
			b = false
			return false
		}
		return true
	})
	return
}

// If s is a supper set of set, return true.
func (s *Set[T]) IsSupperSet(set *Set[T]) (b bool) {
	return set.IsSubSet(s)
}

// Contains all elements in the original s but not in the set
func (s *Set[T]) Difference(set *Set[T]) (diff *Set[T]) {
	diff = New[T]()
	s.Range(func(val T) bool {
		if !set.Contain(val) {
			diff.Insert(val)
		}
		return true
	})
	return
}

// Contains all elements in s or set but not both.
func (s *Set[T]) SymmetricDifference(set *Set[T]) (diff *Set[T]) {
	diff = New[T]()
	s.Range(func(val T) bool {
		if !set.Contain(val) {
			diff.Insert(val)
		}
		return true
	})
	set.Range(func(val T) bool {
		if !s.Contain(val) {
			diff.Insert(val)
		}
		return true
	})
	return
}

// Contains all elements in s and set.
func (s *Set[T]) Intersection(set *Set[T]) (intersection *Set[T]) {
	intersection = New[T]()
	s.Range(func(val T) bool {
		if set.Contain(val) {
			intersection.Insert(val)
		}
		return true
	})
	return
}

// Union all sets.
func (s *Set[T]) Union(set ...*Set[T]) (union *Set[T]) {
	union = New[T]()
	s.Range(func(val T) bool {
		union.Insert(val)
		return true
	})
	for _, v := range set {
		v.Range(func(val T) bool {
			union.Insert(val)
			return true
		})
	}
	return
}
