package interfaces

import (
	"context"
	"time"
)

// CaptinInterface - Captin Interface
type CaptinInterface interface {
	Execute(ctx context.Context, e IncomingEventInterface) (bool, []ErrorInterface)
	IsRunning() bool
}

// EventSenderInterface - Event Sender Interface
type EventSenderInterface interface {
	SendEvent(ctx context.Context, e IncomingEventInterface, d DestinationInterface) error
}

// ThrottleInterface - interface for a throttle object
// Throttle event flow:
// Mutex Lock:		1
// Payload Store:	p p p p (Expired * 2)
// Events:				x x x x -
// Throttle: 			  t     t2
type ThrottleInterface interface {
	// CanTrigger - Check if can trigger
	CanTrigger(ctx context.Context, id string, period time.Duration) (canTrigger bool, ttl time.Duration, err error)
}

type ErrorHandlerInterface interface {
	Exec(ctx context.Context, e ErrorInterface)
}
