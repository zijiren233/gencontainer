package hashring

import (
	"testing"
)

func TestGetNode(t *testing.T) {
	hr := New[string](1024).AddNodes("node1", "node2", "node3", "node4")
	if hr.rawNodes.Len() != 4 {
		t.Errorf("New error %d", hr.rawNodes.Len())
	}
	if hr.nodes.Len() != 1024*4 {
		t.Errorf("New error %d", hr.nodes.Len())
	}
	s := hr.GetNode("somehash to get node")
	if s == "" {
		t.Errorf("GetNode error %s", s)
	}
}

func TestAddNode(t *testing.T) {
	hr := New[string](10).AddNodes("node1")
	hr.AddNodes("node2")
	if hr.nodes.Len() != 20 {
		t.Error("AddNodes error")
	}
	hr.RemoveNodes("node2")
	if hr.nodes.Len() != 10 {
		t.Error("RemoveNodes error")
	}
}
