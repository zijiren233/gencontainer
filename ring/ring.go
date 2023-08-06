package ring

type Ring[T any] struct {
	next, prev *Ring[T]
	Value      T
}

func (r *Ring[T]) init() *Ring[T] {
	r.next = r
	r.prev = r
	return r
}

func (r *Ring[T]) Next() *Ring[T] {
	if r.next == nil {
		return r.init()
	}
	return r.next
}

func (r *Ring[T]) Prev() *Ring[T] {
	if r.next == nil {
		return r.init()
	}
	return r.prev
}

func (r *Ring[T]) Move(n int) *Ring[T] {
	if r.next == nil {
		return r.init()
	}
	switch {
	case n < 0:
		for ; n < 0; n++ {
			r = r.prev
		}
	case n > 0:
		for ; n > 0; n-- {
			r = r.next
		}
	}
	return r
}

func New[T any](n int) *Ring[T] {
	if n <= 0 {
		return nil
	}
	r := new(Ring[T])
	p := r
	for i := 1; i < n; i++ {
		p.next = &Ring[T]{prev: p}
		p = p.next
	}
	p.next = r
	r.prev = p
	return r
}

func (r *Ring[T]) Link(s *Ring[T]) *Ring[T] {
	n := r.Next()
	if s != nil {
		p := s.Prev()

		r.next = s
		s.prev = r
		n.prev = p
		p.next = n
	}
	return n
}

func (r *Ring[T]) Unlink(n int) *Ring[T] {
	if n <= 0 {
		return nil
	}
	return r.Link(r.Move(n + 1))
}

func (r *Ring[T]) Len() int {
	n := 0
	if r != nil {
		n = 1
		for p := r.Next(); p != r; p = p.next {
			n++
		}
	}
	return n
}

func (r *Ring[T]) Do(f func(T)) {
	if r != nil {
		f(r.Value)
		for p := r.Next(); p != r; p = p.next {
			f(p.Value)
		}
	}
}
