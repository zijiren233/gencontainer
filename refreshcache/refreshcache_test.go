package refreshcache_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/zijiren233/gencontainer/refreshcache"
)

func TestRefreshCache(t *testing.T) {
	c := refreshcache.NewRefreshCache[int](func(context.Context) (int, error) {
		return time.Now().Second(), nil
	}, time.Second)
	fmt.Println(c.Get(context.Background()))
	fmt.Println(c.Get(context.Background()))
	time.Sleep(time.Second)
	fmt.Println(c.Get(context.Background()))
}

func TestRefreshCacheStatic(t *testing.T) {
	c := refreshcache.NewRefreshCache[int](func(context.Context) (int, error) {
		return time.Now().Second(), nil
	}, 0)
	fmt.Println(c.Get(context.Background()))
	fmt.Println(c.Get(context.Background()))
	time.Sleep(time.Second)
	fmt.Println(c.Get(context.Background()))
}

func TestRefreshData(t *testing.T) {
	d := refreshcache.NewRefreshData[int](time.Second)
	fmt.Println(d.Get(context.Background(), func(context.Context) (int, error) {
		return time.Now().Second(), nil
	}))
	fmt.Println(d.Get(context.Background(), func(context.Context) (int, error) {
		return time.Now().Second(), nil
	}))
	time.Sleep(time.Second)
	fmt.Println(d.Get(context.Background(), func(context.Context) (int, error) {
		return time.Now().Second(), nil
	}))
}
