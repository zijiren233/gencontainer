package pqueue

import "github.com/zijiren233/gencontainer/heap"

var _ heap.Interface[*Item[int]] = (*maxqueue[int])(&[]*Item[int]{})

type maxqueue[T any] []*Item[T]

func (pq maxqueue[T]) Len() int {
	return len(pq)
}

func (pq maxqueue[T]) Less(i, j int) bool {
	return pq[i].priority > pq[j].priority
}

func (pq maxqueue[T]) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *maxqueue[T]) Push(i *Item[T]) {
	*pq = append(*pq, i)
}

func (pq *maxqueue[T]) Pop() *Item[T] {
	n := len(*pq)
	i := (*pq)[n-1]
	*pq = (*pq)[:n-1]
	return i
}
