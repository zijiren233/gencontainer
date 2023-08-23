package rbtree

import "golang.org/x/exp/constraints"

type Color bool

const (
	RED   Color = false
	BLACK Color = true
)

type Node[T constraints.Ordered] struct {
	Value               T
	Color               Color
	Parent, Left, Right *Node[T]
}

func NewNode[T constraints.Ordered](value T) *Node[T] {
	return &Node[T]{Value: value}
}

type Tree[T constraints.Ordered] struct {
	Root *Node[T]
}

func NewTree[T constraints.Ordered]() *Tree[T] {
	return &Tree[T]{}
}

