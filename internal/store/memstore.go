package store

import (
	"sync"
	"time"

	interfaces "github.com/shoplineapp/captin/interfaces"
)

type item struct {
	value      string
	createDate time.Time
	ttl        time.Duration
}

// MemoryStore - In-app memory storage
type MemoryStore struct {
	interfaces.StoreInterface
	m    map[string]*item
	lock sync.Mutex
}

// NewMemoryStore - Create new MemoryStore
func NewMemoryStore() *MemoryStore {
	m := &MemoryStore{
		m: make(map[string]*item),
	}

	go func() {
		for range time.Tick(100 * time.Millisecond) {
			m.lock.Lock()
			for k, v := range m.m {
				duration := time.Since(v.createDate)
				if duration > v.ttl {
					delete(m.m, k)
				}
			}
			m.lock.Unlock()
		}
	}()

	return m
}

// Get - Get value from store, return with remaining time
func (ms *MemoryStore) Get(key string) (string, time.Duration, error) {
	ms.lock.Lock()
	defer ms.lock.Unlock()
	if it, ok := ms.m[key]; ok {
		return it.value, time.Since(it.createDate), nil
	}
	return "", 0, nil
}

// Set - Set value into store with ttl
func (ms *MemoryStore) Set(key string, value string, ttl time.Duration) (bool, error) {
	ms.lock.Lock()
	it, ok := ms.m[key]
	if !ok {
		it = &item{value: value}
		ms.m[key] = it
	}
	it.createDate = time.Now()
	ms.lock.Unlock()
	return true, nil
}

// Len - Get memory size
func (ms *MemoryStore) Len() int {
	return len(ms.m)
}
