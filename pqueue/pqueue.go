// PriorityQueue
package pqueue

import "github.com/zijiren233/gencontainer/heap"

type Item[T any] struct {
	priority int
	value    T
}

func newItem[T any](priority int, value T) *Item[T] {
	return &Item[T]{priority, value}
}

type PQueue[T any] struct {
	heap heap.Interface[*Item[T]]
}

func NewMinPriorityQueue[T any]() *PQueue[T] {
	return &PQueue[T]{&minqueue[T]{}}
}

func NewMaxPriorityQueue[T any]() *PQueue[T] {
	return &PQueue[T]{&maxqueue[T]{}}
}

func (pq *PQueue[T]) Len() int {
	return pq.heap.Len()
}

func (pq *PQueue[T]) Push(priority int, value T) {
	heap.Push(pq.heap, newItem(priority, value))
}

func (pq *PQueue[T]) Pop() (int, T) {
	item := heap.Pop(pq.heap)
	return item.priority, item.value
}
