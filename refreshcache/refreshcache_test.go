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
	c := refreshcache.NewRefreshCache[int](func(ctx context.Context, args ...int) (int, error) {
		if len(args) == 0 {
			return 0, nil
		}
		return args[0] + time.Now().Second(), nil
	}, time.Second)
	fmt.Println(c.Get(context.Background(), 1))
	fmt.Println(c.Refresh(context.Background()))
}

func TestRefreshCacheStatic(t *testing.T) {
	c := refreshcache.NewRefreshCache[int](func(context.Context, ...int) (int, error) {
		return time.Now().Second(), nil
	}, 0)
	fmt.Println(c.Get(context.Background()))
	fmt.Println(c.Get(context.Background()))
	time.Sleep(time.Second)
	fmt.Println(c.Get(context.Background()))
}

func TestRefreshData(t *testing.T) {
	d := refreshcache.NewRefreshData[int, any](time.Second)
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
	fmt.Printf("LastTime: %v\n", d.LastTime())
}

func TestCacheError(t *testing.T) {
	c := refreshcache.NewRefreshCache[int](func(context.Context, ...any) (int, error) {
		fmt.Println("refresh")
		return 0, fmt.Errorf("error")
	}, time.Second, refreshcache.WithErrAge[int, any](time.Second*3))
	fmt.Println(c.Get(context.Background()))
	fmt.Println(c.Get(context.Background()))
	time.Sleep(time.Second)
	fmt.Println(c.Get(context.Background()))
	time.Sleep(time.Second * 2)
	fmt.Println(c.Get(context.Background()))
}

func TestRaw(t *testing.T) {
	c := refreshcache.NewRefreshCache[int](func(context.Context, ...any) (int, error) {
		return 3, nil
	}, time.Second)
	fmt.Println(c.Raw())
	fmt.Println(c.Refresh(context.Background()))
	fmt.Println(c.Raw())
	fmt.Println(c.Data().Refresh(context.Background(), func(ctx context.Context, args ...any) (int, error) {
		return 4, fmt.Errorf("error")
	}))
	fmt.Println(c.Raw())
	fmt.Println(c.Data().Refresh(context.Background(), func(ctx context.Context, args ...any) (int, error) {
		return 5, nil
	}))
	fmt.Println(c.Raw())
}

func TestOldVal(t *testing.T) {
	c := refreshcache.NewRefreshCache[int](func(ctx context.Context, _ ...any) (int, error) {
		t.Log(ctx.Value(refreshcache.OldValKey))
		return 3, nil
	}, time.Second)
	_, _ = c.Get(context.Background())
	_, _ = c.Refresh(context.Background())
}
