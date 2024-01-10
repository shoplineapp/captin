package interfaces

import (
	"context"
	"time"
)

// StoreInterface - Store for throttle events
type StoreInterface interface {
	// Get - Get value from store, return with remaining time
	Get(ctx context.Context, key string) (payload string, exists bool, ttl time.Duration, err error)

	// Set - Set value into store with ttl
	Set(ctx context.Context, key string, value string, ttl time.Duration) (bool, error)

	// Update - Update value for key
	Update(ctx context.Context, key string, value string) (bool, error)

	// Remove - Remove value for key
	Remove(ctx context.Context, key string) (bool, error)

	DataKey(ctx context.Context, e IncomingEventInterface, dest DestinationInterface, prefix string, suffix string) string

	Enqueue(ctx context.Context, key string, value string, ttl time.Duration) (bool, error)

	GetQueue(ctx context.Context, key string) (values []string, exists bool, ttl time.Duration, err error)
}
