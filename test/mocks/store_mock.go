package mocks

import (
	"context"
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
func (s *StoreMock) Get(ctx context.Context, key string) (string, bool, time.Duration, error) {
	args := s.Called(ctx, key)
	return args.String(0), args.Bool(1), args.Get(2).(time.Duration), args.Error(3)
}

// Set - Set value into store with ttl
func (s *StoreMock) Set(ctx context.Context, key string, value string, ttl time.Duration) (bool, error) {
	args := s.Called(ctx, key, value, ttl)
	return args.Bool(0), args.Error(1)
}

// Update - Update value for key
func (s *StoreMock) Update(ctx context.Context, key string, value string) (bool, error) {
	args := s.Called(ctx, key, value)
	return args.Bool(0), args.Error(1)
}

func (s *StoreMock) Enqueue(ctx context.Context, key string, value string, ttl time.Duration) (bool, error) {
	args := s.Called(ctx, key, value, ttl)
	return args.Bool(0), args.Error(1)
}

func (s *StoreMock) GetQueue(ctx context.Context, key string) ([]string, bool, time.Duration, error) {
	args := s.Called(ctx, key)
	return args.Get(0).([]string), args.Bool(1), args.Get(2).(time.Duration), args.Error(3)
}

// Remove - Remove value for key
func (s *StoreMock) Remove(ctx context.Context, key string) (bool, error) {
	args := s.Called(ctx, key)
	return args.Bool(0), args.Error(1)
}

// DataKey - Generate DataKey with events and destination (Won't Mock)
func (s *StoreMock) DataKey(ctx context.Context, ie interfaces.IncomingEventInterface, idest interfaces.DestinationInterface, prefix string, suffix string) string {
	e := ie.(models.IncomingEvent)
	dest := idest.(models.Destination)
	return fmt.Sprintf("%s%s.%s.%s%s", prefix, e.Key, dest.Config.GetName(), e.TargetId, suffix)
}
