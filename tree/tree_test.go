package tree

import (
	"testing"
)

func TestLayer(t *testing.T) {
	root := NewNode(0)
	node1 := NewNode(1)
	node2 := NewNode(1)
	root.SetLeft(node1)
	root.SetRight(node2)

	node3 := NewNode(1)
	node4 := NewNode(1)
	node5 := NewNode(1)
	node6 := NewNode(1)
	node1.SetLeft(node3)
	node1.SetRight(node4)
	node2.SetLeft(node5)
	node2.SetRight(node6)

	layer := root.Layer()
	if layer != 0 {
		t.Errorf("layer length should be 0, but got %d", layer)
	}
	layer = root.Left.Layer()
	if layer != 1 {
		t.Errorf("layer length should be 1, but got %d", layer)
	}
	layer = root.Right.Layer()
	if layer != 1 {
		t.Errorf("layer length should be 1, but got %d", layer)
	}
	layer = root.Left.Left.Layer()
	if layer != 2 {
		t.Errorf("layer length should be 2, but got %d", layer)
	}
	layer = root.Left.Right.Layer()
	if layer != 2 {
		t.Errorf("layer length should be 2, but got %d", layer)
	}
	layer = root.Right.Left.Layer()
	if layer != 2 {
		t.Errorf("layer length should be 2, but got %d", layer)
	}
	node2.SetRight(nil)
	layer = root.Right.Right.Layer()
	if layer != -1 {
		t.Errorf("layer length should be -1, but got %d", layer)
	}
}

func TestRelation(t *testing.T) {
	root := NewNode(0)
	node1 := NewNode(1)
	node2 := NewNode(1)
	root.SetLeft(node1)
	root.SetRight(node2)

	node3 := NewNode(1)
	node4 := NewNode(1)
	node5 := NewNode(1)
	node6 := NewNode(1)
	node1.SetLeft(node3)
	node1.SetRight(node4)
	node2.SetLeft(node5)
	node2.SetRight(node6)

	if !root.IsRoot() {
		t.Errorf("root should be root")
	}
	if root.IsLeaf() {
		t.Errorf("root should not be leaf")
	}
	if node1.IsRoot() {
		t.Errorf("node1 should not be root")
	}
	if node1.IsRight() {
		t.Errorf("node1 should not be right")
	}
	if !node1.IsLeft() {
		t.Errorf("node1 should be left")
	}
	if !root.HasChild() {
		t.Errorf("root should have child")
	}
	if root.HasParent() {
		t.Errorf("root should not have parent")
	}
	if !node1.HasParent() {
		t.Errorf("node1 should have parent")
	}
	if !node1.HasChild() {
		t.Errorf("node1 should have child")
	}
	if node1.Sibling() != node2 {
		t.Errorf("node1's sibling should be node2")
	}
	if node2.Sibling() != node1 {
		t.Errorf("node2's sibling should be node1")
	}
	if node3.Sibling() != node4 {
		t.Errorf("node3's sibling should be node4")
	}
	if node4.Sibling() != node3 {
		t.Errorf("node4's sibling should be node3")
	}
	if node5.Sibling() != node6 {
		t.Errorf("node5's sibling should be node6")
	}
	if node6.Sibling() != node5 {
		t.Errorf("node6's sibling should be node5")
	}
	if !root.IsAncestorRelation(node1) {
		t.Errorf("root should not be ancestor of node1")
	}
	if !root.IsAncestorRelation(node2) {
		t.Errorf("root should not be ancestor of node2")
	}
	if !node1.IsAncestorRelation(node3) {
		t.Errorf("node1 should be ancestor of node3")
	}
	if !node1.IsAncestorRelation(node4) {
		t.Errorf("node1 should be ancestor of node4")
	}
	if !node2.IsAncestorRelation(node5) {
		t.Errorf("node2 should be ancestor of node5")
	}
	if !node2.IsAncestorRelation(node6) {
		t.Errorf("node2 should be ancestor of node6")
	}
	if !node1.IsDescendantRelation(node3) {
		t.Errorf("node1 should be descendant of node3")
	}
	// IsGrandparentRelation
	if !node3.IsGrandparentRelation(root) {
		t.Errorf("root should be grandparent of node3")
	}
	if !node4.IsGrandparentRelation(root) {
		t.Errorf("root should be grandparent of node4")
	}
	if !node5.IsGrandparentRelation(root) {
		t.Errorf("root should be grandparent of node5")
	}
	if !node6.IsGrandparentRelation(root) {
		t.Errorf("root should be grandparent of node6")
	}
	// IsUncleRelation
	if !node3.IsUncleRelation(node2) {
		t.Errorf("node6 should be uncle of node2")
	}
	if !node4.IsUncleRelation(node2) {
		t.Errorf("node6 should be uncle of node2")
	}
	if !node5.IsUncleRelation(node1) {
		t.Errorf("node5 should be uncle of node1")
	}
	if !node6.IsUncleRelation(node1) {
		t.Errorf("node5 should be uncle of node1")
	}
}

func TestTraversalDFS(t *testing.T) {
	root := NewNode(0)
	node1 := NewNode(1)
	node2 := NewNode(2)
	root.SetLeft(node1)
	root.SetRight(node2)

	node3 := NewNode(3)
	node4 := NewNode(4)
	node5 := NewNode(5)
	node6 := NewNode(6)
	node1.SetLeft(node3)
	node1.SetRight(node4)
	node2.SetLeft(node5)
	node2.SetRight(node6)

	res := root.TraversalDFS()
	if len(res) != 7 {
		t.Errorf("length of result should be 7, but got %d", len(res))
	}
	if res[0] != root {
		t.Errorf("root should be first")
	}
	if res[1] != node1 {
		t.Errorf("node1 should be second")
	}
	if res[2] != node3 {
		t.Errorf("node3 should be third")
	}
	if res[3] != node4 {
		t.Errorf("node4 should be fourth")
	}
	if res[4] != node2 {
		t.Errorf("node2 should be fifth")
	}
	if res[5] != node5 {
		t.Errorf("node5 should be sixth")
	}
	if res[6] != node6 {
		t.Errorf("node6 should be seventh")
	}
}

func TestTraversalBFS(t *testing.T) {
	root := NewNode(0)
	node1 := NewNode(1)
	node2 := NewNode(2)
	root.SetLeft(node1)
	root.SetRight(node2)

	node3 := NewNode(3)
	node4 := NewNode(4)
	node5 := NewNode(5)
	node6 := NewNode(6)
	node1.SetLeft(node3)
	node1.SetRight(node4)
	node2.SetLeft(node5)
	node2.SetRight(node6)

	res := root.TraversalBFS()
	if len(res) != 7 {
		t.Errorf("length of result should be 7, but got %d", len(res))
	}
	if res[0] != root {
		t.Errorf("root should be first")
	}
	if res[1] != node1 {
		t.Errorf("node1 should be second")
	}
	if res[2] != node2 {
		t.Errorf("node2 should be third")
	}
	if res[3] != node3 {
		t.Errorf("node3 should be fourth")
	}
	if res[4] != node4 {
		t.Errorf("node4 should be fifth")
	}
	if res[5] != node5 {
		t.Errorf("node5 should be sixth")
	}
	if res[6] != node6 {
		t.Errorf("node6 should be seventh")
	}
}

func TestBuildTree(t *testing.T) {
	preorder := []int{0, 1, 3, 4, 2, 5, 6}
	root := BuildTree(preorder)
	if root.Val != 0 {
		t.Errorf("root value should be 0")
	}
	if root.Left.Val != 1 {
		t.Errorf("root's left value should be 1")
	}
	if root.Right.Val != 3 {
		t.Errorf("root's right value should be 2")
	}
	if root.Left.Left.Val != 4 {
		t.Errorf("root's left's left value should be 4")
	}
	if root.Left.Right.Val != 2 {
		t.Errorf("root's left's right value should be 2")
	}
	if root.Right.Left.Val != 5 {
		t.Errorf("root's right's left value should be 5")
	}
	if root.Right.Right.Val != 6 {
		t.Errorf("root's right's right value should be 6")
	}
}
