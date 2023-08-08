// doubly linked list
package dllist

import (
	"testing"
)

func checkListLen[T any](t *testing.T, l *Dllist[T], len int) bool {
	if n := l.Len(); n != len {
		t.Errorf("l.Len() = %d, want %d", n, len)
		return false
	}
	return true
}

func checkListPointers[T any](t *testing.T, l *Dllist[T], es []*Element[T]) {
	root := &l.root

	if !checkListLen(t, l, len(es)) {
		return
	}

	if len(es) == 0 {
		if l.root.next != nil && l.root.next != root || l.root.prev != nil && l.root.prev != root {
			t.Errorf("l.root.next = %p, l.root.prev = %p; both should both be nil or %p", l.root.next, l.root.prev, root)
		}
		return
	}

	for i, e := range es {
		prev := root
		Prev := (*Element[T])(nil)
		if i > 0 {
			prev = es[i-1]
			Prev = prev
		}
		if p := e.prev; p != prev {
			t.Errorf("elt[%d](%p).prev = %p, want %p", i, e, p, prev)
		}
		if p := e.Prev(); p != Prev {
			t.Errorf("elt[%d](%p).Prev() = %p, want %p", i, e, p, Prev)
		}

		next := root
		Next := (*Element[T])(nil)
		if i < len(es)-1 {
			next = es[i+1]
			Next = next
		}
		if n := e.next; n != next {
			t.Errorf("elt[%d](%p).next = %p, want %p", i, e, n, next)
		}
		if n := e.Next(); n != Next {
			t.Errorf("elt[%d](%p).Next() = %p, want %p", i, e, n, Next)
		}
	}
}

func checkList(t *testing.T, l *Dllist[int], es []any) {
	if !checkListLen(t, l, len(es)) {
		return
	}

	i := 0
	l.Range(func(e *Element[int]) (Continue bool) {
		le := e.Value
		if le != es[i] {
			t.Errorf("elt[%d].Value = %v, want %v", i, le, es[i])
		}
		i++
		return true
	})
}

func TestExtending(t *testing.T) {
	l1 := &Dllist[int]{}
	l2 := &Dllist[int]{}

	l1.PushBack(1)
	l1.PushBack(2)
	l1.PushBack(3)

	l2.PushBack(4)
	l2.PushBack(5)

	l3 := &Dllist[int]{}
	l3.PushBackList(l1)
	checkList(t, l3, []any{1, 2, 3})
	l3.PushBackList(l2)
	checkList(t, l3, []any{1, 2, 3, 4, 5})

	l3 = &Dllist[int]{}
	l3.PushFrontList(l2)
	checkList(t, l3, []any{4, 5})
	l3.PushFrontList(l1)
	checkList(t, l3, []any{1, 2, 3, 4, 5})

	checkList(t, l1, []any{1, 2, 3})
	checkList(t, l2, []any{4, 5})

	l3 = &Dllist[int]{}
	l3.PushBackList(l1)
	checkList(t, l3, []any{1, 2, 3})
	l3.PushBackList(l3)
	checkList(t, l3, []any{1, 2, 3, 1, 2, 3})

	l3 = &Dllist[int]{}
	l3.PushFrontList(l1)
	checkList(t, l3, []any{1, 2, 3})
	l3.PushFrontList(l3)
	checkList(t, l3, []any{1, 2, 3, 1, 2, 3})

	l3 = &Dllist[int]{}
	l1.PushBackList(l3)
	checkList(t, l1, []any{1, 2, 3})
	l1.PushFrontList(l3)
	checkList(t, l1, []any{1, 2, 3})
}

func TestRemove(t *testing.T) {
	l := &Dllist[int]{}
	e1 := l.PushBack(1)
	e2 := l.PushBack(2)
	checkListPointers(t, l, []*Element[int]{e1, e2})
	e := l.Front()
	l.Remove(e)
	checkListPointers(t, l, []*Element[int]{e2})
	l.Remove(e)
	checkListPointers(t, l, []*Element[int]{e2})
}

func TestIssue4103(t *testing.T) {
	l1 := &Dllist[int]{}
	l1.PushBack(1)
	l1.PushBack(2)

	l2 := New[int]()
	l2.PushBack(3)
	l2.PushBack(4)

	e := l1.Front()
	l2.Remove(e)
	if n := l2.Len(); n != 2 {
		t.Errorf("l2.Len() = %d, want 2", n)
	}

	l1.InsertBefore(8, e)
	if n := l1.Len(); n != 3 {
		t.Errorf("l1.Len() = %d, want 3", n)
	}
}

func TestIssue6349(t *testing.T) {
	l := &Dllist[int]{}
	l.PushBack(1)
	l.PushBack(2)

	e := l.Front()
	l.Remove(e)
	if e.Value != 1 {
		t.Errorf("e.value = %d, want 1", e.Value)
	}
	if e.Next() != nil {
		t.Errorf("e.Next() != nil")
	}
	if e.Prev() != nil {
		t.Errorf("e.Prev() != nil")
	}
}

func TestMove(t *testing.T) {
	l := &Dllist[int]{}
	e1 := l.PushBack(1)
	e2 := l.PushBack(2)
	e3 := l.PushBack(3)
	e4 := l.PushBack(4)

	l.MoveAfter(e3, e3)
	checkListPointers(t, l, []*Element[int]{e1, e2, e3, e4})
	l.MoveBefore(e2, e2)
	checkListPointers(t, l, []*Element[int]{e1, e2, e3, e4})

	l.MoveAfter(e3, e2)
	checkListPointers(t, l, []*Element[int]{e1, e2, e3, e4})
	l.MoveBefore(e2, e3)
	checkListPointers(t, l, []*Element[int]{e1, e2, e3, e4})

	l.MoveBefore(e2, e4)
	checkListPointers(t, l, []*Element[int]{e1, e3, e2, e4})
	e2, e3 = e3, e2

	l.MoveBefore(e4, e1)
	checkListPointers(t, l, []*Element[int]{e4, e1, e2, e3})
	e1, e2, e3, e4 = e4, e1, e2, e3

	l.MoveAfter(e4, e1)
	checkListPointers(t, l, []*Element[int]{e1, e4, e2, e3})
	e2, e3, e4 = e4, e2, e3

	l.MoveAfter(e2, e3)
	checkListPointers(t, l, []*Element[int]{e1, e3, e2, e4})
}

func TestSwapElement(t *testing.T) {
	l := &Dllist[int]{}
	e1 := l.PushBack(1)
	e2 := l.PushBack(2)
	e3 := l.PushBack(3)
	l.Swap(e2, e2)
	checkListPointers(t, l, []*Element[int]{e1, e2, e3})
	l.Swap(e2, e3)
	checkListPointers(t, l, []*Element[int]{e1, e3, e2})
	l.Swap(e1, e3)
	checkListPointers(t, l, []*Element[int]{e3, e1, e2})
	l.Swap(e1, e1)
	checkListPointers(t, l, []*Element[int]{e3, e1, e2})
}

func TestPushAfter(t *testing.T) {
	l := &Dllist[int]{}
	l.PushBack(1)
	e := l.PushBack(2)
	l.PushBack(3)
	e.InsertAfter(4)
	checkList(t, l, []any{1, 2, 4, 3})
}

func TestSort(t *testing.T) {
	l := &Dllist[int]{}
	l.PushBack(3)
	l.PushBack(1)
	l.PushBack(2)
	l.Sort(func(a, b int) bool { return a < b })
	checkList(t, l, []any{1, 2, 3})
}
