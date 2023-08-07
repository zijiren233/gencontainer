// linked hash map
package lhashmap

import "github.com/zijiren233/gencontainer/dllist"

type Lhashmap[K comparable, V any] struct {
	m map[K]*LhashmapNode[K, V]
	l dllist.Dllist[*LhashmapNode[K, V]]
}

type LHashMapConf[K comparable, V any] func(*Lhashmap[K, V])

func WithCap[K comparable, V any](cap int) LHashMapConf[K, V] {
	return func(m *Lhashmap[K, V]) {
		m.m = make(map[K]*LhashmapNode[K, V], cap)
	}
}

func New[K comparable, V any](conf ...LHashMapConf[K, V]) *Lhashmap[K, V] {
	m := &Lhashmap[K, V]{
		l: *dllist.New[*LhashmapNode[K, V]](nil),
	}
	for _, c := range conf {
		c(m)
	}
	if m.m == nil {
		m.m = make(map[K]*LhashmapNode[K, V])
	}
	return m
}

func (m *Lhashmap[K, V]) Len() int {
	return m.l.Len()
}

func (m *Lhashmap[K, V]) Get(k K) (v V, b bool) {
	if n, ok := m.m[k]; ok {
		return n.v, true
	}
	return
}

// If the key already exists, the value will be updated, and the key will be moved to the end of the list.
func (m *Lhashmap[K, V]) Set(k K, v V) {
	if n, ok := m.m[k]; ok {
		n.v = v
		m.l.MoveToBack(n.e)
		return
	}
	m.m[k] = NewNode(k, v).PushBack(m)
}

func (m *Lhashmap[K, V]) Delete(k K) {
	if n, ok := m.m[k]; ok {
		delete(m.m, k)
		m.l.Remove(n.e)
	}
}

// according to the order of insertion
func (m *Lhashmap[K, V]) Keys() []K {
	keys := make([]K, 0, m.Len())
	for e := m.l.Front(); e != nil; e = e.Next() {
		keys = append(keys, e.Value.k)
	}
	return keys
}

// according to the order of insertion
func (m *Lhashmap[K, V]) Values() []V {
	values := make([]V, 0, m.Len())
	for e := m.l.Front(); e != nil; e = e.Next() {
		values = append(values, e.Value.v)
	}
	return values
}

func (m *Lhashmap[K, V]) IsEmpty() bool {
	return m.Len() == 0
}

func (m *Lhashmap[K, V]) Clear() {
	m.m = make(map[K]*LhashmapNode[K, V])
	m.l.Clear()
}

func (m *Lhashmap[K, V]) Range(cbk func(e *LhashmapNode[K, V]) (Continue bool)) {
	for e := m.l.Front(); e != nil; e = e.Next() {
		if !cbk(e.Value) {
			return
		}
	}
}
