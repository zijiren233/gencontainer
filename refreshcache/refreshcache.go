package refreshcache

import (
	"context"
	"sync"
	"sync/atomic"
	"time"
)

type RefreshCache[T any] struct {
	data        *RefreshData[T]
	refreshFunc func(context.Context) (T, error)
}

func NewRefreshCache[T any](refreshFunc func(context.Context) (T, error), maxAge time.Duration) *RefreshCache[T] {
	return &RefreshCache[T]{
		refreshFunc: refreshFunc,
		data:        NewRefreshData[T](maxAge),
	}
}

func (r *RefreshCache[T]) Get(ctx context.Context) (data T, err error) {
	return r.data.Get(ctx, r.refreshFunc)
}

func (r *RefreshCache[T]) Refresh(ctx context.Context) (data T, err error) {
	return r.data.Refresh(ctx, r.refreshFunc)
}

func (r *RefreshCache[T]) Clear() {
	r.data.Clear()
}

func (r *RefreshCache[T]) Last() int64 {
	return r.data.Last()
}

func (r *RefreshCache[T]) MaxAge() int64 {
	return r.data.MaxAge()
}

func (r *RefreshCache[T]) Raw() T {
	return r.data.Raw()
}

func (r *RefreshCache[T]) Data() *RefreshData[T] {
	return r.data
}

type RefreshData[T any] struct {
	lock   sync.Mutex
	last   int64
	maxAge int64
	data   atomic.Pointer[T]
}

func NewRefreshData[T any](maxAge time.Duration) *RefreshData[T] {
	return &RefreshData[T]{
		maxAge: int64(maxAge),
	}
}

func (r *RefreshData[T]) Get(ctx context.Context, refreshFunc func(context.Context) (T, error)) (data T, err error) {
	if (r.maxAge <= 0 && atomic.LoadInt64(&r.last) > 0) || (time.Now().UnixMicro()-atomic.LoadInt64(&r.last) < r.maxAge) {
		return *r.data.Load(), nil
	}
	r.lock.Lock()
	defer r.lock.Unlock()
	if (r.maxAge <= 0 && r.last > 0) || (time.Now().UnixMicro()-r.last < r.maxAge) {
		return *r.data.Load(), nil
	}
	defer func() {
		if err == nil {
			r.data.Store(&data)
			atomic.StoreInt64(&r.last, time.Now().UnixMicro())
		}
	}()
	return refreshFunc(ctx)
}

func (r *RefreshData[T]) Refresh(ctx context.Context, refreshFunc func(context.Context) (T, error)) (data T, err error) {
	r.lock.Lock()
	defer r.lock.Unlock()
	defer func() {
		if err == nil {
			r.data.Store(&data)
			atomic.StoreInt64(&r.last, time.Now().UnixMicro())
		}
	}()
	return refreshFunc(ctx)
}

func (r *RefreshData[T]) Clear() {
	atomic.StoreInt64(&r.last, 0)
}

func (r *RefreshData[T]) Last() int64 {
	return atomic.LoadInt64(&r.last)
}

func (r *RefreshData[T]) MaxAge() int64 {
	return r.maxAge
}

func (r *RefreshData[T]) Raw() T {
	return *r.data.Load()
}
