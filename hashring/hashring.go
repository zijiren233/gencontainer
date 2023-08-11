package hashring

import (
	"fmt"
	"hash/crc32"

	"github.com/zijiren233/gencontainer/set"
	"github.com/zijiren233/gencontainer/vec"
	"golang.org/x/exp/constraints"
)

type HashRing[Node constraints.Ordered] struct {
	replicas int
	rawNodes set.Set[Node]
	nodes    *vec.Vec[node[Node]]
}

type node[Node constraints.Ordered] struct {
	Node Node
	hash uint32
}

func New[Node constraints.Ordered](replicas int) *HashRing[Node] {
	hr := &HashRing[Node]{
		replicas: replicas,
		rawNodes: set.New[Node](),
		nodes: vec.New[node[Node]](
			vec.WithCmpLess(func(t1, t2 node[Node]) bool {
				return t1.hash < t2.hash
			}),
			vec.WithCmpEqual(func(t1, t2 node[Node]) bool {
				return t1.hash == t2.hash
			})),
	}
	return hr
}

func (hr *HashRing[Node]) AddNodes(nodes ...Node) *HashRing[Node] {
	if len(nodes) == 0 {
		return hr
	}
	s := set.New[Node]().Push(nodes...).Difference(hr.rawNodes)
	s.Range(func(val Node) (Continue bool) {
		for v := 0; v < hr.replicas; v++ {
			hr.nodes.Push(node[Node]{
				Node: val,
				hash: hr.hashKey(val, v),
			})
		}
		return true
	})
	hr.rawNodes.Push(s.Slice()...)
	hr.nodes.Sort()
	return hr
}

func (hr *HashRing[Node]) ResetNodes(nodes ...Node) {
	hr.nodes.Clear()
	hr.rawNodes.Clear()
	hr.AddNodes(nodes...)
}

// GetNode returns the node that a given key maps to
func (hr *HashRing[Node]) GetNode(key Node) (n Node) {
	if hr.nodes.Len() == 0 {
		return
	}
	hash := hr.hashKey(key, 0)
	if i, ok := hr.nodes.BinarySearch(node[Node]{
		hash: hash,
		Node: key,
	}); ok {
		return key
	} else {
		if i == hr.nodes.Len() {
			n, _ := hr.nodes.First()
			return n.Node
		} else {
			n, _ := hr.nodes.Get(i)
			return n.Node
		}
	}
}

func (hr *HashRing[Node]) hashKey(key Node, index int) uint32 {
	return crc32.ChecksumIEEE([]byte(fmt.Sprintf("%v-%d", key, index)))
}

func (hr *HashRing[Node]) RemoveNodes(nodes ...Node) {
	if len(nodes) == 0 {
		return
	}
	s := hr.rawNodes.Intersection(set.New[Node]().Push(nodes...))
	hr.ResetNodes(hr.rawNodes.Difference(s).Slice()...)
}
