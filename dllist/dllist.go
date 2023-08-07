// doubly linked list
package dllist

type Element[Item any] struct {
	next, prev *Element[Item]

	list *Dllist[Item]

	Value Item
}

func (e *Element[T]) Next() *Element[T] {
	if p := e.next; e.list != nil && p != &e.list.root {
		return p
	}
	return nil
}

func (e *Element[T]) Prev() *Element[T] {
	if p := e.prev; e.list != nil && p != &e.list.root {
		return p
	}
	return nil
}

// Dllist represents a doubly linked list.
type Dllist[Item any] struct {
	root Element[Item]
	len  int
}

func (l *Dllist[T]) Clear() *Dllist[T] {
	l.root.next = &l.root
	l.root.prev = &l.root
	l.len = 0
	return l
}

func New[T any]() *Dllist[T] {
	return new(Dllist[T]).Clear()
}

func (l *Dllist[T]) Len() int { return l.len }

func (l *Dllist[T]) Front() *Element[T] {
	if l.len == 0 {
		return nil
	}
	return l.root.next
}

func (l *Dllist[T]) Back() *Element[T] {
	if l.len == 0 {
		return nil
	}
	return l.root.prev
}

func (l *Dllist[T]) insert(e, at *Element[T]) *Element[T] {
	e.prev = at
	e.next = at.next
	e.prev.next = e
	e.next.prev = e
	e.list = l
	l.len++
	return e
}

func (l *Dllist[T]) insertValue(v T, at *Element[T]) *Element[T] {
	return l.insert(&Element[T]{Value: v}, at)
}

func (l *Dllist[T]) remove(e *Element[T]) {
	e.prev.next = e.next
	e.next.prev = e.prev
	e.next = nil
	e.prev = nil
	e.list = nil
	l.len--
}

func (l *Dllist[T]) move(e, at *Element[T]) {
	if e == at {
		return
	}
	e.prev.next = e.next
	e.next.prev = e.prev

	e.prev = at
	e.next = at.next
	e.prev.next = e
	e.next.prev = e
}

func (l *Dllist[T]) Remove(e *Element[T]) T {
	if e.list == l {

		l.remove(e)
	}
	return e.Value
}

func (l *Dllist[T]) PushFront(v T) *Element[T] {
	return l.insertValue(v, &l.root)
}

func (l *Dllist[T]) PushBack(v T) *Element[T] {
	return l.insertValue(v, l.root.prev)
}

func (l *Dllist[T]) InsertBefore(v T, mark *Element[T]) *Element[T] {
	if mark.list != l {
		return nil
	}

	return l.insertValue(v, mark.prev)
}

func (l *Dllist[T]) InsertAfter(v T, mark *Element[T]) *Element[T] {
	if mark.list != l {
		return nil
	}

	return l.insertValue(v, mark)
}

func (l *Dllist[T]) MoveToFront(e *Element[T]) {
	if e.list != l || l.root.next == e {
		return
	}

	l.move(e, &l.root)
}

func (l *Dllist[T]) MoveToBack(e *Element[T]) {
	if e.list != l || l.root.prev == e {
		return
	}

	l.move(e, l.root.prev)
}

func (l *Dllist[T]) MoveBefore(e, mark *Element[T]) {
	if e.list != l || e == mark || mark.list != l {
		return
	}
	l.move(e, mark.prev)
}

func (l *Dllist[T]) MoveAfter(e, mark *Element[T]) {
	if e.list != l || e == mark || mark.list != l {
		return
	}
	l.move(e, mark)
}

func (l *Dllist[T]) PushBackList(other *Dllist[T]) {
	for i, e := other.Len(), other.Front(); i > 0; i, e = i-1, e.Next() {
		l.insertValue(e.Value, l.root.prev)
	}
}

func (l *Dllist[T]) PushFrontList(other *Dllist[T]) {
	for i, e := other.Len(), other.Back(); i > 0; i, e = i-1, e.Prev() {
		l.insertValue(e.Value, &l.root)
	}
}

func (l *Dllist[T]) Swap(a, b *Element[T]) {
	if a.list != l || b.list != l {
		return
	}
	if a == b {
		return
	}
	if a.next == b {
		a.prev.next = b
		b.next.prev = a
		a.next = b.next
		b.prev = a.prev
		a.prev = b
		b.next = a
		return
	}
	if b.next == a {
		b.prev.next = a
		a.next.prev = b
		b.next = a.next
		a.prev = b.prev
		b.prev = a
		a.next = b
		return
	}
	a.prev.next = b
	a.next.prev = b
	b.prev.next = a
	b.next.prev = a
	a.prev, b.prev = b.prev, a.prev
	a.next, b.next = b.next, a.next
}
