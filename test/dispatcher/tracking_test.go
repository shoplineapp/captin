package dispatcher_delayers_test

import (
	"testing"
	"time"

	"github.com/shoplineapp/captin/v2/dispatcher"
	"github.com/stretchr/testify/assert"
)

func TestTrackAfterFuncJob(t *testing.T) {
	c := make(chan int)

	mock := func() {
		c <- 1
	}

	dispatcher.TrackAfterFuncJob(1*time.Millisecond, mock)

	assert.Equal(t, 1, <-c)
	assert.EqualValues(t, 0, dispatcher.PendingJobCount())

}

func TestTrackGoRoutine(t *testing.T) {
	c := make(chan int)

	mock := func() {
		c <- 1
	}

	dispatcher.TrackGoRoutine(mock)

	assert.Equal(t, 1, <-c)
	assert.EqualValues(t, 0, dispatcher.PendingJobCount())

}
