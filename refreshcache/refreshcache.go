package refreshcache

import (
	"context"
	"sync"
	"sync/atomic"
	"time"
)

type RefreshFunc[T any, A any] func(ctx context.Context, args ...A) (T, error)

type RefreshCache[T any, A any] struct {
	refreshData *RefreshData[T, A]
	refreshFunc atomic.Pointer[RefreshFunc[T, A]]
}

func NewRefreshCache[T any, A any](refreshFunc RefreshFunc[T, A], maxAge time.Duration) *RefreshCache[T, A] {
	c := RefreshCache[T, A]{
		refreshData: NewRefreshData[T, A](maxAge),
	}
	c.SetRefreshFunc(refreshFunc)
	return &c
}

func (r *RefreshCache[T, A]) GetRefreshFunc() RefreshFunc[T, A] {
	return *r.refreshFunc.Load()
}

func (r *RefreshCache[T, A]) SetRefreshFunc(refreshFunc RefreshFunc[T, A]) {
	r.refreshFunc.Store(&refreshFunc)
}

func (r *RefreshCache[T, A]) Get(ctx context.Context, args ...A) (data T, err error) {
	return r.refreshData.Get(ctx, *r.refreshFunc.Load(), args...)
}

func (r *RefreshCache[T, A]) Refresh(ctx context.Context, args ...A) (data T, err error) {
	return r.refreshData.Refresh(ctx, *r.refreshFunc.Load(), args...)
}

func (r *RefreshCache[T, A]) Clear() {
	r.refreshData.Clear()
}

func (r *RefreshCache[T, A]) Last() int64 {
	return r.refreshData.Last()
}

func (r *RefreshCache[T, A]) LastTime() time.Time {
	return r.refreshData.LastTime()
}

func (r *RefreshCache[T, A]) MaxAge() int64 {
	return r.refreshData.MaxAge()
}

func (r *RefreshCache[T, A]) Raw() T {
	return r.refreshData.Raw()
}

func (r *RefreshCache[T, A]) Data() *RefreshData[T, A] {
	return r.refreshData
}

type RefreshData[T any, A any] struct {
	last   int64
	maxAge int64
	lock   sync.Mutex
	data   atomic.Pointer[T]
}

func NewRefreshData[T any, A any](maxAge time.Duration) *RefreshData[T, A] {
	return &RefreshData[T, A]{
		maxAge: int64(maxAge),
	}
}

func (r *RefreshData[T, A]) Get(ctx context.Context, refreshFunc RefreshFunc[T, A], args ...A) (data T, err error) {
	if (r.maxAge <= 0 && atomic.LoadInt64(&r.last) > 0) || (time.Now().UnixNano()-atomic.LoadInt64(&r.last) < r.maxAge) {
		return *r.data.Load(), nil
	}
	r.lock.Lock()
	defer r.lock.Unlock()
	if (r.maxAge <= 0 && r.last > 0) || (time.Now().UnixNano()-r.last < r.maxAge) {
		return *r.data.Load(), nil
	}
	defer func() {
		if err == nil {
			r.data.Store(&data)
			atomic.StoreInt64(&r.last, time.Now().UnixNano())
		}
	}()
	return refreshFunc(ctx, args...)
}

func (r *RefreshData[T, A]) Refresh(ctx context.Context, refreshFunc RefreshFunc[T, A], args ...A) (data T, err error) {
	r.lock.Lock()
	defer r.lock.Unlock()
	defer func() {
		if err == nil {
			r.data.Store(&data)
			atomic.StoreInt64(&r.last, time.Now().UnixNano())
		}
	}()
	return refreshFunc(ctx, args...)
}

func (r *RefreshData[T, A]) Clear() {
	atomic.StoreInt64(&r.last, 0)
}

func (r *RefreshData[T, A]) Last() int64 {
	return atomic.LoadInt64(&r.last)
}

func (r *RefreshData[T, A]) LastTime() time.Time {
	return time.Unix(0, atomic.LoadInt64(&r.last))
}

func (r *RefreshData[T, A]) MaxAge() int64 {
	return r.maxAge
}

func (r *RefreshData[T, A]) Raw() T {
	return *r.data.Load()
}
