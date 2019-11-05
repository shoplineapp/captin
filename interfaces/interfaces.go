package interfaces

import (
	"time"

	models "github.com/shoplineapp/captin/models"
)

// CaptinInterface - Captin Interface
type CaptinInterface interface {
	Execute(e models.IncomingEvent) (bool, []error)
}

// EventSenderInterface - Event Sender Interface
type EventSenderInterface interface {
	SendEvent(e models.IncomingEvent, d models.Destination) error
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
