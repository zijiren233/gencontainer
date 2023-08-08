// doubly linked list
package dllist

// Dllist represents a doubly linked list.
type Dllist[Item any] struct {
	root Element[Item]
	len  int

	inited bool
}

type DllistConf[Item any] func(*Dllist[Item])

func New[T any]() *Dllist[T] {
	return (&Dllist[T]{}).Clear()
}

func (l *Dllist[T]) Clear() *Dllist[T] {
	l.root.next = &l.root
	l.root.prev = &l.root
	l.len = 0
	l.inited = true
	return l
}

// only call on before use root element
func (l *Dllist[T]) lazyInit() {
	if !l.inited {
		l.Clear()
	}
}

func (l *Dllist[T]) Get(i int) *Element[T] {
	if i < 0 || i >= l.len {
		return nil
	}
	l.lazyInit()

	var e *Element[T]

	if i < l.len/2 {
		e = l.root.next
		for ; i > 0; i-- {
			e = e.next
		}
	} else {
		e = &l.root
		for ; i < l.len; i++ {
			e = e.prev
		}
	}

	return e
}

func (l *Dllist[T]) Len() int { return l.len }

func (l *Dllist[T]) Front() *Element[T] {
	if l.len == 0 {
		return nil
	}
	l.lazyInit()
	return l.root.next
}

func (l *Dllist[T]) Back() *Element[T] {
	if l.len == 0 {
		return nil
	}
	l.lazyInit()
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
	l.lazyInit()
	return l.insertValue(v, &l.root)
}

func (l *Dllist[T]) PushBack(v T) *Element[T] {
	l.lazyInit()
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
	l.lazyInit()
	l.move(e, &l.root)
}

func (l *Dllist[T]) MoveToBack(e *Element[T]) {
	if e.list != l || l.root.prev == e {
		return
	}
	l.lazyInit()
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
	l.lazyInit()
	for i, e := other.len, other.Front(); i > 0; i, e = i-1, e.Next() {
		l.insertValue(e.Value, l.root.prev)
	}
}

func (l *Dllist[T]) PushBackListTo(other *Dllist[T], mark *Element[T]) {
	if mark.list != l {
		return
	}
	for i, e := other.len, other.Front(); i > 0; i, e = i-1, e.Next() {
		l.insertValue(e.Value, mark)
	}
}

func (l *Dllist[T]) PushFrontList(other *Dllist[T]) {
	l.lazyInit()
	for i, e := other.len, other.Back(); i > 0; i, e = i-1, e.Prev() {
		l.insertValue(e.Value, &l.root)
	}
}

func (l *Dllist[T]) PushFrontListTo(other *Dllist[T], mark *Element[T]) {
	if mark.list != l {
		return
	}
	for i, e := other.len, other.Back(); i > 0; i, e = i-1, e.Prev() {
		l.insertValue(e.Value, mark)
	}
}

func (l *Dllist[T]) Swap(a, b *Element[T]) {
	if a.list != l || b.list != l || a == b {
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

// dont remove element when iterating, it will stop the iteration
func (l *Dllist[T]) Range(f func(e *Element[T]) (Continue bool)) (RangeAll bool) {
	for e := l.Front(); e != nil; e = e.Next() {
		if !f(e) {
			return
		}
	}
	return true
}

func (l *Dllist[T]) RangeFrom(start *Element[T], f func(e *Element[T]) (Continue bool)) (RangeAll bool) {
	for e := start; e != nil; e = e.Next() {
		if !f(e) {
			return
		}
	}
	return true
}

// DeepClone returns a new Dllist with a copy of l's elements.
func (l *Dllist[T]) Slice() []T {
	s := make([]T, 0, l.len)
	l.Range(func(e *Element[T]) bool {
		s = append(s, e.Value)
		return true
	})
	return s
}

// stack overflow
func (l *Dllist[T]) Sort(cmpLess func(T, T) bool) {
	if l.len <= 1 {
		return
	}

	pivot := l.Front().Value
	smaller := New[T]()
	larger := New[T]()

	l.RangeFrom(l.Front().Next(), func(e *Element[T]) bool {
		if cmpLess(e.Value, pivot) {
			smaller.PushBack(e.Value)
		} else {
			larger.PushBack(e.Value)
		}
		return true
	})

	smaller.Sort(cmpLess)
	larger.Sort(cmpLess)

	l.Clear()

	smaller.Range(func(e *Element[T]) (Continue bool) {
		l.PushBack(e.Value)
		return true
	})

	l.PushBack(pivot)

	larger.Range(func(e *Element[T]) (Continue bool) {
		l.PushBack(e.Value)
		return true
	})
}
