package vec

import (
	"testing"

	"github.com/zijiren233/gencontainer/utils"
)

func TestVec(t *testing.T) {
	v := New[int](WithValues(1, 2, 3), WithValues(4, 5, 6), WithCap[int](10))
	if v.Len() != 6 {
		t.Fatal("wrong length")
	}
	if v.Cap() != 10 {
		t.Fatal("wrong capacity")
	}
	if !utils.EqualSlice(v.Slice(), []int{1, 2, 3, 4, 5, 6}) {
		t.Fatal("wrong values")
	}
	v.Push(7, 8)
	if v.Len() != 8 {
		t.Fatal("wrong length")
	}
	if !utils.EqualSlice(v.Slice(), []int{1, 2, 3, 4, 5, 6, 7, 8}) {
		t.Fatal("wrong values")
	}
}

func TestPop(t *testing.T) {
	v := New[int](WithValues(1, 2, 3), WithValues(4, 5, 6), WithCap[int](10))
	if e, ok := v.Pop(); !ok || e != 6 {
		t.Fatal("wrong pop")
	}
	if v.Len() != 5 {
		t.Fatal("wrong length")
	}
	if !utils.EqualSlice(v.Slice(), []int{1, 2, 3, 4, 5}) {
		t.Fatal("wrong values")
	}
}

func TestResize(t *testing.T) {
	v := New[int](WithValues(1, 2, 3), WithValues(4, 5, 6), WithCap[int](10))
	v.Resize(2)
	if v.Len() != 2 {
		t.Fatal("wrong length")
	}
	if v.Cap() != 10 {
		t.Fatal("wrong capacity")
	}
	if !utils.EqualSlice(v.Slice(), []int{1, 2}) {
		t.Fatal("wrong values")
	}
	v.Resize(20)
	if v.Len() != 2 {
		t.Fatal("wrong length")
	}
	if v.Cap() != 20 {
		t.Fatal("wrong capacity")
	}
	if !utils.EqualSlice(v.Slice(), []int{1, 2}) {
		t.Fatal("wrong values")
	}
}

func TestInsert(t *testing.T) {
	v := New[int](WithValues(1, 2, 3), WithValues(4, 5, 6), WithCap[int](10))
	v.Insert(1, 10, 11, 12)
	if v.Len() != 9 {
		t.Fatal("wrong length")
	}
	if v.Cap() != 10 {
		t.Fatal("wrong capacity")
	}
	if !utils.EqualSlice(v.Slice(), []int{1, 10, 11, 12, 2, 3, 4, 5, 6}) {
		t.Fatal("wrong values")
	}
}

func TestContain(t *testing.T) {
	v := New[int](WithValues(1, 2, 3), WithValues(4, 5, 6), WithCap[int](10))
	if !v.Contain(func(v int) (matched bool) {
		return v == 1
	}) {
		t.Fatal("wrong contain")
	}
	if v.Contain(func(v int) (matched bool) {
		return v == 10
	}) {
		t.Fatal("wrong contain")
	}
}

func TestSwap(t *testing.T) {
	v := New[int](WithValues(1, 2, 3), WithValues(4, 5, 6), WithCap[int](10))
	v.Swap(0, 5)
	if !utils.EqualSlice(v.Slice(), []int{6, 2, 3, 4, 5, 1}) {
		t.Fatal("wrong values")
	}
}

func TestSplitOff(t *testing.T) {
	v := New[int](WithValues(1, 2, 3), WithValues(4, 5, 6), WithCap[int](10))
	v2, ok := v.SplitOff(2)
	if !ok {
		t.Fatal("wrong split off")
	}
	if !utils.EqualSlice(v.Slice(), []int{1, 2}) {
		t.Fatal("wrong values")
	}
	if !utils.EqualSlice(v2.Slice(), []int{3, 4, 5, 6}) {
		t.Fatal("wrong values")
	}
}

func TestClear(t *testing.T) {
	v := New[int](WithValues(1, 2, 3), WithValues(4, 5, 6), WithCap[int](10))
	v.Clear()
	if v.Len() != 0 {
		t.Fatal("wrong length")
	}
	if !utils.EqualSlice(v.Slice(), []int{}) {
		t.Fatal("wrong values")
	}
}
