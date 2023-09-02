package heap

import "golang.org/x/exp/constraints"

var _ Interface[int] = (*MinHeap[int])(&[]int{})

type MinHeap[T constraints.Ordered] []T

func (h MinHeap[T]) Len() int {
	return len(h)
}

func (h MinHeap[T]) Less(i, j int) bool {
	return h[i] < h[j]
}

func (h MinHeap[T]) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *MinHeap[T]) Push(x T) {
	*h = append(*h, x)
}

func (h *MinHeap[T]) Pop() T {
	n := len(*h)
	x := (*h)[n-1]
	*h = (*h)[:n-1]
	return x
}
