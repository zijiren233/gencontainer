package lhashmap

import "testing"

func TestAll(t *testing.T) {
	m := New[string, int]()
	m.Store("foo", 1)
	m.Store("bar", 2)
	m.Store("baz", 3)
	if m.Len() != 3 {
		t.Errorf("Len() = %v, want 3", m.Len())
	}
	if v, ok := m.Load("foo"); !ok || v != 1 {
		t.Errorf("Load(\"foo\") = %v, %v, want 1, true", v, ok)
	}
	if v, ok := m.Load("bar"); !ok || v != 2 {
		t.Errorf("Load(\"bar\") = %v, %v, want 2, true", v, ok)
	}
	if v, ok := m.Load("baz"); !ok || v != 3 {
		t.Errorf("Load(\"baz\") = %v, %v, want 3, true", v, ok)
	}
}
