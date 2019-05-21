package interfaces

import "time"

// StoreInterface - Store for throttle events
type StoreInterface interface {
	// Get - Get value from store, return with remaining time
	Get(key string) (string, bool, time.Duration, error)

	// Set - Set value into store with ttl
	Set(key string, value string, ttl time.Duration) (bool, error)

	// Update - Update value for key
	Update(key string, value string) (bool, error)

	// Remove - Remove value for key
	Remove(key string) (bool, error)
}
