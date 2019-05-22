package core

import (
	"fmt"

	interfaces "github.com/shoplineapp/captin/interfaces"
	outgoing "github.com/shoplineapp/captin/internal/outgoing"
	outgoing_filters "github.com/shoplineapp/captin/internal/outgoing/filters"
	senders "github.com/shoplineapp/captin/internal/senders"
	models "github.com/shoplineapp/captin/models"

	throttles "github.com/shoplineapp/captin/internal/throttles"
)

// ExecutionError - Error on executing events
type ExecutionError struct {
	Cause string
}

func (e *ExecutionError) Error() string {
	return fmt.Sprintf("ExecutionError: caused by %s", e.Cause)
}

// Captin - Captin instance
type Captin struct {
	ConfigMap   interfaces.ConfigMapperInterface
	filters     []interfaces.DestinationFilter
	middlewares []interfaces.DestinationMiddleware
	sender      interfaces.EventSenderInterface
	store       interfaces.StoreInterface
	throttler   interfaces.ThrottleInterface
}

// NewCaptin - Create Captin instance with default http senders and time throttler
func NewCaptin(
	configMap interfaces.ConfigMapperInterface,
	store interfaces.StoreInterface) *Captin {
	c := Captin{
		ConfigMap: configMap,
		filters: []interfaces.DestinationFilter{
			outgoing_filters.ValidateFilter{},
			outgoing_filters.SourceFilter{},
		},
		sender:    &senders.HTTPEventSender{},
		store:     store,
		throttler: throttles.NewThrottler(store),
	}
	return &c
}

// SetThrottler - Set throttle
func (c *Captin) SetThrottler(throttle interfaces.ThrottleInterface) {
	c.throttler = throttle
}

// SetDestinationFilters - Set filters
func (c *Captin) SetDestinationFilters(filters []interfaces.DestinationFilter) {
	c.filters = filters
}

// SetDestinationMiddlewares - Set middlewares
func (c *Captin) SetDestinationMiddlewares(middlewares []interfaces.DestinationMiddleware) {
	c.middlewares = middlewares
}

// Execute - Execute for events
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
	dispatcher := outgoing.NewDispatcherWithDestinations(destinations, c.sender)
	dispatcher.Dispatch(e, c.store, c.throttler)

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
