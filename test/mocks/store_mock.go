package mocks

import (
	"fmt"
	"time"

	"github.com/shoplineapp/captin/interfaces"
	"github.com/shoplineapp/captin/models"
	"github.com/stretchr/testify/mock"
)

// StoreMock - Mock of StoreInterface
type StoreMock struct {
	interfaces.StoreInterface
	mock.Mock
}

// Get - Get value from store, return with remaining time
func (s *StoreMock) Get(key string) (string, bool, time.Duration, error) {
	args := s.Called(key)
	return args.String(0), args.Bool(1), args.Get(2).(time.Duration), args.Error(3)
}

// Set - Set value into store with ttl
func (s *StoreMock) Set(key string, value string, ttl time.Duration) (bool, error) {
	args := s.Called(key, value, ttl)
	return args.Bool(0), args.Error(1)
}

// Update - Update value for key
func (s *StoreMock) Update(key string, value string) (bool, error) {
	args := s.Called(key, value)
	return args.Bool(0), args.Error(1)
}

// Remove - Remove value for key
func (s *StoreMock) Remove(key string) (bool, error) {
	args := s.Called(key)
	return args.Bool(0), args.Error(1)
}

// DataKey - Generate DataKey with events and destination (Won't Mock)
func (s *StoreMock) DataKey(e models.IncomingEvent, dest models.Destination, prefix string, suffix string) string {
	return fmt.Sprintf("%s%s.%s.%s%s", prefix, e.Key, dest.Config.Name, e.TargetId, suffix)
}
