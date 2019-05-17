package outgoing

import (
	"fmt"

	models "github.com/shoplineapp/captin/internal/models"
)

// Dispatcher - Event Dispatcher
type Dispatcher struct {
	callbacks []chan models.IncomingEvent
}

// NewDispatcherWithConfig - Create Outgoing event dispatcher with config
func NewDispatcherWithConfig(configs []models.Configuration) *Dispatcher {
	result := Dispatcher{
		callbacks: []chan models.IncomingEvent{},
	}

	for _, config := range configs {
		ch := make(chan models.IncomingEvent)
		go func(conf models.Configuration) {
			evt := <-ch
			fmt.Println("Configuration: \t\t", conf.Name)
			fmt.Println("Process Event ID: \t", evt.TargetId)
			fmt.Println("Process Event Type: \t", evt.TargetType)
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
