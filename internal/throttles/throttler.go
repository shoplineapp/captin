package throttles

import (
	"context"
	"time"

	interfaces "github.com/shoplineapp/captin/v2/interfaces"
	log "github.com/sirupsen/logrus"
)

var tLogger = log.WithFields(log.Fields{"class": "Throttler"})

var _ interfaces.ThrottleInterface = &Throttler{}

// Throttler - Event Throttler
type Throttler struct {
	store interfaces.StoreInterface
}

// NewThrottler - Create new Throttler
func NewThrottler(store interfaces.StoreInterface) *Throttler {
	return &Throttler{
		store: store,
	}
}

func (t *Throttler) CanTrigger(ctx context.Context, id string, period time.Duration) (bool, time.Duration, error) {
	// ignore throttle if no period is given
	if period == time.Duration(0) {
		return true, time.Duration(0), nil
	}

	val, ok, duration, err := t.store.Get(ctx, id)

	if err != nil {
		return true, time.Duration(0), err
	}
	tLogger.WithFields(log.Fields{"value": val}).Debug("Check throttle value on CanTrigger")
	if !ok {
		tLogger.WithFields(log.Fields{"period": period}).Debug("Throttle value not set, creating...")
		t.store.Set(ctx, id, "1", period)
		return true, time.Duration(0), nil
	}

	return false, duration, nil
}
