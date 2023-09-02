package pqueue

import "github.com/zijiren233/gencontainer/heap"

var _ heap.Interface[*Item[int]] = (*minqueue[int])(&[]*Item[int]{})

type minqueue[T any] []*Item[T]

func (pq minqueue[T]) Len() int {
	return len(pq)
}

func (pq minqueue[T]) Less(i, j int) bool {
	return pq[i].priority < pq[j].priority
}

func (pq minqueue[T]) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *minqueue[T]) Push(i *Item[T]) {
	*pq = append(*pq, i)
}

func (pq *minqueue[T]) Pop() *Item[T] {
	n := len(*pq)
	i := (*pq)[n-1]
	*pq = (*pq)[:n-1]
	return i
}
