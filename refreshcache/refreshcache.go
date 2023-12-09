package refreshcache

import (
	"sync"
	"sync/atomic"
	"time"
)

type RefreshCache[T any] struct {
	data        *RefreshData[T]
	refreshFunc func() (T, error)
}

func NewRefreshCache[T any](refreshFunc func() (T, error), maxAge time.Duration) *RefreshCache[T] {
	if refreshFunc == nil {
		panic("refreshFunc cannot be nil")
	}
	return &RefreshCache[T]{
		refreshFunc: refreshFunc,
		data:        NewRefreshData[T](maxAge),
	}
}

func (r *RefreshCache[T]) Get() (data T, err error) {
	return r.data.Get(r.refreshFunc)
}

func (r *RefreshCache[T]) Refresh() (data T, err error) {
	return r.data.Refresh(r.refreshFunc)
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

func (r *RefreshData[T]) Get(refreshFunc func() (T, error)) (data T, err error) {
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
	return refreshFunc()
}

func (r *RefreshData[T]) Refresh(refreshFunc func() (T, error)) (data T, err error) {
	r.lock.Lock()
	defer r.lock.Unlock()
	defer func() {
		if err == nil {
			r.data.Store(&data)
			atomic.StoreInt64(&r.last, time.Now().UnixNano())
		}
	}()
	return refreshFunc()
}
