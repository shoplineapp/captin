package throttles_test

import (
	"errors"
	"testing"
	"time"

	interfaces "github.com/shoplineapp/captin/interfaces"
	throttles "github.com/shoplineapp/captin/internal/throttles"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type storeMock struct {
	interfaces.StoreInterface
	mock.Mock
}

// Get - Get value from store, return with remaining time
func (s *storeMock) Get(key string) (string, bool, time.Duration, error) {
	args := s.Called(key)
	return args.String(0), args.Bool(1), args.Get(2).(time.Duration), args.Error(3)
}

// Set - Set value into store with ttl
func (s *storeMock) Set(key string, value string, ttl time.Duration) (bool, error) {
	args := s.Called(key, value, ttl)
	return args.Bool(0), args.Error(1)
}

// Update - Update value for key
func (s *storeMock) Update(key string, value string) (bool, error) {
	args := s.Called(key, value)
	return args.Bool(0), args.Error(1)
}

// Remove - Remove value for key
func (s *storeMock) Remove(key string) (bool, error) {
	args := s.Called(key)
	return args.Bool(0), args.Error(1)
}

func setup() (string, time.Duration, *storeMock) {
	throttleID := "123"
	throttlePeriod := time.Millisecond * 10

	store := new(storeMock)
	return throttleID, throttlePeriod, store
}
func TestThrottler_CreateStoreRecord(t *testing.T) {
	throttleID, throttlePeriod, store := setup()

	store.On("Get", throttleID).Return("", false, time.Millisecond*10, nil)
	store.On("Set", throttleID, "1", throttlePeriod).Return(true, nil)

	subject := throttles.NewThrottler(throttlePeriod, store)
	result, duration, err := subject.CanTrigger(throttleID)
	assert.True(t, result)
	assert.Equal(t, time.Duration(0), duration)
	assert.Nil(t, err)

	store.AssertCalled(t, "Get", throttleID)
	store.AssertCalled(t, "Set", throttleID, "1", throttlePeriod)
}

func TestThrottler_Reject(t *testing.T) {
	throttleID, throttlePeriod, store := setup()

	store.On("Get", throttleID).Return("1", true, time.Millisecond*10, nil)

	subject := throttles.NewThrottler(throttlePeriod, store)
	result, duration, err := subject.CanTrigger(throttleID)
	assert.False(t, result)
	assert.Equal(t, time.Millisecond*10, duration)
	assert.Nil(t, err)

	store.AssertCalled(t, "Get", throttleID)
}

func TestThrottler_Error(t *testing.T) {
	throttleID, throttlePeriod, store := setup()

	store.On("Get", throttleID).Return("", false, time.Duration(0), errors.New("some error"))

	subject := throttles.NewThrottler(throttlePeriod, store)
	result, duration, err := subject.CanTrigger(throttleID)
	assert.True(t, result)
	assert.Equal(t, time.Duration(0), duration)
	assert.EqualError(t, err, "some error")

	store.AssertCalled(t, "Get", throttleID)
}
