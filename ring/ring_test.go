package ring

import (
	"testing"
)

func verify(t *testing.T, r *Ring[int], N int, sum int) {
	n := r.Len()
	if n != N {
		t.Errorf("r.Len() == %d; expected %d", n, N)
	}

	n = 0
	s := 0
	r.Do(func(p int) {
		n++
		s += p
	})
	if n != N {
		t.Errorf("number of forward iterations == %d; expected %d", n, N)
	}
	if sum >= 0 && s != sum {
		t.Errorf("forward ring sum = %d; expected %d", s, sum)
	}

	if r == nil {
		return
	}

	if r.next != nil {
		var p *Ring[int]
		for q := r; p == nil || q != r; q = q.next {
			if p != nil && p != q.prev {
				t.Errorf("prev = %p, expected q.prev = %p\n", p, q.prev)
			}
			p = q
		}
		if p != r.prev {
			t.Errorf("prev = %p, expected r.prev = %p\n", p, r.prev)
		}
	}

	if r.Next() != r.next {
		t.Errorf("r.Next() != r.next")
	}
	if r.Prev() != r.prev {
		t.Errorf("r.Prev() != r.prev")
	}

	if r.Move(0) != r {
		t.Errorf("r.Move(0) != r")
	}
	if r.Move(N) != r {
		t.Errorf("r.Move(%d) != r", N)
	}
	if r.Move(-N) != r {
		t.Errorf("r.Move(%d) != r", -N)
	}
	for i := 0; i < 10; i++ {
		ni := N + i
		mi := ni % N
		if r.Move(ni) != r.Move(mi) {
			t.Errorf("r.Move(%d) != r.Move(%d)", ni, mi)
		}
		if r.Move(-ni) != r.Move(-mi) {
			t.Errorf("r.Move(%d) != r.Move(%d)", -ni, -mi)
		}
	}
}

func TestCornerCases(t *testing.T) {
	var (
		r0 *Ring[int]
		r1 Ring[int]
	)

	verify(t, r0, 0, 0)
	verify(t, &r1, 1, 0)

	r1.Link(r0)
	verify(t, r0, 0, 0)
	verify(t, &r1, 1, 0)

	r1.Link(r0)
	verify(t, r0, 0, 0)
	verify(t, &r1, 1, 0)

	r1.Unlink(0)
	verify(t, &r1, 1, 0)
}

func makeN(n int) *Ring[int] {
	r := New[int](n)
	for i := 1; i <= n; i++ {
		r.Value = i
		r = r.Next()
	}
	return r
}

func sumN(n int) int { return (n*n + n) / 2 }

func TestNew(t *testing.T) {
	for i := 0; i < 10; i++ {
		r := New[int](i)
		verify(t, r, i, -1)
	}
	for i := 0; i < 10; i++ {
		r := makeN(i)
		verify(t, r, i, sumN(i))
	}
}

func TestLink1(t *testing.T) {
	r1a := makeN(1)
	var r1b Ring[int]
	r2a := r1a.Link(&r1b)
	verify(t, r2a, 2, 1)
	if r2a != r1a {
		t.Errorf("a) 2-element link failed")
	}

	r2b := r2a.Link(r2a.Next())
	verify(t, r2b, 2, 1)
	if r2b != r2a.Next() {
		t.Errorf("b) 2-element link failed")
	}

	r1c := r2b.Link(r2b)
	verify(t, r1c, 1, 1)
	verify(t, r2b, 1, 0)
}

func TestLink2(t *testing.T) {
	var r0 *Ring[int]
	r1a := &Ring[int]{Value: 42}
	r1b := &Ring[int]{Value: 77}
	r10 := makeN(10)

	r1a.Link(r0)
	verify(t, r1a, 1, 42)

	r1a.Link(r1b)
	verify(t, r1a, 2, 42+77)

	r10.Link(r0)
	verify(t, r10, 10, sumN(10))

	r10.Link(r1a)
	verify(t, r10, 12, sumN(10)+42+77)
}

func TestLink3(t *testing.T) {
	var r Ring[int]
	n := 1
	for i := 1; i < 10; i++ {
		n += i
		verify(t, r.Link(New[int](i)), n, -1)
	}
}

func TestUnlink(t *testing.T) {
	r10 := makeN(10)
	s10 := r10.Move(6)

	sum10 := sumN(10)

	verify(t, r10, 10, sum10)
	verify(t, s10, 10, sum10)

	r0 := r10.Unlink(0)
	verify(t, r0, 0, 0)

	r1 := r10.Unlink(1)
	verify(t, r1, 1, 2)
	verify(t, r10, 9, sum10-2)

	r9 := r10.Unlink(9)
	verify(t, r9, 9, sum10-2)
	verify(t, r10, 9, sum10-2)
}

func TestLinkUnlink(t *testing.T) {
	for i := 1; i < 4; i++ {
		ri := New[int](i)
		for j := 0; j < i; j++ {
			rj := ri.Unlink(j)
			verify(t, rj, j, -1)
			verify(t, ri, i-j, -1)
			ri.Link(rj)
			verify(t, ri, i, -1)
		}
	}
}

func TestMoveEmptyRing(t *testing.T) {
	var r Ring[int]

	r.Move(1)
	verify(t, &r, 1, 0)
}
