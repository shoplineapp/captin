package mocks

import (
	"context"
	"time"

	"github.com/shoplineapp/captin/v2/interfaces"
	"github.com/stretchr/testify/mock"
)

var _ interfaces.ThrottleInterface = &ThrottleMock{}

// ThrottleMock - Mock ThrottleInterface
type ThrottleMock struct {
	mock.Mock
}

func (t *ThrottleMock) CanTrigger(ctx context.Context, id string, period time.Duration) (bool, time.Duration, error) {
	args := t.Called(ctx, id, period)
	return args.Bool(0), args.Get(1).(time.Duration), args.Error(2)
}
