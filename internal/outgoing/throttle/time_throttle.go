package throttles

import (
	"sync"
	"time"
)

// TimeThrottle - Throttle by time
type TimeThrottle struct {
	cond     *sync.Cond
	period   time.Duration
	last     time.Time
	trailing bool
	waiting  bool
	stop     bool
}

// NewTimeThrottle - Create time throttle
func NewTimeThrottle(period time.Duration, trailing bool) *TimeThrottle {
	time.Tick(200)
	return &TimeThrottle{
		period:   period,
		trailing: trailing,
		cond:     sync.NewCond(&sync.Mutex{}),
	}
}

// Trigger - Trigger throttle
func (t *TimeThrottle) Trigger() {
	t.cond.L.Lock()
	defer t.cond.L.Unlock()

	if !t.waiting && !t.stop {

		delta := time.Now().Sub(t.last)

		if delta > t.period {
			t.waiting = true
			t.cond.Broadcast()
		} else if t.trailing {
			t.waiting = true
			time.AfterFunc(t.period-delta, t.cond.Broadcast)
		}
	}
}

// Next - Returns true at most once per period
func (t *TimeThrottle) Next() bool {
	t.cond.L.Lock()
	defer t.cond.L.Unlock()
	for !t.waiting && !t.stop {
		t.cond.Wait()
	}
	if !t.stop {
		t.waiting = false
		t.last = time.Now()
	}
	return !t.stop
}

// Stop - Stop this throttler
func (t *TimeThrottle) Stop() {
	t.cond.L.Lock()
	defer t.cond.L.Unlock()
	t.stop = true
	t.cond.Broadcast()
}
