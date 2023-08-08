package genmap

import "testing"

func TestStore(t *testing.T) {
	m := GenMap[string, int]{}
	m.Store("hello", 1)
	m.Store("world", 2)
	if v, ok := m.Load("hello"); !ok || v != 1 {
		t.Error("Load error")
	}
	if v, ok := m.Load("world"); !ok || v != 2 {
		t.Error("Load error")
	}
	if _, ok := m.Load("not exist"); ok {
		t.Error("Load error")
	}
}

func TestClone(t *testing.T) {
	m := GenMap[string, int]{}
	m.Store("hello", 1)
	m.Store("world", 2)
	if v, ok := m.Load("hello"); !ok || v != 1 {
		t.Error("Load error")
	}
	if v, ok := m.Load("world"); !ok || v != 2 {
		t.Error("Load error")
	}
	gm := m.Clone()
	gm.Clear()
	if v, ok := m.Load("hello"); !ok || v != 1 {
		t.Error("Load error")
	}
	if v, ok := m.Load("world"); !ok || v != 2 {
		t.Error("Load error")
	}
}

func TestMap(t *testing.T) {
	m := GenMap[string, int]{}
	m.Store("hello", 1)
	m.Store("world", 2)
	if v, ok := m.Load("hello"); !ok || v != 1 {
		t.Error("Load error")
	}
	if v, ok := m.Load("world"); !ok || v != 2 {
		t.Error("Load error")
	}
	gm := m.Map()
	for k := range gm {
		delete(gm, k)
	}
	if _, ok := m.Load("hello"); ok {
		t.Error("Load error")
	}
	if _, ok := m.Load("world"); ok {
		t.Error("Load error")
	}
}

func TestFrom(t *testing.T) {
	m := From[string, int](map[string]int{"hello": 1, "world": 2})
	if v, ok := m.Load("hello"); !ok || v != 1 {
		t.Error("Load error")
	}
	if v, ok := m.Load("world"); !ok || v != 2 {
		t.Error("Load error")
	}
}

func TestCopyFrom(t *testing.T) {
	m := GenMap[string, int]{}
	m.CopyFrom(From[string, int](map[string]int{"hello": 1, "world": 2}))
	if v, ok := m.Load("hello"); !ok || v != 1 {
		t.Error("Load error")
	}
	if v, ok := m.Load("world"); !ok || v != 2 {
		t.Error("Load error")
	}
}

func TestCopyFrom2(t *testing.T) {
	m := GenMap[string, int]{}
	gm := From[string, int](map[string]int{"hello": 1, "world": 2})
	m.CopyFrom(gm)
	gm.Clear()
	if v, ok := m.Load("hello"); !ok || v != 1 {
		t.Error("Load error")
	}
	if v, ok := m.Load("world"); !ok || v != 2 {
		t.Error("Load error")
	}
}
