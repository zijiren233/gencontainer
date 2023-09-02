package heap

import "golang.org/x/exp/constraints"

var _ Interface[int] = (*MaxHeap[int])(&[]int{})

type MaxHeap[T constraints.Ordered] []T

func (h MaxHeap[T]) Len() int {
	return len(h)
}

func (h MaxHeap[T]) Less(i, j int) bool {
	return h[i] > h[j]
}

func (h MaxHeap[T]) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *MaxHeap[T]) Push(x T) {
	*h = append(*h, x)
}

func (h *MaxHeap[T]) Pop() T {
	n := len(*h)
	x := (*h)[n-1]
	*h = (*h)[:n-1]
	return x
}
