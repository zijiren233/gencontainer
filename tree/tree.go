package tree

import (
	"fmt"
	"strings"

	"github.com/zijiren233/gencontainer/utils"
)

type Node[T any] struct {
	Val                 T
	Parent, Left, Right *Node[T]
}

func NewNode[T any](val T) *Node[T] {
	return &Node[T]{Val: val}
}

func BuildTree[T any](val []T) (root *Node[T]) {
	if len(val) == 0 {
		return nil
	}
	root = NewNode(val[0])
	queue := []*Node[T]{root}
	for i := 1; i < len(val); {
		node := queue[0]
		queue = queue[1:]
		if i < len(val) {
			node.SetLeft(NewNode(val[i]))
			queue = append(queue, node.Left)
			i++
		}
		if i < len(val) {
			node.SetRight(NewNode(val[i]))
			queue = append(queue, node.Right)
			i++
		}
	}

	return
}

func (n *Node[T]) DrawTree() string {
	builder := strings.Builder{}
	layer := 0
	count := 0
	for _, node := range n.TraversalBFS() {
		if layer == 0 {
			builder.WriteString(fmt.Sprintf("%v\n", node.Val))
			layer++
			continue
		}
		builder.WriteString(fmt.Sprintf("%v\t", node.Val))
		count++
		if layer*2 == count {
			builder.WriteRune('\n')
			layer++
			count = 0
			continue
		}
	}
	return builder.String()
}

func (n *Node[T]) SetLeft(left *Node[T]) (old *Node[T]) {
	old = n.Left
	if n.Left != nil {
		n.Left.Parent = nil
	}
	n.Left = left
	if left != nil {
		left.Parent = n
	}
	return old
}

func (n *Node[T]) SetRight(right *Node[T]) (old *Node[T]) {
	old = n.Right
	if n.Right != nil {
		n.Right.Parent = nil
	}
	n.Right = right
	if right != nil {
		right.Parent = n
	}
	return old
}

func (n *Node[T]) IsLeaf() bool {
	return !n.HasChild()
}

func (n *Node[T]) IsRoot() bool {
	return !n.HasParent()
}

func (n *Node[T]) IsLeft() bool {
	return n.HasParent() && n.Parent.Left == n
}

func (n *Node[T]) IsRight() bool {
	return n.HasParent() && n.Parent.Right == n
}

func (n *Node[T]) Sibling() *Node[T] {
	if n.IsRoot() {
		return nil
	}
	if n.IsLeft() {
		return n.Parent.Right
	}
	return n.Parent.Left
}

func (n *Node[T]) Uncle() *Node[T] {
	if n.IsRoot() {
		return nil
	}
	return n.Parent.Sibling()
}

func (n *Node[T]) Grandparent() *Node[T] {
	if n.IsRoot() {
		return nil
	}
	return n.Parent.Parent
}

func (n *Node[T]) IsLeftRelation(node *Node[T]) bool {
	return node != nil && node.IsLeft() && node.Parent == n
}

func (n *Node[T]) IsRightRelation(node *Node[T]) bool {
	return node != nil && node.IsRight() && node.Parent == n
}

func (n *Node[T]) IsChildRelation(node *Node[T]) bool {
	return n.IsLeftRelation(node) || n.IsRightRelation(node)
}

func (n *Node[T]) IsUncleRelation(node *Node[T]) bool {
	return n.Uncle() == node
}

func (n *Node[T]) IsSiblingRelation(node *Node[T]) bool {
	return n.Sibling() == node
}

func (n *Node[T]) IsGrandparentRelation(node *Node[T]) bool {
	return n.Grandparent() == node
}

func (n *Node[T]) IsAncestorRelation(node *Node[T]) bool {
	if node == nil {
		return false
	}
	for node != nil {
		if node == n {
			return true
		}
		node = node.Parent
	}
	return false
}

func (n *Node[T]) IsDescendantRelation(node *Node[T]) bool {
	return !node.IsAncestorRelation(n)
}

func (n *Node[T]) HasLeft() bool {
	return n.Left != nil
}

func (n *Node[T]) HasRight() bool {
	return n.Right != nil
}

func (n *Node[T]) HasChild() bool {
	return n.HasLeft() || n.HasRight()
}

func (n *Node[T]) HasBothChild() bool {
	return n.HasLeft() && n.HasRight()
}

func (n *Node[T]) HasParent() bool {
	return n.Parent != nil
}

func (n *Node[T]) HasSibling() bool {
	return n.Sibling() != nil
}

func (n *Node[T]) HasUncle() bool {
	return n.Uncle() != nil
}

func (n *Node[T]) HasGrandparent() bool {
	return n.Grandparent() != nil
}

func (n *Node[T]) Layer() int {
	if n == nil {
		return -1
	}
	if n.IsRoot() {
		return 0
	}
	return n.Parent.Layer() + 1
}

func (n *Node[T]) MaxDepth() int {
	if n == nil || n.IsLeaf() {
		return 0
	}
	return utils.Max(n.Left.MaxDepth(), n.Right.MaxDepth()) + 1
}

func (n *Node[T]) TraversalDFS() (res []*Node[T]) {
	queue := []*Node[T]{n}
	if n.IsLeaf() {
		return queue
	}
	res = make([]*Node[T], 0, 1)
	for len(queue) > 0 {
		node := queue[len(queue)-1]
		queue = queue[:len(queue)-1]
		res = append(res, node)
		if node.HasRight() {
			queue = append(queue, node.Right)
		}
		if node.HasLeft() {
			queue = append(queue, node.Left)
		}
	}
	return
}

func (n *Node[T]) TraversalBFS() (res []*Node[T]) {
	r := n.traversalBFS(nil, 0)
	res = make([]*Node[T], 0, len(r)*2)
	for _, v := range r {
		res = append(res, v...)
	}
	return
}

func (n *Node[T]) traversalBFS(val [][]*Node[T], levle int) (res [][]*Node[T]) {
	if n == nil {
		return val
	}
	if len(val) == levle {
		val = append(val, []*Node[T]{n})
	} else {
		val[levle] = append(val[levle], n)
	}
	val = n.Left.traversalBFS(val, levle+1)
	val = n.Right.traversalBFS(val, levle+1)
	return val
}
