package vec

import (
	"fmt"
	"testing"

	"golang.org/x/exp/slices"
)

func TestVec(t *testing.T) {
	var v Vec[int]
	v.Push(1, 2, 3, 4, 5, 6).Pop()
	fmt.Printf("v: %v\n", v)

}

func TestPop(t *testing.T) {
	v := New[int]().Push(1, 2, 3, 4, 5, 6)
	if e, ok := v.Pop(); !ok || e != 6 {
		t.Fatal("wrong pop")
	}
	if v.Len() != 5 {
		t.Fatal("wrong length")
	}
	if !slices.Equal(v.Slice(), []int{1, 2, 3, 4, 5}) {
		t.Fatal("wrong values")
	}
}

func TestInsert(t *testing.T) {
	v := New[int]().Push(1, 2, 3, 4, 5, 6)
	v.Insert(1, 10, 11, 12)
	if v.Len() != 9 {
		t.Fatal("wrong length")
	}
	if !slices.Equal(v.Slice(), []int{1, 10, 11, 12, 2, 3, 4, 5, 6}) {
		t.Fatal("wrong values")
	}
}

func TestContain(t *testing.T) {
	v := New[int]().Push(1, 2, 3, 4, 5, 6)
	if !v.Contains(1) {
		t.Fatal("wrong contain")
	}
	if v.Contains(10) {
		t.Fatal("wrong contain")
	}
}

func TestSwap(t *testing.T) {
	v := New[int]().Push(1, 2, 3, 4, 5, 6)
	v.Swap(0, 5)
	if !slices.Equal(v.Slice(), []int{6, 2, 3, 4, 5, 1}) {
		t.Fatal("wrong values")
	}

}

func TestSplitOff(t *testing.T) {
	v := New[int]().Push(1, 2, 3, 4, 5, 6)
	v2 := v.SplitOff(2)
	if !slices.Equal(v.Slice(), []int{1, 2}) {
		t.Fatal("wrong values")
	}
	if !slices.Equal(v2.Slice(), []int{3, 4, 5, 6}) {
		t.Fatal("wrong values")
	}
}

func TestClear(t *testing.T) {
	v := New[int]().Push(1, 2, 3, 4, 5, 6)
	v.Clear()
	if v.Len() != 0 {
		t.Fatal("wrong length")
	}
	if !slices.Equal(v.Slice(), []int{}) {
		t.Fatal("wrong values")
	}
}

func TestGetSet(t *testing.T) {
	v := New[int]().Push(1, 2, 3, 4, 5, 6)
	if e, ok := v.Get(0); !ok || e != 1 {
		t.Fatal("wrong get")
	}
	v.Set(0, 10)
	if e, ok := v.Get(0); !ok || e != 10 {
		t.Fatal("wrong get")
	}
}

func TestRemove(t *testing.T) {
	v := New[int]().Push(1, 2, 3, 4, 5, 6)
	v.Insert(1, 10, 11, 12)
	if e, ok := v.Remove(1); !ok || e != 10 {
		t.Fatal("wrong remove")
	}
	if !slices.Equal(v.Slice(), []int{1, 11, 12, 2, 3, 4, 5, 6}) {
		t.Fatal("wrong values")
	}
}

func TestSort(t *testing.T) {
	v := New[int]().Push(1, 2, 3, 4, 5, 6)
	v.Insert(1, 10, 11, 12)
	if !v.Sort().EqualSlice([]int{1, 2, 3, 4, 5, 6, 10, 11, 12}) {
		t.Fatal("wrong values")
	}
}

func TestReverse(t *testing.T) {
	v := New[int]().Push(1, 2, 3, 4, 5, 6)
	if !v.Insert(1, 10, 11, 12).Sort().Reverse().EqualSlice([]int{12, 11, 10, 6, 5, 4, 3, 2, 1}) {
		t.Fatal("wrong values")
	}
}

func TestBinarySearch(t *testing.T) {
	v := New[int]().Push(1, 2, 3, 4, 5, 6)
	if i, ok := v.BinarySearch(3); !ok || i != 2 {
		t.Fatal("wrong binary search")
	}
	if i, ok := v.BinarySearch(10); ok || i != 6 {
		t.Fatal("wrong binary search")
	}
}

func TestCompact(t *testing.T) {
	v := New[int]().Push(1, 2, 3, 4, 5, 6, 6, 5, 4, 3, 2, 1)
	if !v.Insert(1, 10, 11, 12).Reverse().Sort().Compact().EqualSlice([]int{1, 2, 3, 4, 5, 6, 10, 11, 12}) {
		t.Fatalf("wrong values %v", v.Slice())
	}
}

func TestReplace(t *testing.T) {
	v := New[int]().Push(1, 2, 3, 4, 5, 6)
	v.Replace(0, 3, 10, 11, 12)
	if !slices.Equal(v.Slice(), []int{10, 11, 12, 4, 5, 6}) {
		t.Fatal("wrong values")
	}
}
