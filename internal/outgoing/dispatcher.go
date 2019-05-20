package outgoing

import (
	models "github.com/shoplineapp/captin/internal/models"
	sender "github.com/shoplineapp/captin/internal/senders"
)

// DispatcherError - Error when send events
type DispatcherError struct {
	msg         string
	Event       models.IncomingEvent
	Destination models.Destination
}

func (e *DispatcherError) Error() string {
	return e.msg
}

// Dispatcher - Event Dispatcher
type Dispatcher struct {
	destinations []models.Destination
	sender       sender.EventSenderInterface
	Errors       []error
}

// NewDispatcherWithDestinations - Create Outgoing event dispatcher with destinations
func NewDispatcherWithDestinations(destinations []models.Destination, sender sender.EventSenderInterface) *Dispatcher {
	result := Dispatcher{
		destinations: destinations,
		sender:       sender,
		Errors:       []error{},
	}

	return &result
}

// Dispatch - Dispatch an event to outgoing webhook
func (d *Dispatcher) Dispatch(e models.IncomingEvent) error {
	for _, destination := range d.destinations {
		go func(evt models.IncomingEvent, destination models.Destination) {
			err := d.sender.SendEvent(evt, destination)
			if err != nil {
				d.Errors = append(d.Errors, &DispatcherError{
					msg:         err.Error(),
					Destination: destination,
					Event:       evt,
				})
			}
		}(e, destination)
	}
	return nil
}
