package hashring

import (
	"fmt"
	"hash/crc32"
	"sort"

	"github.com/zijiren233/gencontainer/dllist"
	"github.com/zijiren233/gencontainer/restrictions"
	"github.com/zijiren233/gencontainer/vec"
)

type HashRing[Node restrictions.Ordered] struct {
	replicas    int
	rawNoods    *vec.Vec[Node]
	sortedNodes *dllist.Dllist[node[Node]]
}

type node[Node restrictions.Ordered] struct {
	Node Node
	hash uint32
}

type HashRingConf[Node restrictions.Ordered] func(*HashRing[Node])

func WithNodes[Node restrictions.Ordered](nodes ...Node) HashRingConf[Node] {
	return func(hr *HashRing[Node]) {
		hr.AddNodes(nodes...)
	}
}

func New[Node restrictions.Ordered](replicas int, conf ...HashRingConf[Node]) *HashRing[Node] {
	hr := &HashRing[Node]{
		replicas: replicas,
		rawNoods: vec.New[Node](),
		sortedNodes: dllist.New[node[Node]](dllist.WithLessFunc[node[Node]](
			func(a, b node[Node]) bool {
				return a.hash < b.hash
			},
		)),
	}
	for _, c := range conf {
		c(hr)
	}
	return hr
}

func (hr *HashRing[Node]) AddNodes(nodes ...Node) {
	if len(nodes) == 0 {
		return
	}
	hr.rawNoods.Push(nodes...)
	for _, n := range nodes {
		for i := 0; i < hr.replicas; i++ {
			key := hr.hashKey(n, i)
			hr.sortedNodes.PushBack(node[Node]{
				Node: n,
				hash: key,
			})
		}
	}
	hr.sortedNodes.Sort()
}

func (hr *HashRing[Node]) ResetNodes(nodes ...Node) {
	hr.sortedNodes.Clear()
	hr.rawNoods.Clear()
	hr.AddNodes(nodes...)
}

// GetNode returns the node that a given key maps to
func (hr *HashRing[Node]) GetNode(key Node) (n Node) {
	if hr.sortedNodes.Len() == 0 {
		return
	}
	hash := hr.hashKey(key, 0)

	index := sort.Search(hr.sortedNodes.Len(), func(i int) bool {
		return hr.sortedNodes.Get(i).Value.hash >= hash
	})
	if index == hr.sortedNodes.Len() {
		index = 0
	}
	return hr.sortedNodes.Get(index).Value.Node
}

func (hr *HashRing[Node]) hashKey(key Node, index int) uint32 {
	return crc32.ChecksumIEEE([]byte(fmt.Sprintf("%v-%d", key, index)))
}

func (hr *HashRing[Node]) RemoveNodes(nodes ...Node) {
	if len(nodes) == 0 {
		return
	}
	keys := make(map[uint32]struct{})
	for _, n := range nodes {
		for _, v := range hr.rawNoods.FindAll(n) {
			for i := 0; i < hr.replicas; i++ {
				keys[hr.hashKey(n, i)] = struct{}{}
			}
			hr.rawNoods.Remove(v)
		}
	}
	removeds := make([]*dllist.Element[node[Node]], 0, len(keys))
	hr.sortedNodes.Range(func(e *dllist.Element[node[Node]]) bool {
		if _, ok := keys[e.Value.hash]; ok {
			removeds = append(removeds, e)
		}
		return true
	})
	for _, e := range removeds {
		e.Remove()
	}
}
