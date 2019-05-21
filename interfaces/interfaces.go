package interfaces

import (
	models "github.com/shoplineapp/captin/models"
)

// CaptinInterface
type CaptinInterface interface {
	Execute(e models.IncomingEvent) (bool, error)
}

// EventSenderInterface - Event Sender Interface
type EventSenderInterface interface {
	SendEvent(e models.IncomingEvent, d models.Destination) error
}
