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
	rawNoods set.Set[Node]
	noods    *vec.Vec[node[Node]]
}

type node[Node constraints.Ordered] struct {
	Node Node
	hash uint32
}

func New[Node constraints.Ordered](replicas int) *HashRing[Node] {
	hr := &HashRing[Node]{
		replicas: replicas,
		rawNoods: set.New[Node](),
		noods: vec.New[node[Node]](
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
	s := set.New[Node]().Push(nodes...).Difference(hr.rawNoods)
	s.Range(func(val Node) (Continue bool) {
		for v := 0; v < hr.replicas; v++ {
			hr.noods.Push(node[Node]{
				Node: val,
				hash: hr.hashKey(val, v),
			})
		}
		return true
	})
	hr.rawNoods.Push(s.Slice()...)
	hr.noods.Sort()
	return hr
}

func (hr *HashRing[Node]) ResetNodes(nodes ...Node) {
	hr.noods.Clear()
	hr.rawNoods.Clear()
	hr.AddNodes(nodes...)
}

// GetNode returns the node that a given key maps to
func (hr *HashRing[Node]) GetNode(key Node) (n Node) {
	if hr.noods.Len() == 0 {
		return
	}
	hash := hr.hashKey(key, 0)
	if i, ok := hr.noods.BinarySearch(node[Node]{
		hash: hash,
		Node: key,
	}); ok {
		return key
	} else {
		if i == hr.noods.Len() {
			n, _ := hr.noods.First()
			return n.Node
		} else {
			n, _ := hr.noods.Get(i)
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
	s := hr.rawNoods.Intersection(set.New[Node]().Push(nodes...))
	hr.ResetNodes(hr.rawNoods.Difference(s).Slice()...)
}
