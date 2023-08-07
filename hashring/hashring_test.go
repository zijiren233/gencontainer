package hashring

import (
	"testing"
)

func TestGetNode(t *testing.T) {
	hr := New[string](1024, WithNodes("node1", "node2", "node3", "node4"))
	if hr.rawNoods.Len() != 4 {
		t.Errorf("New error %d", hr.rawNoods.Len())
	}
	if hr.sortedNodes.Len() != 1024*4 {
		t.Errorf("New error %d", hr.sortedNodes.Len())
	}
	s := hr.GetNode("somehash to get node")
	if s == "" {
		t.Errorf("GetNode error %s", s)
	}
}

func TestAddNode(t *testing.T) {
	hr := New[string](10, WithNodes("node1"))
	hr.AddNodes("node2")
	if hr.sortedNodes.Len() != 20 {
		t.Error("AddNodes error")
	}
	hr.RemoveNodes("node2")
	if hr.sortedNodes.Len() != 10 {
		t.Error("RemoveNodes error")
	}
}
