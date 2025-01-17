package refreshcache0

import (
	"context"
	"sync"
	"sync/atomic"
	"time"
)

type RefreshFunc[T any] func(ctx context.Context) (T, error)

type ClearFunc[T any] func(ctx context.Context) error

type RefreshCache[T any] struct {
	*RefreshData[T]
	RefreshFunc atomic.Pointer[RefreshFunc[T]]
	ClearFunc   atomic.Pointer[ClearFunc[T]]
}

func NewRefreshCache[T any](refreshFunc RefreshFunc[T], maxAge time.Duration, opts ...RefreshDataOption[T]) *RefreshCache[T] {
	c := RefreshCache[T]{
		RefreshData: NewRefreshData[T](maxAge, opts...),
	}
	c.SetRefreshFunc(refreshFunc)
	return &c
}

func (r *RefreshCache[T]) GetClearFunc() ClearFunc[T] {
	p := r.ClearFunc.Load()
	if p == nil {
		return nil
	}
	return *p
}

func (r *RefreshCache[T]) SetClearFunc(clearFunc ClearFunc[T]) {
	r.ClearFunc.Store(&clearFunc)
}

func (r *RefreshCache[T]) GetRefreshFunc() RefreshFunc[T] {
	return *r.RefreshFunc.Load()
}

func (r *RefreshCache[T]) SetRefreshFunc(refreshFunc RefreshFunc[T]) {
	r.RefreshFunc.Store(&refreshFunc)
}

func (r *RefreshCache[T]) Get(ctx context.Context) (data T, err error) {
	return r.RefreshData.Get(ctx, *r.RefreshFunc.Load())
}

func (r *RefreshCache[T]) Refresh(ctx context.Context) (data T, err error) {
	return r.RefreshData.Refresh(ctx, *r.RefreshFunc.Load())
}

func (r *RefreshCache[T]) Data() *RefreshData[T] {
	return r.RefreshData
}

func (r *RefreshCache[T]) Clear(ctx context.Context) error {
	return r.RefreshData.Clear(ctx, r.GetClearFunc())
}

type RefreshData[T any] struct {
	err     atomic.Pointer[error]
	data    atomic.Pointer[T]
	last    int64
	maxAge  int64
	lastErr int64
	errAge  int64
	lock    sync.Mutex
}

type RefreshDataOption[T any] func(*RefreshData[T])

func WithErrAge[T any](age time.Duration) RefreshDataOption[T] {
	return func(r *RefreshData[T]) {
		r.errAge = int64(age)
	}
}

func NewRefreshData[T any](maxAge time.Duration, opts ...RefreshDataOption[T]) *RefreshData[T] {
	rd := &RefreshData[T]{
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

func (r *RefreshData[T]) Get(ctx context.Context, refreshFunc RefreshFunc[T]) (data T, err error) {
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
	return refreshFunc(context.WithValue(ctx, OldValKey, *r.data.Load()))
}

func (r *RefreshData[T]) Refresh(ctx context.Context, refreshFunc RefreshFunc[T]) (data T, err error) {
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
	return refreshFunc(context.WithValue(ctx, OldValKey, *r.data.Load()))
}

func (r *RefreshData[T]) Clear(ctx context.Context, clearFunc ClearFunc[T]) error {
	if clearFunc != nil {
		data := r.data.Load()
		atomic.StoreInt64(&r.last, 0)
		return clearFunc(context.WithValue(ctx, OldValKey, *data))
	}
	atomic.StoreInt64(&r.last, 0)
	return nil
}

func (r *RefreshData[T]) Last() int64 {
	return atomic.LoadInt64(&r.last)
}

func (r *RefreshData[T]) LastTime() time.Time {
	return time.Unix(0, atomic.LoadInt64(&r.last))
}

func (r *RefreshData[T]) LastErr() int64 {
	return atomic.LoadInt64(&r.lastErr)
}

func (r *RefreshData[T]) LastErrTime() time.Time {
	return time.Unix(0, atomic.LoadInt64(&r.lastErr))
}

func (r *RefreshData[T]) MaxAge() int64 {
	return r.maxAge
}

func (r *RefreshData[T]) MaxErrAge() int64 {
	return r.errAge
}

func (r *RefreshData[T]) Raw() (data T, err error) {
	last := atomic.LoadInt64(&r.last)
	if (r.maxAge < 0 && last > 0) || (time.Now().UnixNano()-last < r.maxAge) {
		return *r.data.Load(), nil
	}
	if r.errAge > 0 && (time.Now().UnixNano()-atomic.LoadInt64(&r.lastErr) < r.errAge) {
		return data, *r.err.Load()
	}
	return *r.data.Load(), *r.err.Load()
}
