package heap

import (
	"fmt"
)

type Item struct {
	value    string
	priority int

	index int
}

type PriorityQueue []*Item

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].priority > pq[j].priority
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x *Item) {
	n := len(*pq)
	x.index = n
	*pq = append(*pq, x)
}

func (pq *PriorityQueue) Pop() *Item {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	item.index = -1
	*pq = old[0 : n-1]
	return item
}

func (pq *PriorityQueue) update(item *Item, value string, priority int) {
	item.value = value
	item.priority = priority
	Fix[*Item](pq, item.index)
}

func Example_priorityQueue() {
	items := map[string]int{
		"banana": 3, "apple": 2, "pear": 4,
	}

	pq := make(PriorityQueue, len(items))
	i := 0
	for value, priority := range items {
		pq[i] = &Item{
			value:    value,
			priority: priority,
			index:    i,
		}
		i++
	}
	Init[*Item](&pq)

	item := &Item{
		value:    "orange",
		priority: 1,
	}
	Push[*Item](&pq, item)
	pq.update(item, item.value, 5)

	for pq.Len() > 0 {
		item := Pop[*Item](&pq)
		fmt.Printf("%.2d:%s ", item.priority, item.value)
	}
}
