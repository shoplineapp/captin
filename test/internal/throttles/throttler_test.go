package throttles_test

import (
	"context"
	"errors"
	"testing"
	"time"

	throttles "github.com/shoplineapp/captin/v2/internal/throttles"

	mocks "github.com/shoplineapp/captin/v2/test/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setup() (string, time.Duration, *mocks.StoreMock) {
	throttleID := "123"
	throttlePeriod := time.Millisecond * 10

	store := new(mocks.StoreMock)
	return throttleID, throttlePeriod, store
}
func TestThrottler_CreateStoreRecord(t *testing.T) {
	throttleID, throttlePeriod, store := setup()

	store.On("Get", mock.Anything, throttleID).Return("", false, time.Millisecond*10, nil)
	store.On("Set", mock.Anything, throttleID, "1", throttlePeriod).Return(true, nil)

	subject := throttles.NewThrottler(store)
	result, duration, err := subject.CanTrigger(context.Background(), throttleID, throttlePeriod)
	assert.True(t, result)
	assert.Equal(t, time.Duration(0), duration)
	assert.Nil(t, err)

	store.AssertCalled(t, "Get", mock.Anything, throttleID)
	store.AssertCalled(t, "Set", mock.Anything, throttleID, "1", throttlePeriod)
}

func TestThrottler_NoThrottleSet(t *testing.T) {
	throttleID, _, store := setup()
	throttlePeriod := time.Duration(0)

	subject := throttles.NewThrottler(store)
	result, duration, err := subject.CanTrigger(context.Background(), throttleID, throttlePeriod)
	assert.True(t, result)
	assert.Equal(t, time.Duration(0), duration)
	assert.Nil(t, err)

	store.AssertNotCalled(t, "Get", mock.Anything, throttleID)
	store.AssertNotCalled(t, "Set", mock.Anything, throttleID, "1", throttlePeriod)
}

func TestThrottler_Reject(t *testing.T) {
	throttleID, throttlePeriod, store := setup()

	store.On("Get", mock.Anything, throttleID).Return("1", true, time.Millisecond*10, nil)

	subject := throttles.NewThrottler(store)
	result, duration, err := subject.CanTrigger(context.Background(), throttleID, throttlePeriod)
	assert.False(t, result)
	assert.Equal(t, time.Millisecond*10, duration)
	assert.Nil(t, err)

	store.AssertCalled(t, "Get", mock.Anything, throttleID)
}

func TestThrottler_Error(t *testing.T) {
	throttleID, throttlePeriod, store := setup()

	store.On("Get", mock.Anything, throttleID).Return("", false, time.Duration(0), errors.New("some error"))

	subject := throttles.NewThrottler(store)
	result, duration, err := subject.CanTrigger(context.Background(), throttleID, throttlePeriod)
	assert.True(t, result)
	assert.Equal(t, time.Duration(0), duration)
	assert.EqualError(t, err, "some error")

	store.AssertCalled(t, "Get", mock.Anything, throttleID)
}
