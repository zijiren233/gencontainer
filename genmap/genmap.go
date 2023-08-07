package genmap

import (
	"golang.org/x/exp/maps"
)

type GenMap[K comparable, V any] struct {
	m   map[K]V
	cap int
}

type GenMapConf[K comparable, V any] func(*GenMap[K, V])

func WithCap[K comparable, V any](cap int) GenMapConf[K, V] {
	return func(m *GenMap[K, V]) {
		m.cap = cap
	}
}

func New[K comparable, V any](conf ...GenMapConf[K, V]) *GenMap[K, V] {
	m := &GenMap[K, V]{}
	for _, c := range conf {
		c(m)
	}

	return m.Init()
}

func (m *GenMap[K, V]) Init() *GenMap[K, V] {
	if m.m == nil {
		m.m = make(map[K]V, m.cap)
	}
	return m
}

func (m *GenMap[K, V]) Len() int {
	return len(m.m)
}

func (m *GenMap[K, V]) Load(k K) (v V, ok bool) {
	v, ok = m.m[k]
	return
}

func (m *GenMap[K, V]) Store(k K, v V) {
	m.m[k] = v
}

func (m *GenMap[K, V]) Delete(k K) {
	delete(m.m, k)
}

// deep copy
func (m *GenMap[K, V]) Clone() *GenMap[K, V] {
	return New(WithCap[K, V](m.Len())).CopyFrom(m)
}

// shallow copy
func (m *GenMap[K, V]) Map() map[K]V {
	return m.m
}

func (m *GenMap[K, V]) Keys() []K {
	return maps.Keys(m.m)
}

func (m *GenMap[K, V]) Values() []V {
	return maps.Values(m.m)
}

func (m *GenMap[K, V]) Range(f func(k K, v V) (Continue bool)) (RangeAll bool) {
	for k, v := range m.m {
		if !f(k, v) {
			return
		}
	}
	return true
}

func (m *GenMap[K, V]) LoadOrStore(k K, v V) (actual V, loaded bool) {
	actual, loaded = m.m[k]
	if loaded {
		return
	}
	m.m[k] = v
	return
}

func (m *GenMap[K, V]) LoadAndDelete(k K) (v V, loaded bool) {
	v, loaded = m.m[k]
	if loaded {
		delete(m.m, k)
	}
	return
}

func (m *GenMap[K, V]) LoadAndDeleteAll() map[K]V {
	mm := m.m
	m.m = make(map[K]V)
	return mm
}

func (m *GenMap[K, V]) LoadAndStore(k K, v V) (actual V, loaded bool) {
	actual, loaded = m.m[k]
	m.m[k] = v
	return
}

func (m *GenMap[K, V]) EqualFunc(m2 *GenMap[K, V], eq func(V, V) bool) bool {
	return maps.EqualFunc(m.m, m2.m, eq)
}

func (m *GenMap[K, V]) Clear() {
	maps.Clear(m.m)
}

func (m *GenMap[K, V]) CopyTo(dst *GenMap[K, V]) (src *GenMap[K, V]) {
	maps.Copy(dst.m, m.m)
	return m
}

func (m *GenMap[K, V]) CopyToRaw(dst map[K]V) (src *GenMap[K, V]) {
	maps.Copy(dst, m.m)
	return m
}

func (m *GenMap[K, V]) CopyFrom(src *GenMap[K, V]) (dst *GenMap[K, V]) {
	maps.Copy(m.m, src.m)
	return m
}

func (m *GenMap[K, V]) CopyFromRaw(src map[K]V) (dst *GenMap[K, V]) {
	maps.Copy(m.m, src)
	return m
}

func (m *GenMap[K, V]) DeleteFunc(f func(k K, v V) bool) {
	maps.DeleteFunc(m.m, f)
}
