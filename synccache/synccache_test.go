package synccache_test

import (
	"testing"
	"time"

	"github.com/zijiren233/gencontainer/synccache"
)

func TestCompareValueAndDelete(t *testing.T) {
	s := synccache.NewSyncCache[int, string](time.Second)
	s.Store(1, "hello", time.Second*3)
	success := s.CompareValueAndDelete(
		1,
		"hello",
	)
	if !success {
		t.Fatal("compare value and delete failed")
	}
}
