package core

import (
	destination_filters "github.com/shoplineapp/captin/destinations/filters"
	interfaces "github.com/shoplineapp/captin/interfaces"
	outgoing "github.com/shoplineapp/captin/internal/outgoing"
	models "github.com/shoplineapp/captin/models"
	senders "github.com/shoplineapp/captin/senders"

	captin_errors "github.com/shoplineapp/captin/errors"
	stores "github.com/shoplineapp/captin/internal/stores"
	throttles "github.com/shoplineapp/captin/internal/throttles"
	log "github.com/sirupsen/logrus"
)

var cLogger = log.WithFields(log.Fields{"class": "Captin"})

// Captin - Captin instance
type Captin struct {
	ConfigMap     interfaces.ConfigMapperInterface
	filters       []interfaces.DestinationFilter
	middlewares   []interfaces.DestinationMiddleware
	SenderMapping map[string]interfaces.EventSenderInterface
	store         interfaces.StoreInterface
	throttler     interfaces.ThrottleInterface
}

// NewCaptin - Create Captin instance with default http senders and time throttler
func NewCaptin(configMap interfaces.ConfigMapperInterface) *Captin {
	store := stores.NewMemoryStore()
	senderMapping := map[string]interfaces.EventSenderInterface{
		"http":       &senders.HTTPEventSender{},
		"beanstalkd": &senders.BeanstalkdSender{},
	}
	c := Captin{
		ConfigMap: configMap,
		filters: []interfaces.DestinationFilter{
			destination_filters.ValidateFilter{},
			destination_filters.SourceFilter{},
			destination_filters.DesiredHookFilter{},
		},
		SenderMapping: senderMapping,
		store:         store,
		throttler:     throttles.NewThrottler(store),
	}
	return &c
}

// SetStore - Set store
func (c *Captin) SetStore(store interfaces.StoreInterface) {
	c.store = store
	c.throttler = throttles.NewThrottler(store)
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

func (c *Captin) SetSenderMapping(senderMapping map[string]interfaces.EventSenderInterface) {
	c.SenderMapping = senderMapping
}

// Execute - Execute for events
func (c Captin) Execute(e models.IncomingEvent) (bool, []captin_errors.ErrorInterface) {
	if e.IsValid() != true {
		return false, []captin_errors.ErrorInterface{&captin_errors.ExecutionError{Cause: "invalid incoming event object"}}
	}

	configs := c.ConfigMap.ConfigsForKey(e.Key)

	destinations := []models.Destination{}

	for _, config := range configs {
		destinations = append(destinations, models.Destination{Config: config})
	}

	destinations = outgoing.Custom{}.Sift(&e, destinations, c.filters, c.middlewares)
	cLogger.WithFields(log.Fields{
		"event":        e,
		"destinations": destinations,
	}).Info("Ready to dispatch event with destinations")

	// Create dispatcher and dispatch events
	dispatcher := outgoing.NewDispatcherWithDestinations(destinations, c.SenderMapping)
	dispatcher.Dispatch(e, c.store, c.throttler)

	for _, err := range dispatcher.Errors {
		switch dispatcherErr := err.(type) {
		case *captin_errors.DispatcherError:
			cLogger.WithFields(log.Fields{
				"event":       dispatcherErr.Event,
				"destination": dispatcherErr.Destination,
				"reason":      dispatcherErr.Error(),
			}).Error("Failed to dispatch event")
		default:
			cLogger.WithFields(log.Fields{"error": e}).Error("Unhandled error on dispatcher")
		}
	}

	return true, dispatcher.Errors
}
