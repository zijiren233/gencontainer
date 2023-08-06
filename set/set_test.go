package set

import "testing"

func TestNew(t *testing.T) {
	s := New[int](WithCap[int](1024), WithValues(1, 2, 3, 4, 5, 4, 3, 2, 1))
	if s == nil {
		t.Errorf("New[int]() = %p, want not nil", s)
	}
	if s.Len() != 5 {
		t.Errorf("s.Len() = %d, want %d", s.Len(), 5)
	}
	if !s.Contain(5) {
		t.Errorf("s.Contain(5) = %t, want %t", s.Contain(5), true)
	}
}

func TestInsert(t *testing.T) {
	s := New[int]()
	for i := 0; i < 10; i++ {
		s.Insert(i)
	}
	if s.Len() != 10 {
		t.Errorf("s.Len() = %d, want %d", s.Len(), 10)
	}
	if !s.Contain(5) {
		t.Errorf("s.Contain(5) = %t, want %t", s.Contain(5), true)
	}
}

func TestRemove(t *testing.T) {
	s := New[int]()
	for i := 0; i < 10; i++ {
		s.Insert(i)
	}
	s.Remove(5)
	if s.Len() != 9 {
		t.Errorf("s.Len() = %d, want %d", s.Len(), 9)
	}
	if s.Contain(5) {
		t.Errorf("s.Contain(5) = %t, want %t", s.Contain(5), false)
	}
}

func TestClear(t *testing.T) {
	s := New[int]()
	for i := 0; i < 10; i++ {
		s.Insert(i)
	}
	s.Clear()
	if s.Len() != 0 {
		t.Errorf("s.Len() = %d, want %d", s.Len(), 0)
	}
	if s.Contain(5) {
		t.Errorf("s.Contain(5) = %t, want %t", s.Contain(5), false)
	}
}

func TestContainAndRemove(t *testing.T) {
	s := New[int]()
	for i := 0; i < 10; i++ {
		s.Insert(i)
	}
	if !s.ContainAndRemove(5) {
		t.Errorf("s.Contain(5) = %t, want %t", s.Contain(5), true)
	}
	if s.Contain(5) {
		t.Errorf("s.Contain(5) = %t, want %t", s.Contain(5), false)
	}
}

func TestEqual(t *testing.T) {
	s1 := New[int]()
	s2 := New[int]()
	for i := 0; i < 10; i++ {
		s1.Insert(i)
		s2.Insert(i)
	}
	if !s1.Equal(s2) {
		t.Errorf("s1.Equal(s2) = %t, want %t", s1.Equal(s2), true)
	}
	s2.Remove(5)
	if s1.Equal(s2) {
		t.Errorf("s1.Equal(s2) = %t, want %t", s1.Equal(s2), false)
	}
}

func TestIsSubSet(t *testing.T) {
	s1 := New[int]()
	s2 := New[int]()
	for i := 0; i < 10; i++ {
		s1.Insert(i)
		s2.Insert(i)
	}
	if !s1.IsSubSet(s2) {
		t.Errorf("s1.IsSubset(s2) = %t, want %t", s1.IsSubSet(s2), true)
	}
	s2.Remove(5)
	if s1.IsSubSet(s2) {
		t.Errorf("s1.IsSubset(s2) = %t, want %t", s1.IsSubSet(s2), false)
	}
}

func TestDifference(t *testing.T) {
	s1 := New[int]()
	s2 := New[int]()
	for i := 0; i < 10; i++ {
		s1.Insert(i)
		s2.Insert(i)
	}
	s2.Remove(5)
	s3 := s1.Difference(s2)
	if s3.Len() != 1 {
		t.Errorf("s3.Len() = %d, want %d", s3.Len(), 1)
	}
	if !s3.Contain(5) {
		t.Errorf("s3.Contain(5) = %t, want %t", s3.Contain(5), true)
	}
}

func TestSymmetricDifference(t *testing.T) {
	s1 := New[int]()
	s2 := New[int]()
	for i := 0; i < 10; i++ {
		s1.Insert(i)
		s2.Insert(i)
	}
	s1.Remove(0)
	s2.Remove(2)
	s3 := s1.SymmetricDifference(s2)
	if s3.Len() != 2 {
		t.Errorf("s3.Len() = %d, want %d", s3.Len(), 2)
	}
}

func TestIntersection(t *testing.T) {
	s1 := New[int]()
	s2 := New[int]()
	for i := 0; i < 10; i++ {
		s1.Insert(i)
		s2.Insert(i)
	}
	s1.Remove(0)
	s2.Remove(2)
	s3 := s1.Intersection(s2)
	if s3.Len() != 8 {
		t.Errorf("s3.Len() = %d, want %d", s3.Len(), 8)
	}
	if !s3.Contain(1) {
		t.Errorf("s3.Contain(1) = %t, want %t", s3.Contain(1), true)
	}
	if s3.Contain(0) {
		t.Errorf("s3.Contain(0) = %t, want %t", s3.Contain(0), false)
	}
}

func TestUnion(t *testing.T) {
	s1 := New[int]()
	s2 := New[int]()
	for i := 0; i < 10; i++ {
		s1.Insert(i)
		s2.Insert(i)
	}
	s1.Remove(0)
	s2.Remove(2)
	s3 := s1.Union(s2)
	if s3.Len() != 10 {
		t.Errorf("s3.Len() = %d, want %d", s3.Len(), 10)
	}
	if !s3.Contain(0) {
		t.Errorf("s3.Contain(0) = %t, want %t", s3.Contain(0), true)
	}
	if !s3.Contain(2) {
		t.Errorf("s3.Contain(2) = %t, want %t", s3.Contain(2), true)
	}
}
