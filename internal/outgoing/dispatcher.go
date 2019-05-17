package outgoing

import (
	models "github.com/shoplineapp/captin/internal/models"
	sender "github.com/shoplineapp/captin/internal/senders"
)

// DispatcherError - Error when send events
type DispatcherError struct {
	msg           string
	Event         models.IncomingEvent
	Configuration models.Configuration
}

func (e *DispatcherError) Error() string {
	return e.msg
}

// Dispatcher - Event Dispatcher
type Dispatcher struct {
	configs []models.Configuration
	sender  sender.EventSenderInterface
	Errors  []error
}

// NewDispatcherWithConfig - Create Outgoing event dispatcher with config
func NewDispatcherWithConfig(configs []models.Configuration, sender sender.EventSenderInterface) *Dispatcher {
	result := Dispatcher{
		configs: configs,
		sender:  sender,
		Errors:  []error{},
	}

	return &result
}

// Dispatch - Dispatch an event to outgoing webhook
func (d *Dispatcher) Dispatch(e models.IncomingEvent) error {
	for _, config := range d.configs {
		go func(evt models.IncomingEvent, conf models.Configuration) {
			err := d.sender.SendEvent(evt, conf)
			if err != nil {
				d.Errors = append(d.Errors, &DispatcherError{
					msg:           err.Error(),
					Configuration: conf,
					Event:         evt,
				})
			}
		}(e, config)
	}
	return nil
}
