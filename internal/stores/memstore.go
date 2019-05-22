package stores

import (
	"fmt"
	"sync"
	"time"

	interfaces "github.com/shoplineapp/captin/interfaces"
	"github.com/shoplineapp/captin/models"
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
func (ms *MemoryStore) Get(key string) (string, bool, time.Duration, error) {
	ms.lock.Lock()
	defer ms.lock.Unlock()
	fmt.Println("[Memstore] Get Key: ", key)
	if it, ok := ms.m[key]; ok {
		return it.value, true, time.Since(it.createDate), nil
	}
	return "", false, 0, nil
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
	it.ttl = ttl
	ms.lock.Unlock()
	return true, nil
}

// Update - Update value into store
func (ms *MemoryStore) Update(key string, value string) (bool, error) {
	ms.lock.Lock()
	defer ms.lock.Unlock()

	it, ok := ms.m[key]
	if !ok {
		return false, nil
	}
	it.value = value
	ms.m[key] = it
	return true, nil
}

// Remove - Remove value in store
func (ms *MemoryStore) Remove(key string) (bool, error) {
	ms.lock.Lock()
	delete(ms.m, key)
	ms.lock.Unlock()
	return true, nil
}

// Len - Get memory size
func (ms *MemoryStore) Len() int {
	return len(ms.m)
}

// DataKey - Generate DataKey with events and destination
func (ms *MemoryStore) DataKey(e models.IncomingEvent, dest models.Destination, prefix string, suffix string) string {
	return fmt.Sprintf("%s%s.%s.%s%s", prefix, e.Key, dest.Config.Name, e.TargetId, suffix)
}
