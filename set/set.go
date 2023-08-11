// hash set
package set

import (
	"golang.org/x/exp/constraints"
	"golang.org/x/exp/maps"
)

type Set[T constraints.Ordered] map[T]struct{}

// Set is a container that stores unique values.
func New[T constraints.Ordered]() Set[T] {
	return make(Set[T])
}

// Len returns the length of the set.
func (s Set[T]) Len() int {
	return len(s)
}

func (s Set[T]) Slice() []T {
	return maps.Keys(s)
}

// DeepClone returns a deep copy of the set.
func (s Set[T]) Clone() Set[T] {
	return maps.Clone(s)
}

func (s Set[T]) Push(val ...T) Set[T] {
	for _, v := range val {
		s[v] = struct{}{}
	}
	return s
}

// Slice returns a slice of the set.
// If value is not in the set, return true.
// If value is in the set, return false.
func (s Set[T]) Insert(val T) bool {
	if s.Contain(val) {
		return false
	}
	s[val] = struct{}{}
	return true
}

// If value is not in the set, return true.
// If value is in the set, return false.
func (s Set[T]) Contain(val T) bool {
	_, ok := s[val]
	return ok
}

// Remove removes the value from the set.
// If value is in the set, return true.
// If value is not in the set, return false.
func (s Set[T]) Remove(val T) bool {
	if s.Contain(val) {
		delete(s, val)
		return true
	}
	return false
}

func (s Set[T]) IsEmpty() bool {
	return s.Len() == 0
}

func (s Set[T]) Clear() Set[T] {
	maps.Clear(s)
	return s
}

func (s Set[T]) CopyFrom(set Set[T]) Set[T] {
	maps.Copy(s, set)
	return s
}

// If value is in the set, return true and remove it.
func (s Set[T]) ContainAndRemove(val T) (ok bool) {
	return s.Remove(val)
}

func (s Set[T]) Range(f func(val T) (Continue bool)) {
	for k := range s {
		if !f(k) {
			break
		}
	}
}

func (s Set[T]) Equal(set Set[T]) (b bool) {
	return maps.Equal(s, set)
}

// If s is a subset of set, return true.
func (s Set[T]) IsSubSet(set Set[T]) (b bool) {
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
func (s Set[T]) IsSupperSet(set Set[T]) (b bool) {
	return !set.IsSubSet(s)
}

// Contains all elements in the original s but not in the set
func (s Set[T]) Difference(set Set[T]) (diff Set[T]) {
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
func (s Set[T]) SymmetricDifference(set Set[T]) (diff Set[T]) {
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
func (s Set[T]) Intersection(set Set[T]) (intersection Set[T]) {
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
func (s Set[T]) Union(set ...Set[T]) (union Set[T]) {
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
