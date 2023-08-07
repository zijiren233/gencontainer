package rwmap

import "testing"

func TestLoad(t *testing.T) {
	m := New[int, int]()
	m.Store(1, 1)
	if v, ok := m.Load(1); !ok || v != 1 {
		t.Errorf("Load(1) = %v, %v, want 1, true", v, ok)
	}
}

func TestStore(t *testing.T) {
	m := New[int, int]()
	m.Store(1, 1)
	if v, ok := m.Load(1); !ok || v != 1 {
		t.Errorf("Load(1) = %v, %v, want 1, true", v, ok)
	}
	m.Store(1, 2)
	if v, ok := m.Load(1); !ok || v != 2 {
		t.Errorf("Load(1) = %v, %v, want 2, true", v, ok)
	}
}

func TestDelete(t *testing.T) {
	m := New[int, int]()
	m.Store(1, 1)
	m.Delete(1)
	if v, ok := m.Load(1); ok || v != 0 {
		t.Errorf("Load(1) = %v, %v, want 0, false", v, ok)
	}
}

func TestClear(t *testing.T) {
	m := New[int, int]()
	m.Store(1, 1)
	m.Clear()
	if v, ok := m.Load(1); ok || v != 0 {
		t.Errorf("Load(1) = %v, %v, want 0, false", v, ok)
	}
}

func TestLen(t *testing.T) {
	m := New[int, int]()
	m.Store(1, 1)
	if l := m.Len(); l != 1 {
		t.Errorf("Len() = %v, want 1", l)
	}
}

func TestWithCap(t *testing.T) {
	m := New[int, int](WithCap[int, int](10))
	if l := m.Len(); l != 0 {
		t.Errorf("Len() = %v, want 0", l)
	}
}

func TestWithValues(t *testing.T) {
	m := New[int, int](WithValues[int, int](map[int]int{1: 1}))
	if l := m.Len(); l != 1 {
		t.Errorf("Len() = %v, want 1", l)
	}
}

func TestWithValuesOverride(t *testing.T) {
	m := New[int, int](WithValues[int, int](map[int]int{1: 1}), WithValues[int, int](map[int]int{1: 2}))
	if l := m.Len(); l != 1 {
		t.Errorf("Len() = %v, want 1", l)
	}
	if v, ok := m.Load(1); !ok || v != 2 {
		t.Errorf("Load(1) = %v, %v, want 2, true", v, ok)
	}
}

func TestWithValuesOverride2(t *testing.T) {
	m := New[int, int](WithValues[int, int](map[int]int{1: 1}), WithValues[int, int](map[int]int{2: 2}))
	if l := m.Len(); l != 2 {
		t.Errorf("Len() = %v, want 2", l)
	}
	if v, ok := m.Load(1); !ok || v != 1 {
		t.Errorf("Load(1) = %v, %v, want 1, true", v, ok)
	}
	if v, ok := m.Load(2); !ok || v != 2 {
		t.Errorf("Load(2) = %v, %v, want 2, true", v, ok)
	}
}

func TestLoadOrStore(t *testing.T) {
	m := New[int, int]()
	if v, ok := m.LoadOrStore(1, 1); ok || v == 1 {
		t.Errorf("LoadOrStore(1, 1) = %v, %v, want 0, false", v, ok)
	}
	if v, ok := m.LoadOrStore(1, 2); !ok || v != 1 {
		t.Errorf("LoadOrStore(1, 2) = %v, %v, want 1, true", v, ok)
	}
}

func TestLoadAndDelete(t *testing.T) {
	m := New[int, int]()
	m.Store(1, 1)
	if v, ok := m.LoadAndDelete(1); !ok || v != 1 {
		t.Errorf("LoadAndDelete(1) = %v, %v, want 1, true", v, ok)
	}
	if v, ok := m.LoadAndDelete(1); ok || v != 0 {
		t.Errorf("LoadAndDelete(1) = %v, %v, want 0, false", v, ok)
	}
}

func TestRange(t *testing.T) {
	m := New[int, int]()
	m.Store(1, 1)
	m.Store(2, 2)
	m.Store(3, 3)
	m.Range(func(k int, v int) bool {
		if k != v {
			t.Errorf("Range() = %v, %v, want %v, %v", k, v, v, v)
		}
		return true
	})
}

func TestRangeBreak(t *testing.T) {
	m := New[int, int]()
	m.Store(1, 1)
	m.Store(2, 2)
	m.Store(3, 3)
	if m.Range(func(k int, v int) bool {
		if k != v {
			t.Errorf("Range() = %v, %v, want %v, %v", k, v, v, v)
		}
		return false
	}) {
		t.Errorf("Range() = true, want false")
	}
}

func TestLoadAndDeleteAll(t *testing.T) {
	m := New[int, int]()
	m.Store(1, 1)
	m.Store(2, 2)
	m.Store(3, 3)
	if v := m.LoadAndDeleteAll(); len(v) != 3 {
		t.Errorf("LoadAndDeleteAll() = %v, want 3", v)
	}
	if v := m.LoadAndDeleteAll(); len(v) != 0 {
		t.Errorf("LoadAndDeleteAll() = %v, want 0", v)
	}
}
