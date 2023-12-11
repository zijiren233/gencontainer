package refreshcache_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/zijiren233/gencontainer/refreshcache"
)

func TestRefreshCache(t *testing.T) {
	c := refreshcache.NewRefreshCache[int](func(context.Context, ...any) (int, error) {
		return time.Now().Second(), nil
	}, time.Second)
	fmt.Println(c.Get(context.Background()))
	fmt.Println(c.Get(context.Background()))
	time.Sleep(time.Second)
	fmt.Println(c.Get(context.Background()))
}

func TestRefreshCacheWithArgs(t *testing.T) {
	c := refreshcache.NewRefreshCache[int](func(ctx context.Context, args ...any) (int, error) {
		if len(args) == 0 {
			return 0, nil
		}
		return args[0].(int) + time.Now().Second(), nil
	}, time.Second)
	fmt.Println(c.Get(context.Background(), 1))
	fmt.Println(c.Refresh(context.Background()))
}

func TestRefreshCacheStatic(t *testing.T) {
	c := refreshcache.NewRefreshCache[int](func(context.Context, ...any) (int, error) {
		return time.Now().Second(), nil
	}, 0)
	fmt.Println(c.Get(context.Background()))
	fmt.Println(c.Get(context.Background()))
	time.Sleep(time.Second)
	fmt.Println(c.Get(context.Background()))
}

func TestRefreshData(t *testing.T) {
	d := refreshcache.NewRefreshData[int](time.Second)
	fmt.Println(d.Get(context.Background(), func(context.Context, ...any) (int, error) {
		return time.Now().Second(), nil
	}))
	fmt.Println(d.Get(context.Background(), func(context.Context, ...any) (int, error) {
		return time.Now().Second(), nil
	}))
	time.Sleep(time.Second)
	fmt.Println(d.Get(context.Background(), func(context.Context, ...any) (int, error) {
		return time.Now().Second(), nil
	}))
}
