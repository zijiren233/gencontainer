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

func (e *Element[T]) Remove() T {
	if e.list != nil {
		e.list.remove(e)
	}
	e.next = nil
	e.prev = nil
	return e.Value
}

func (e *Element[T]) InsertAfter(v T) *Element[T] {
	if e.list == nil {
		return nil
	}
	return e.list.insertValue(v, e)
}

func (e *Element[T]) InsertBefore(v T) *Element[T] {
	if e.list == nil {
		return nil
	}
	return e.list.insertValue(v, e.prev)
}

// dont remove element when iterating, it will stop the iteration
func (e *Element[T]) Range(f func(e *Element[T]) (Continue bool)) {
	for ; e != nil; e = e.Next() {
		if !f(e) {
			return
		}
	}
}
