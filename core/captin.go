package core

import (
	"fmt"

	interfaces "github.com/shoplineapp/captin/interfaces"
	outgoing "github.com/shoplineapp/captin/internal/outgoing"
	outgoing_filters "github.com/shoplineapp/captin/internal/outgoing/filters"
	senders "github.com/shoplineapp/captin/internal/senders"
	models "github.com/shoplineapp/captin/models"
)

type ExecutionError struct {
	Cause string
}

func (e *ExecutionError) Error() string {
	return fmt.Sprintf("ExecutionError: caused by %s", e.Cause)
}

type Captin struct {
	ConfigMap   interfaces.ConfigMapperInterface
	filters     []interfaces.DestinationFilter
	middlewares []interfaces.DestinationMiddleware
}

func NewCaptin(configMap interfaces.ConfigMapperInterface) *Captin {
	c := Captin{
		ConfigMap: configMap,
		filters: []interfaces.DestinationFilter{
			outgoing_filters.ValidateFilter{},
			outgoing_filters.SourceFilter{},
		},
		middlewares: []interfaces.DestinationMiddleware{},
	}
	return &c
}

func (c *Captin) SetDestinationFilters(filters []interfaces.DestinationFilter) {
	c.filters = filters
}

func (c *Captin) SetDestinationMiddlewares(middlewares []interfaces.DestinationMiddleware) {
	c.middlewares = middlewares
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

	destinations = outgoing.Custom{}.Sift(e, destinations, c.filters, c.middlewares)

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
