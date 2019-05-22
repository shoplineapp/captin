package throttles

import (
	"fmt"
	"time"

	interfaces "github.com/shoplineapp/captin/interfaces"
)

// Throttler - Event Throttler
type Throttler struct {
	interfaces.ThrottleInterface
	store interfaces.StoreInterface
}

// NewThrottler - Create new Throttler
func NewThrottler(store interfaces.StoreInterface) *Throttler {
	return &Throttler{
		store: store,
	}
}

// CanTrigger - Check if can trigger
func (t *Throttler) CanTrigger(id string, period time.Duration) (bool, time.Duration, error) {
	val, ok, duration, err := t.store.Get(id)

	if err != nil {
		return true, time.Duration(0), err
	}
	fmt.Println("[Throttler] Value: ", val)
	if !ok {
		fmt.Println("[Throttler] Create throttle in store with period, ", period)
		t.store.Set(id, "1", period)
		return true, time.Duration(0), nil
	}

	return false, duration, nil
}
