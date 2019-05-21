package throttles_test

import (
	"testing"
	"time"

	interfaces "github.com/shoplineapp/captin/interfaces"

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
	return args.String(0), args.Bool(1), args.Get(3).(time.Duration), args.Error(4)
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

func TestThrottler(t *testing.T) {
	assert.True(t, true, true)
}
