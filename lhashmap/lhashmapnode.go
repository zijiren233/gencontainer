package lhashmap

import "github.com/zijiren233/gencontainer/dllist"

type LhashmapNode[K comparable, V any] struct {
	k K
	v V

	e *dllist.Element[*LhashmapNode[K, V]]
}

func NewNode[K comparable, V any](k K, v V) *LhashmapNode[K, V] {
	return &LhashmapNode[K, V]{k: k, v: v}
}

func (n *LhashmapNode[K, V]) PushBack(m *Lhashmap[K, V]) *LhashmapNode[K, V] {
	n.e = m.l.PushBack(n)
	m.m[n.k] = n
	return n
}

func (n *LhashmapNode[K, V]) Next() *LhashmapNode[K, V] {
	if e := n.e.Next(); e != nil {
		return e.Value
	}
	return nil
}

func (n *LhashmapNode[K, V]) Prev() *LhashmapNode[K, V] {
	if e := n.e.Prev(); e != nil {
		return e.Value
	}
	return nil
}

func (n *LhashmapNode[K, V]) Key() K {
	return n.k
}

func (n *LhashmapNode[K, V]) Value() V {
	return n.v
}
