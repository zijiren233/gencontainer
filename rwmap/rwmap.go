package rwmap

import (
	"sync"

	"golang.org/x/exp/maps"
)

type RWMap[K comparable, V any] struct {
	m map[K]V
	l sync.RWMutex
}

type RWMapConf[K comparable, V any] func(*RWMap[K, V])

func WithCap[K comparable, V any](cap int) RWMapConf[K, V] {
	return func(m *RWMap[K, V]) {
		if m.m == nil {
			m.m = make(map[K]V, cap)
		}
	}
}

func WithValues[K comparable, V any](val map[K]V) RWMapConf[K, V] {
	return func(m *RWMap[K, V]) {
		if m.m == nil {
			m.m = val
		} else {
			for k, v := range val {
				m.m[k] = v
			}
		}
	}
}

func New[K comparable, V any](conf ...RWMapConf[K, V]) *RWMap[K, V] {
	m := &RWMap[K, V]{}
	for _, c := range conf {
		c(m)
	}
	if m.m == nil {
		m.m = make(map[K]V)
	}
	return m
}

func (m *RWMap[K, V]) Len() int {
	m.l.RLock()
	defer m.l.RUnlock()
	return len(m.m)
}

func (m *RWMap[K, V]) Load(k K) (v V, ok bool) {
	m.l.RLock()
	defer m.l.RUnlock()
	v, ok = m.m[k]
	return
}

func (m *RWMap[K, V]) Store(k K, v V) {
	m.l.Lock()
	defer m.l.Unlock()
	m.m[k] = v
}

func (m *RWMap[K, V]) Delete(k K) {
	m.l.Lock()
	defer m.l.Unlock()
	delete(m.m, k)
}

func (m *RWMap[K, V]) Clone() *RWMap[K, V] {
	m.l.RLock()
	defer m.l.RUnlock()

	return New(WithValues(maps.Clone(m.m)))
}

func (m *RWMap[K, V]) Keys() []K {
	m.l.RLock()
	defer m.l.RUnlock()
	return maps.Keys(m.m)
}

func (m *RWMap[K, V]) Values() []V {
	m.l.RLock()
	defer m.l.RUnlock()
	return maps.Values(m.m)
}

func (m *RWMap[K, V]) Range(f func(k K, v V) (Continue bool)) (RangeAll bool) {
	m.l.RLock()
	defer m.l.RUnlock()
	for k, v := range m.m {
		if !f(k, v) {
			return
		}
	}
	return true
}

func (m *RWMap[K, V]) LoadOrStore(k K, v V) (actual V, loaded bool) {
	m.l.Lock()
	defer m.l.Unlock()
	actual, loaded = m.m[k]
	if loaded {
		return
	}
	m.m[k] = v
	return
}

func (m *RWMap[K, V]) LoadAndDelete(k K) (v V, loaded bool) {
	m.l.Lock()
	defer m.l.Unlock()
	v, loaded = m.m[k]
	if loaded {
		delete(m.m, k)
	}
	return
}

func (m *RWMap[K, V]) LoadAndDeleteAll() map[K]V {
	m.l.Lock()
	defer m.l.Unlock()
	mm := m.m
	m.m = make(map[K]V)
	return mm
}

func (m *RWMap[K, V]) LoadAndStore(k K, v V) (actual V, loaded bool) {
	m.l.Lock()
	defer m.l.Unlock()
	actual, loaded = m.m[k]
	m.m[k] = v
	return
}

func (m *RWMap[K, V]) EqualFunc(m2 *RWMap[K, V], eq func(V, V) bool) bool {
	m.l.RLock()
	defer m.l.RUnlock()
	m2.l.RLock()
	defer m2.l.RUnlock()
	return maps.EqualFunc(m.m, m2.m, eq)
}

func (m *RWMap[K, V]) Clear() {
	m.l.Lock()
	defer m.l.Unlock()
	maps.Clear(m.m)
}

func (m *RWMap[K, V]) Copy(src *RWMap[K, V]) {
	m.l.Lock()
	defer m.l.Unlock()
	src.l.RLock()
	defer src.l.RUnlock()
	maps.Copy(m.m, src.m)
}

func (m *RWMap[K, V]) CopyTo(dst *RWMap[K, V]) {
	m.l.RLock()
	defer m.l.RUnlock()
	dst.l.Lock()
	defer dst.l.Unlock()
	maps.Copy(dst.m, m.m)
}

func (m *RWMap[K, V]) DeleteFunc(f func(k K, v V) bool) {
	m.l.Lock()
	defer m.l.Unlock()
	maps.DeleteFunc(m.m, f)
}
