package throttles

import (
	"time"

	interfaces "github.com/shoplineapp/captin/interfaces"
)

// MemstoreThrottle - Throttle by memory store
type MemstoreThrottle struct {
	interfaces.ThrottleInterface
}

// NewMemstoreThrottle - Create Memstore Throttle
func NewMemstoreThrottle() (*MemstoreThrottle, error) {
	return &MemstoreThrottle{}, nil
}

// CanTrigger - Check can be trigger
func (m *MemstoreThrottle) CanTrigger(id string) (bool, time.Duration, error) {
	return true, -1, nil
}
