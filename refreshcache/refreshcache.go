package refreshcache

import (
	"context"
	"sync"
	"sync/atomic"
	"time"
)

type RefreshFunc[T any, A any] func(ctx context.Context, args ...A) (T, error)

type RefreshCache[T any, A any] struct {
	*RefreshData[T, A]
	RefreshFunc atomic.Pointer[RefreshFunc[T, A]]
}

func NewRefreshCache[T any, A any](refreshFunc RefreshFunc[T, A], maxAge time.Duration, opts ...RefreshDataOption[T, A]) *RefreshCache[T, A] {
	c := RefreshCache[T, A]{
		RefreshData: NewRefreshData[T, A](maxAge, opts...),
	}
	c.SetRefreshFunc(refreshFunc)
	return &c
}

func (r *RefreshCache[T, A]) GetRefreshFunc() RefreshFunc[T, A] {
	return *r.RefreshFunc.Load()
}

func (r *RefreshCache[T, A]) SetRefreshFunc(refreshFunc RefreshFunc[T, A]) {
	r.RefreshFunc.Store(&refreshFunc)
}

func (r *RefreshCache[T, A]) Get(ctx context.Context, args ...A) (data T, err error) {
	return r.RefreshData.Get(ctx, *r.RefreshFunc.Load(), args...)
}

func (r *RefreshCache[T, A]) Refresh(ctx context.Context, args ...A) (data T, err error) {
	return r.RefreshData.Refresh(ctx, *r.RefreshFunc.Load(), args...)
}

func (r *RefreshCache[T, A]) Data() *RefreshData[T, A] {
	return r.RefreshData
}

type RefreshData[T any, A any] struct {
	last    int64
	maxAge  int64
	lastErr int64
	errAge  int64
	err     atomic.Pointer[error]
	lock    sync.Mutex
	data    atomic.Pointer[T]
}

type RefreshDataOption[T any, A any] func(*RefreshData[T, A])

func WithErrAge[T any, A any](age time.Duration) RefreshDataOption[T, A] {
	return func(r *RefreshData[T, A]) {
		r.errAge = int64(age)
	}
}

func NewRefreshData[T any, A any](maxAge time.Duration, opts ...RefreshDataOption[T, A]) *RefreshData[T, A] {
	rd := &RefreshData[T, A]{
		maxAge: int64(maxAge),
	}
	rd.data.Store(new(T))
	rd.err.Store(new(error))
	for _, opt := range opts {
		opt(rd)
	}
	return rd
}

type oldVal struct{}

var OldValKey = oldVal{}

func (r *RefreshData[T, A]) Get(ctx context.Context, refreshFunc RefreshFunc[T, A], args ...A) (data T, err error) {
	last := atomic.LoadInt64(&r.last)
	if (r.maxAge < 0 && last > 0) || (time.Now().UnixNano()-last < r.maxAge) {
		return *r.data.Load(), nil
	}
	if r.errAge > 0 && (time.Now().UnixNano()-atomic.LoadInt64(&r.lastErr) < r.errAge) {
		return data, *r.err.Load()
	}
	r.lock.Lock()
	defer r.lock.Unlock()
	if (r.maxAge < 0 && r.last > 0) || (time.Now().UnixNano()-r.last < r.maxAge) {
		return *r.data.Load(), nil
	}
	if r.errAge > 0 && (time.Now().UnixNano()-r.lastErr < r.errAge) {
		return data, *r.err.Load()
	}
	defer func() {
		if err == nil {
			r.data.Store(&data)
			atomic.StoreInt64(&r.last, time.Now().UnixNano())
			atomic.StoreInt64(&r.lastErr, 0)
		} else {
			r.err.Store(&err)
			atomic.StoreInt64(&r.lastErr, time.Now().UnixNano())
			atomic.StoreInt64(&r.last, 0)
		}
	}()
	return refreshFunc(context.WithValue(ctx, OldValKey, *r.data.Load()), args...)
}

func (r *RefreshData[T, A]) Refresh(ctx context.Context, refreshFunc RefreshFunc[T, A], args ...A) (data T, err error) {
	r.lock.Lock()
	defer r.lock.Unlock()
	defer func() {
		if err == nil {
			r.data.Store(&data)
			atomic.StoreInt64(&r.last, time.Now().UnixNano())
			atomic.StoreInt64(&r.lastErr, 0)
		} else {
			r.err.Store(&err)
			atomic.StoreInt64(&r.lastErr, time.Now().UnixNano())
			atomic.StoreInt64(&r.last, 0)
		}
	}()
	return refreshFunc(context.WithValue(ctx, OldValKey, *r.data.Load()), args...)
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

func (r *RefreshData[T, A]) LastErr() int64 {
	return atomic.LoadInt64(&r.lastErr)
}

func (r *RefreshData[T, A]) LastErrTime() time.Time {
	return time.Unix(0, atomic.LoadInt64(&r.lastErr))
}

func (r *RefreshData[T, A]) MaxAge() int64 {
	return r.maxAge
}

func (r *RefreshData[T, A]) MaxErrAge() int64 {
	return r.errAge
}

func (r *RefreshData[T, A]) Raw() (data T, err error) {
	last := atomic.LoadInt64(&r.last)
	if (r.maxAge < 0 && last > 0) || (time.Now().UnixNano()-last < r.maxAge) {
		return *r.data.Load(), nil
	}
	if r.errAge > 0 && (time.Now().UnixNano()-atomic.LoadInt64(&r.lastErr) < r.errAge) {
		return data, *r.err.Load()
	}
	return *r.data.Load(), *r.err.Load()
}
