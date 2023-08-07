package rwmap

import (
	"sync"

	"github.com/zijiren233/gencontainer/genmap"
)

type RWMap[K comparable, V any] struct {
	m *genmap.GenMap[K, V]
	l *sync.RWMutex
}

func New[K comparable, V any]() *RWMap[K, V] {
	m := &RWMap[K, V]{
		m: genmap.New[K, V](),
		l: new(sync.RWMutex),
	}
	return m
}

func (m *RWMap[K, V]) Len() int {
	m.l.RLock()
	defer m.l.RUnlock()
	return m.m.Len()
}

func (m *RWMap[K, V]) Load(k K) (v V, ok bool) {
	m.l.RLock()
	defer m.l.RUnlock()
	return m.m.Load(k)
}

func (m *RWMap[K, V]) Store(k K, v V) {
	m.l.Lock()
	defer m.l.Unlock()
	m.m.Store(k, v)
}

func (m *RWMap[K, V]) Delete(k K) {
	m.l.Lock()
	defer m.l.Unlock()
	m.m.Delete(k)
}

func (m *RWMap[K, V]) Clone() *RWMap[K, V] {
	m.l.RLock()
	defer m.l.RUnlock()
	return New[K, V]().CopyFrom(m)
}

func (m *RWMap[K, V]) Map() map[K]V {
	m.l.RLock()
	defer m.l.RUnlock()
	return m.m.Map()
}

func (m *RWMap[K, V]) Keys() []K {
	m.l.RLock()
	defer m.l.RUnlock()
	return m.m.Keys()
}

func (m *RWMap[K, V]) Values() []V {
	m.l.RLock()
	defer m.l.RUnlock()
	return m.m.Values()
}

func (m *RWMap[K, V]) Range(f func(k K, v V) (Continue bool)) (RangeAll bool) {
	m.l.RLock()
	defer m.l.RUnlock()
	return m.m.Range(f)
}

func (m *RWMap[K, V]) LoadOrStore(k K, v V) (actual V, loaded bool) {
	m.l.Lock()
	defer m.l.Unlock()
	return m.m.LoadOrStore(k, v)
}

func (m *RWMap[K, V]) LoadAndDelete(k K) (v V, loaded bool) {
	m.l.Lock()
	defer m.l.Unlock()
	return m.m.LoadAndDelete(k)
}

func (m *RWMap[K, V]) LoadAndDeleteAll() map[K]V {
	m.l.Lock()
	defer m.l.Unlock()
	return m.m.LoadAndDeleteAll()
}

func (m *RWMap[K, V]) LoadAndStore(k K, v V) (actual V, loaded bool) {
	m.l.Lock()
	defer m.l.Unlock()
	return m.m.LoadAndStore(k, v)
}

func (m *RWMap[K, V]) EqualFunc(m2 *RWMap[K, V], eq func(V, V) bool) bool {
	m.l.RLock()
	defer m.l.RUnlock()
	m2.l.RLock()
	defer m2.l.RUnlock()
	return m.m.EqualFunc(m2.m, eq)
}

func (m *RWMap[K, V]) Clear() {
	m.l.Lock()
	defer m.l.Unlock()
	m.m.Clear()
}

func (m *RWMap[K, V]) CopyFrom(src *RWMap[K, V]) (dst *RWMap[K, V]) {
	m.l.Lock()
	defer m.l.Unlock()
	src.l.RLock()
	defer src.l.RUnlock()
	m.m.CopyFrom(src.m)
	return m
}

func (m *RWMap[K, V]) CopyFromRaw(src map[K]V) (dst *RWMap[K, V]) {
	m.l.Lock()
	defer m.l.Unlock()
	m.m.CopyFromRaw(src)
	return m
}

func (m *RWMap[K, V]) CopyTo(dst *RWMap[K, V]) (src *RWMap[K, V]) {
	m.l.RLock()
	defer m.l.RUnlock()
	dst.l.Lock()
	defer dst.l.Unlock()
	m.m.CopyTo(dst.m)
	return m
}

func (m *RWMap[K, V]) CopyToRaw(dst map[K]V) (src *RWMap[K, V]) {
	m.l.RLock()
	defer m.l.RUnlock()
	m.m.CopyToRaw(dst)
	return m
}

func (m *RWMap[K, V]) DeleteFunc(f func(k K, v V) bool) {
	m.l.Lock()
	defer m.l.Unlock()
	m.m.DeleteFunc(f)
}
