package dispatcher

import (
	"sync/atomic"
	"time"
)

var pendingJobCount int64 = 0

func PendingJobCount() int64 {
	return atomic.LoadInt64(&pendingJobCount)
}

func TrackAfterFuncJob(d time.Duration, f func()) {
	atomic.AddInt64(&pendingJobCount, 1)

	go func() {
		time.AfterFunc(d, func() {
			f()
			defer atomic.AddInt64(&pendingJobCount, -1)
		})
	}()
}

func TrackGoRoutine(f func()) {
	atomic.AddInt64(&pendingJobCount, 1)

	go func() {
		f()
		defer atomic.AddInt64(&pendingJobCount, -1)
	}()
}
