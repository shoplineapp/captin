package sl_time

import (
	"sync/atomic"
	"time"
)

var pendingJobCount int64 = 0

func PendingJobCount() int64 {
	return pendingJobCount
}

func AfterFunc(d time.Duration, f func()) *time.Timer {
	return time.AfterFunc(d, func() {
		atomic.AddInt64(&pendingJobCount, 1)
		defer atomic.AddInt64(&pendingJobCount, -1)
		f()
	})
}