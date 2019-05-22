package mocks

import (
	"time"

	"github.com/shoplineapp/captin/interfaces"
	"github.com/stretchr/testify/mock"
)

// ThrottleMock - Mock ThrottleInterface
type ThrottleMock struct {
	interfaces.ThrottleInterface
	mock.Mock
}

// CanTrigger - Check if can trigger
func (t *ThrottleMock) CanTrigger(id string, period time.Duration) (bool, time.Duration, error) {
	args := t.Called(id, period)
	return args.Bool(0), args.Get(1).(time.Duration), args.Error(2)
}
