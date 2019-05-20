package senders

import (
	models "captin/internal/models"
)

// EventSenderInterface - Event Sender Interface
type EventSenderInterface interface {
	SendEvent(e models.IncomingEvent, d models.Destination) error
}
