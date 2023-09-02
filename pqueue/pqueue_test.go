package pqueue_test

import (
	"testing"

	"github.com/zijiren233/gencontainer/pqueue"
)

func TestMinPQueue(t *testing.T) {
	pq := pqueue.NewMinPriorityQueue[int]()
	for i := 20; i > 0; i-- {
		pq.Push(i, i)
	}
	for i := 1; pq.Len() > 0; i++ {
		priority, value := pq.Pop()
		if priority != i {
			t.Errorf("%d.th pop got %d; want %d", i, priority, i)
		}
		if value != i {
			t.Errorf("%d.th pop got %d; want %d", i, value, i)
		}
	}
}

func TestMaxPQueue(t *testing.T) {
	pq := pqueue.NewMaxPriorityQueue[int]()
	for i := 1; i <= 20; i++ {
		pq.Push(i, i)
	}
	for i := 20; pq.Len() > 0; i-- {
		priority, value := pq.Pop()
		if priority != i {
			t.Errorf("%d.th pop got %d; want %d", i, priority, i)
		}
		if value != i {
			t.Errorf("%d.th pop got %d; want %d", i, value, i)
		}
	}
}

func TestBench(t *testing.T) {
	pq := pqueue.NewMinPriorityQueue[int]()
	for i := 1; i <= 1_000_000; i++ {
		pq.Push(i, i)
	}
	for i := 1; pq.Len() > 0; i++ {
		priority, value := pq.Pop()
		if priority != i {
			t.Errorf("%d.th pop got %d; want %d", i, priority, i)
		}
		if value != i {
			t.Errorf("%d.th pop got %d; want %d", i, value, i)
		}
	}
}
