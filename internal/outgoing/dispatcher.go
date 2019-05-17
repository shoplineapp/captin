package outgoing

import (
	models "github.com/shoplineapp/captin/internal/models"
	sender "github.com/shoplineapp/captin/internal/senders"
)

// Dispatcher - Event Dispatcher
type Dispatcher struct {
	callbacks []chan models.IncomingEvent
	sender    sender.EventSenderInterface
}

// NewDispatcherWithConfig - Create Outgoing event dispatcher with config
func NewDispatcherWithConfig(configs []models.Configuration, sender sender.EventSenderInterface) *Dispatcher {
	result := Dispatcher{
		callbacks: []chan models.IncomingEvent{},
		sender:    sender,
	}

	for _, config := range configs {
		ch := make(chan models.IncomingEvent)
		go func(conf models.Configuration) {
			evt := <-ch
			sender.SendEvent(evt, config)
		}(config)
		result.callbacks = append(result.callbacks, ch)
	}

	return &result
}

// Dispatch - Dispatch an event to outgoing webhook
func (d *Dispatcher) Dispatch(e models.IncomingEvent) error {

	for _, handler := range d.callbacks {
		go func(handler chan models.IncomingEvent) {
			handler <- e
		}(handler)
	}

	return nil
}
