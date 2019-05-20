package throttle_test

import (
	"sync"
	"testing"
	"time"

	throttle "captin/internal/outgoing/throttle"
)

func TestTimeThrottle(t *testing.T) {
	var wg sync.WaitGroup

	throttle := throttle.NewTimeThrottle(time.Millisecond, false)
	count := 0

	wg.Add(1)
	go func() {
		defer wg.Done()
		for throttle.Next() {
			count += 1
		}
	}()

	for i := 0; i < 5; i++ {
		throttle.Trigger()
	}

	time.Sleep(5 * time.Millisecond)

	throttle.Stop()

	wg.Wait()

	if count != 1 {
		t.Errorf("count = %v", count)
	}
}
