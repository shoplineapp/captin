package core

import (
	models "captin/internal/models"
)

// CustomMiddleware - Interface for third-party application to add extra handling on destinations
type CustomMiddleware interface {
	Apply(e models.IncomingEvent, d []models.Destination) (models.IncomingEvent, []models.Destination)
}

// CustomFilter - Interface for third-party application to filter destination by event
type CustomFilter interface {
	Run(e models.IncomingEvent, c models.Destination) (bool, error)
	Applicable(e models.IncomingEvent, c models.Destination) bool
}

// CaptinInterface
type CaptinInterface interface {
	Execute(e models.IncomingEvent) (bool, error)
}

// EventSenderInterface - Event Sender Interface
type EventSenderInterface interface {
	SendEvent(e models.IncomingEvent, d models.Destination) error
}
