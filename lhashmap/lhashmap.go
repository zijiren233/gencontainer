// linked hash map
package lhashmap

import "github.com/zijiren233/gencontainer/dllist"

type Lhashmap[K comparable, V any] struct {
	m map[K]*dllist.Element[*entry[K, V]]
	l *dllist.Dllist[*entry[K, V]]
}

type entry[K comparable, V any] struct {
	k K
	v V
}

type LHashMapConf[K comparable, V any] func(*Lhashmap[K, V])

func WithCap[K comparable, V any](cap int) LHashMapConf[K, V] {
	return func(m *Lhashmap[K, V]) {
		if m.m == nil {
			m.m = make(map[K]*dllist.Element[*entry[K, V]], cap)
		}
	}
}

func New[K comparable, V any](conf ...LHashMapConf[K, V]) *Lhashmap[K, V] {
	m := &Lhashmap[K, V]{
		l: dllist.New[*entry[K, V]](),
	}
	for _, c := range conf {
		c(m)
	}
	if m.m == nil {
		m.m = make(map[K]*dllist.Element[*entry[K, V]])
	}
	return m
}

func (l *Lhashmap[K, V]) Load(key K) (v V, ok bool) {
	if element, ok := l.m[key]; ok {
		l.l.MoveToFront(element)
		return element.Value.v, true
	}
	return
}

func (l *Lhashmap[K, V]) Store(key K, value V) {
	if element, ok := l.m[key]; ok {
		element.Value.v = value
		l.l.MoveToFront(element)
		return
	}
	element := l.l.PushFront(&entry[K, V]{k: key, v: value})
	l.m[key] = element
}

func (l *Lhashmap[K, V]) Remove(key K) {
	if element, ok := l.m[key]; ok {
		l.l.Remove(element)
		delete(l.m, key)
	}
}

func (l *Lhashmap[K, V]) Len() int {
	return len(l.m)
}

func (l *Lhashmap[K, V]) Range(f func(k K, v V) bool) {
	for e := l.l.Front(); e != nil; e = e.Next() {
		if !f(e.Value.k, e.Value.v) {
			break
		}
	}
}

func (l *Lhashmap[K, V]) Clear() {
	l.m = make(map[K]*dllist.Element[*entry[K, V]])
	l.l.Clear()
}
