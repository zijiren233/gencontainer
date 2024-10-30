package synccache

import (
	"reflect"
	"sync"
	"sync/atomic"
	"time"

	"github.com/zijiren233/gencontainer/rwmap"
)

type trim interface {
	trim()
}

var (
	caches    = rwmap.RWMap[trim, struct{}]{}
	startOnce sync.Once
)

func runTrim() {
	startOnce.Do(func() {
		go func() {
			ticker := time.NewTicker(time.Second * 5)
			defer ticker.Stop()
			for range ticker.C {
				caches.Range(func(key trim, _ struct{}) bool {
					key.trim()
					return true
				})
			}
		}()
	})
}

type SyncCache[K comparable, V any] struct {
	deletedCallback func(v V)
	cache           rwmap.RWMap[K, *Entry[V]]
}

type SyncCacheConfig[K comparable, V any] func(sc *SyncCache[K, V])

func WithDeletedCallback[K comparable, V any](callback func(v V)) SyncCacheConfig[K, V] {
	return func(sc *SyncCache[K, V]) {
		sc.deletedCallback = callback
	}
}

func NewSyncCache[K comparable, V any](trimTime time.Duration, conf ...SyncCacheConfig[K, V]) *SyncCache[K, V] {
	sc := &SyncCache[K, V]{}
	for _, c := range conf {
		c(sc)
	}
	sc.Start()
	return sc
}

func (sc *SyncCache[K, V]) Start() {
	caches.Store(sc, struct{}{})
	runTrim()
}

func (sc *SyncCache[K, V]) Close() {
	caches.Delete(sc)
}

func (sc *SyncCache[K, V]) trim() {
	now := time.Now().UnixMilli()
	sc.cache.Range(func(key K, value *Entry[V]) bool {
		if now > atomic.LoadInt64(&value.expiration) {
			sc.CompareAndDelete(key, value)
		}
		return true
	})
}

func (sc *SyncCache[K, V]) Store(key K, value V, expire time.Duration) {
	sc.cache.Store(key, NewEntry[V](value, expire))
}

func (sc *SyncCache[K, V]) Load(key K) (value *Entry[V], loaded bool) {
	e, ok := sc.cache.Load(key)
	if !ok {
		return nil, false
	}
	if !e.IsExpired() {
		return e, true
	}
	sc.CompareAndDelete(key, e)
	return nil, false
}

func (sc *SyncCache[K, V]) LoadOrStore(key K, value V, expire time.Duration) (actual *Entry[V], loaded bool) {
	for {
		e, loaded := sc.cache.LoadOrStore(key, NewEntry[V](value, expire))
		if !loaded {
			return e, false
		}
		if !e.IsExpired() {
			return e, true
		}
		if sc.CompareAndDelete(key, e) {
			continue
		}
		return e, true
	}
}

func (sc *SyncCache[K, V]) Delete(key K) {
	if e, ok := sc.cache.LoadAndDelete(key); ok && sc.deletedCallback != nil {
		sc.deletedCallback(e.value)
	}
}

func (sc *SyncCache[K, V]) LoadAndDelete(key K) (value *Entry[V], loaded bool) {
	value, loaded = sc.cache.LoadAndDelete(key)
	if !loaded {
		return nil, false
	}
	if value.IsExpired() {
		if sc.deletedCallback != nil {
			sc.deletedCallback(value.value)
		}
		return nil, false
	}
	if sc.deletedCallback != nil {
		sc.deletedCallback(value.value)
	}
	return value, true
}

func (sc *SyncCache[K, V]) CompareAndDelete(key K, oldEntry *Entry[V]) (success bool) {
	if success = sc.cache.CompareAndDelete(key, oldEntry); success && sc.deletedCallback != nil {
		sc.deletedCallback(oldEntry.value)
	}
	return
}

func (sc *SyncCache[K, V]) CompareValueAndDelete(key K, oldValue V) (success bool) {
	if e, ok := sc.Load(key); ok && reflect.ValueOf(oldValue).Equal(reflect.ValueOf(e.value)) {
		return sc.CompareAndDelete(key, e)
	}
	return false
}

func (sc *SyncCache[K, V]) Clear() {
	if sc.deletedCallback == nil {
		sc.cache.Clear()
		return
	}
	sc.cache.Range(func(key K, value *Entry[V]) bool {
		sc.CompareAndDelete(key, value)
		return true
	})
}

func (sc *SyncCache[K, V]) Range(f func(key K, value *Entry[V]) bool) {
	sc.cache.Range(func(key K, value *Entry[V]) bool {
		if !value.IsExpired() {
			return f(key, value)
		}
		sc.CompareAndDelete(key, value)
		return true
	})
}
