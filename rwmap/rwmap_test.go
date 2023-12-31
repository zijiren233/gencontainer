package rwmap_test

import (
	"math/rand"
	"reflect"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/zijiren233/gencontainer/rwmap"
)

func TestConcurrentRange(t *testing.T) {
	const mapSize = 1 << 10

	m := new(rwmap.RWMap[int64, int64])
	for n := int64(1); n <= mapSize; n++ {
		m.Store(n, int64(n))
	}

	done := make(chan struct{})
	var wg sync.WaitGroup
	defer func() {
		close(done)
		wg.Wait()
	}()
	for g := int64(runtime.GOMAXPROCS(0)); g > 0; g-- {
		r := rand.New(rand.NewSource(g))
		wg.Add(1)
		go func(g int64) {
			defer wg.Done()
			for i := int64(0); ; i++ {
				select {
				case <-done:
					return
				default:
				}
				for n := int64(1); n < mapSize; n++ {
					if r.Int63n(mapSize) == 0 {
						m.Store(n, n*i*g)
					} else {
						m.Load(n)
					}
				}
			}
		}(g)
	}

	iters := 1 << 10
	if testing.Short() {
		iters = 16
	}
	for n := iters; n > 0; n-- {
		seen := make(map[int64]bool, mapSize)

		m.Range(func(ki int64, vi int64) bool {
			k, v := ki, vi
			if v%k != 0 {
				t.Fatalf("while Storing multiples of %v, Range saw value %v", k, v)
			}
			if seen[k] {
				t.Fatalf("Range visited key %v twice", k)
			}
			seen[k] = true
			return true
		})

		if len(seen) != mapSize {
			t.Fatalf("Range visited %v elements of %v-element Map", len(seen), mapSize)
		}
	}
}

func TestIssue40999(t *testing.T) {
	var m rwmap.RWMap[*int, struct{}]

	// Since the miss-counting in missLocked (via Delete)
	// compares the miss count with len(m.dirty),
	// add an initial entry to bias len(m.dirty) above the miss count.
	m.Store(nil, struct{}{})

	var finalized uint32

	// Set finalizers that count for collected keys. A non-zero count
	// indicates that keys have not been leaked.
	for atomic.LoadUint32(&finalized) == 0 {
		p := new(int)
		runtime.SetFinalizer(p, func(*int) {
			atomic.AddUint32(&finalized, 1)
		})
		m.Store(p, struct{}{})
		m.Delete(p)
		runtime.GC()
	}
}

func TestMapRangeNestedCall(t *testing.T) { // Issue 46399
	var m rwmap.RWMap[int, string]
	for i, v := range [3]string{"hello", "world", "Go"} {
		m.Store(i, v)
	}
	m.Range(func(key int, value string) bool {
		m.Range(func(key int, value string) bool {
			// We should be able to load the key offered in the Range callback,
			// because there are no concurrent Delete involved in this tested map.
			if v, ok := m.Load(key); !ok || !reflect.DeepEqual(v, value) {
				t.Fatalf("Nested Range loads unexpected value, got %+v want %+v", v, value)
			}

			// We didn't keep 42 and a value into the map before, if somehow we loaded
			// a value from such a key, meaning there must be an internal bug regarding
			// nested range in the Map.
			if _, loaded := m.LoadOrStore(42, "dummy"); loaded {
				t.Fatalf("Nested Range loads unexpected value, want store a new value")
			}

			// Try to Store then LoadAndDelete the corresponding value with the key
			// 42 to the Map. In this case, the key 42 and associated value should be
			// removed from the Map. Therefore any future range won't observe key 42
			// as we checked in above.
			val := "sync.Map"
			m.Store(42, val)
			if v, loaded := m.LoadAndDelete(42); !loaded || !reflect.DeepEqual(v, val) {
				t.Fatalf("Nested Range loads unexpected value, got %v, want %v", v, val)
			}
			return true
		})

		// Remove key from Map on-the-fly.
		m.Delete(key)
		return true
	})

	// After a Range of Delete, all keys should be removed and any
	// further Range won't invoke the callback. Hence length remains 0.
	length := 0
	m.Range(func(key int, value string) bool {
		length++
		return true
	})

	if length != 0 {
		t.Fatalf("Unexpected sync.Map size, got %v want %v", length, 0)
	}
}

func TestCompareAndSwap_NonExistingKey(t *testing.T) {
	m := &rwmap.RWMap[any, any]{}
	if m.CompareAndSwap(m, nil, 42) {
		// See https://go.dev/issue/51972#issuecomment-1126408637.
		t.Fatalf("CompareAndSwap on an non-existing key succeeded")
	}
}

func TestLoadOrStore(t *testing.T) {
	m := &rwmap.RWMap[any, any]{}
	if v, ok := m.LoadOrStore(m, 42); ok {
		t.Fatalf("LoadOrStore on an non-existing key succeeded, got %v", v)
	}
	if v, ok := m.LoadOrStore(m, 42); !ok || v != 42 {
		t.Fatalf("LoadOrStore on an existing key failed, got %v", v)
	}
}

func TestClear(t *testing.T) {
	m := &rwmap.RWMap[any, any]{}
	m.Store(m, 42)
	m.Clear()
	if v, ok := m.Load(m); ok {
		t.Fatalf("Clear failed, got %v", v)
	}
	m.Store(m, 42)
	if v, ok := m.Load(m); !ok || v != 42 {
		t.Fatalf("Clear failed, got %v", v)
	}
}

func TestCompareAndSwap(t *testing.T) {
	m := &rwmap.RWMap[any, any]{}
	m.Store(m, 42)
	if m.CompareAndSwap(m, nil, 42) {
		t.Fatalf("CompareAndSwap on an non-existing key succeeded")
	}
	if !m.CompareAndSwap(m, 42, 43) {
		t.Fatalf("CompareAndSwap on an existing key failed")
	}
	if v, ok := m.Load(m); !ok || v != 43 {
		t.Fatalf("CompareAndSwap failed, got %v", v)
	}
}

func TestLen(t *testing.T) {
	m := &rwmap.RWMap[any, any]{}
	if m.Len() != 0 {
		t.Fatalf("Len failed, got %v", m.Len())
	}
	m.Store(m, 42)
	if m.Len() != 1 {
		t.Fatalf("Len failed, got %v", m.Len())
	}
	m.Delete(m)
	if m.Len() != 0 {
		t.Fatalf("Len failed, got %v", m.Len())
	}
	m.LoadOrStore(m, 1)
	if m.Len() != 1 {
		t.Fatalf("Len failed, got %v", m.Len())
	}
	m.CompareAndSwap(m, nil, 42)
	if m.Len() != 1 {
		t.Fatalf("Len failed, got %v", m.Len())
	}
	m.CompareAndSwap(m, 1, 2)
	if m.Len() != 1 {
		t.Fatalf("Len failed, got %v", m.Len())
	}
	m.CompareAndDelete(m, nil)
	if m.Len() != 1 {
		t.Fatalf("Len failed, got %v", m.Len())
	}
	m.Clear()
	if m.Len() != 0 {
		t.Fatalf("Len failed, got %v", m.Len())
	}
}
