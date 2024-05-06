package stores

import (
	"context"
	"fmt"
	"sync"
	"time"

	interfaces "github.com/shoplineapp/captin/v2/interfaces"
	"github.com/shoplineapp/captin/v2/models"
	log "github.com/sirupsen/logrus"
)

var mLogger = log.WithFields(log.Fields{"class": "MemoryStore"})

type item struct {
	value      interface{}
	createDate time.Time
	ttl        time.Duration
}

var _ interfaces.StoreInterface = &MemoryStore{}

// MemoryStore - In-app memory storage
type MemoryStore struct {
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
func (ms *MemoryStore) Get(_ context.Context, key string) (string, bool, time.Duration, error) {
	ms.lock.Lock()
	defer ms.lock.Unlock()

	mLogger.WithFields(log.Fields{"key": key}).Debug("Get key")
	if it, ok := ms.m[key]; ok {
		return it.value.(string), true, time.Since(it.createDate), nil
	}
	return "", false, 0, nil
}

// Set - Set value into store with ttl
func (ms *MemoryStore) Set(_ context.Context, key string, value string, ttl time.Duration) (bool, error) {
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
func (ms *MemoryStore) Update(_ context.Context, key string, value string) (bool, error) {
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
func (ms *MemoryStore) Remove(_ context.Context, key string) (bool, error) {
	ms.lock.Lock()
	delete(ms.m, key)
	ms.lock.Unlock()
	return true, nil
}

// Enqueue - ttl: optional params for setting the ttl of queue when first element is enqueued
func (ms *MemoryStore) Enqueue(_ context.Context, key string, value string, ttl time.Duration) (bool, error) {
	ms.lock.Lock()
	_, ok := ms.m[key]
	if !ok {
		ms.m[key] = &item{
			value:      []string{},
			createDate: time.Now(),
			ttl:        ttl,
		}
	}
	it := ms.m[key]

	it.value = append(it.value.([]string), value)
	ms.lock.Unlock()
	return true, nil
}

func (ms *MemoryStore) GetQueue(_ context.Context, key string) ([]string, bool, time.Duration, error) {
	ms.lock.Lock()
	defer ms.lock.Unlock()

	mLogger.WithFields(log.Fields{"key": key}).Debug("Get key")
	if it, ok := ms.m[key]; ok {
		return it.value.([]string), true, time.Since(it.createDate), nil
	}
	return []string{}, false, 0, nil
}

// Len - Get memory size
func (ms *MemoryStore) Len() int {
	ms.lock.Lock()
	defer ms.lock.Unlock()

	return len(ms.m)
}

// DataKey - Generate DataKey with events and destination
func (ms *MemoryStore) DataKey(_ context.Context, ev interfaces.IncomingEventInterface, dest interfaces.DestinationInterface, prefix string, suffix string) string {
	e := ev.(models.IncomingEvent)
	config := dest.(models.Destination).Config
	return fmt.Sprintf("%s%s.%s.%s%s", prefix, e.Key, config.GetName(), e.TargetId, suffix)
}
