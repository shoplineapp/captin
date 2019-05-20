package core

import (
	"fmt"

	interfaces "captin/interfaces"
	models "captin/internal/models"
	outgoing "captin/internal/outgoing"
	outgoing_filters "captin/internal/outgoing/filters"
	senders "captin/internal/senders"
)

type ExecutionError struct {
	Cause string
}

func (e *ExecutionError) Error() string {
	return fmt.Sprintf("ExecutionError: caused by %s", e.Cause)
}

type Captin struct {
	ConfigMap models.ConfigurationMapper
	filters   []interfaces.CustomFilter
}

func NewCaptin(configMap models.ConfigurationMapper) *Captin {
	c := Captin{
		ConfigMap: configMap,
		filters: []interfaces.CustomFilter{
			outgoing_filters.ValidateFilter{},
			outgoing_filters.SourceFilter{},
		},
	}
	return &c
}

func (c Captin) Execute(e models.IncomingEvent) (bool, error) {
	if e.IsValid() != true {
		return false, &ExecutionError{Cause: "invalid incoming event object"}
	}

	configs := c.ConfigMap.ConfigsForKey(e.Key)

	destinations := []models.Destination{}

	for _, config := range configs {
		destinations = append(destinations, models.Destination{Config: config})
	}

	destinations = outgoing.Custom{}.Sift(e, c.filters, destinations)

	// TODO: Pass event and destinations into dispatcher

	// Create dispatcher and dispatch events
	sender := senders.HTTPEventSender{}
	dispatcher := outgoing.NewDispatcherWithDestinations(destinations, &sender)
	dispatcher.Dispatch(e)

	for _, err := range dispatcher.Errors {
		switch dispatcherErr := err.(type) {
		case *outgoing.DispatcherError:
			fmt.Println("[Dispatcher] Error on event: ", dispatcherErr.Event.TargetId)
			fmt.Println("[Dispatcher] Error on event type: ", dispatcherErr.Event.TargetType)
		default:
			fmt.Println(e)
		}
	}

	return true, nil
}
