package refreshcache_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/zijiren233/gencontainer/refreshcache"
)

func TestRefreshCache(t *testing.T) {
	c := refreshcache.NewRefreshCache[int](func() (int, error) {
		return time.Now().Second(), nil
	}, time.Second)
	fmt.Println(c.Get())
	fmt.Println(c.Get())
	time.Sleep(time.Second)
	fmt.Println(c.Get())
}

func TestRefreshData(t *testing.T) {
	d := refreshcache.NewRefreshData[int](time.Second)
	fmt.Println(d.Get(func() (int, error) {
		return time.Now().Second(), nil
	}))
	fmt.Println(d.Get(func() (int, error) {
		return time.Now().Second(), nil
	}))
	time.Sleep(time.Second)
	fmt.Println(d.Get(func() (int, error) {
		return time.Now().Second(), nil
	}))
}
