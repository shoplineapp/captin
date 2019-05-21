package interfaces

import "time"

// StoreInterface - Store for throttle events
type StoreInterface interface {
	// Get - Get value from store, return with remaining time
	Get(key string) (string, time.Time, error)

	// Set - Set value into store with ttl
	Set(key string, value string, ttl time.Duration) (bool, error)
}
