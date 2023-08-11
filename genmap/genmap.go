package genmap

import (
	"golang.org/x/exp/maps"
)

type GenMap[K comparable, V any] map[K]V

func New[K comparable, V any]() GenMap[K, V] {
	return make(GenMap[K, V])
}

func From[K comparable, V any](m map[K]V) GenMap[K, V] {
	return m
}

func (m GenMap[K, V]) Len() int {
	return len(m)
}

func (m GenMap[K, V]) Load(k K) (v V, ok bool) {
	v, ok = m[k]
	return
}

func (m GenMap[K, V]) Store(k K, v V) GenMap[K, V] {
	m[k] = v
	return m
}

func (m GenMap[K, V]) Delete(k K) GenMap[K, V] {
	delete(m, k)
	return m
}

// deep copy
func (m GenMap[K, V]) Clone() GenMap[K, V] {
	return New[K, V]().CopyFrom(m)
}

// shallow copy
func (m GenMap[K, V]) Map() map[K]V {
	return m
}

func (m GenMap[K, V]) Keys() []K {
	return maps.Keys(m)
}

func (m GenMap[K, V]) Values() []V {
	return maps.Values(m)
}

func (m GenMap[K, V]) Range(f func(k K, v V) (Continue bool)) (RangeAll bool) {
	for k, v := range m {
		if !f(k, v) {
			return
		}
	}
	return true
}

func (m GenMap[K, V]) LoadOrStore(k K, v V) (actual V, loaded bool) {
	actual, loaded = m[k]
	if loaded {
		return
	}
	m[k] = v
	return
}

func (m GenMap[K, V]) LoadOrStoreFunc(k K, newFunc func() V) (actual V, loaded bool) {
	actual, loaded = m[k]
	if loaded {
		return
	}
	actual = newFunc()
	m[k] = actual
	return
}

func (m GenMap[K, V]) LoadAndDelete(k K) (v V, loaded bool) {
	v, loaded = m[k]
	if loaded {
		delete(m, k)
	}
	return
}

func (m GenMap[K, V]) LoadAndStore(k K, v V) (actual V, loaded bool) {
	actual, loaded = m[k]
	m[k] = v
	return
}

func (m GenMap[K, V]) EqualFunc(m2 GenMap[K, V], eq func(V, V) bool) bool {
	return maps.EqualFunc(m, m2, eq)
}

func (m GenMap[K, V]) Clear() GenMap[K, V] {
	maps.Clear(m)
	return m
}

func (m GenMap[K, V]) CopyTo(dst GenMap[K, V]) GenMap[K, V] {
	maps.Copy(dst, m)
	return m
}

func (m GenMap[K, V]) CopyToRaw(dst map[K]V) GenMap[K, V] {
	maps.Copy(dst, m)
	return m
}

func (m GenMap[K, V]) CopyFrom(src GenMap[K, V]) GenMap[K, V] {
	maps.Copy(m, src)
	return m
}

func (m GenMap[K, V]) CopyFromRaw(src map[K]V) GenMap[K, V] {
	maps.Copy(m, src)
	return m
}

func (m GenMap[K, V]) DeleteFunc(f func(k K, v V) bool) {
	maps.DeleteFunc(m, f)
}
