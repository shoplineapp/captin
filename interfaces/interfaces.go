package interfaces

import (
	"time"
)

// CaptinInterface - Captin Interface
type CaptinInterface interface {
	Execute(e IncomingEventInterface) (bool, []ErrorInterface)
	IsRunning() bool
}

// EventSenderInterface - Event Sender Interface
type EventSenderInterface interface {
	SendEvent(e IncomingEventInterface, d DestinationInterface) error
}

// ThrottleInterface - interface for a throttle object
// Throttle event flow:
// Mutex Lock:		1
// Payload Store:	p p p p (Expired * 2)
// Events:				x x x x -
// Throttle: 			  t     t2
type ThrottleInterface interface {
	// CanTrigger - Check if can trigger
	CanTrigger(id string, period time.Duration) (bool, time.Duration, error)
}

type ErrorHandlerInterface interface {
	Exec(e ErrorInterface)
}
