package core

import (
	"fmt"

	interfaces "github.com/shoplineapp/captin/interfaces"
	outgoing "github.com/shoplineapp/captin/internal/outgoing"
	outgoing_filters "github.com/shoplineapp/captin/internal/outgoing/filters"
	senders "github.com/shoplineapp/captin/internal/senders"
	models "github.com/shoplineapp/captin/models"

	stores "github.com/shoplineapp/captin/internal/stores"
	throttles "github.com/shoplineapp/captin/internal/throttles"
	log "github.com/sirupsen/logrus"
)

var cLogger = log.WithFields(log.Fields{"class": "Captin"})

// ExecutionError - Error on executing events
type ExecutionError struct {
	Cause string
}

func (e *ExecutionError) Error() string {
	return fmt.Sprintf("ExecutionError: caused by %s", e.Cause)
}

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
		"http": &senders.HTTPEventSender{},
		"beanstalkd": &senders.BeanstalkdSender{},
	}
	c := Captin{
		ConfigMap: configMap,
		filters: []interfaces.DestinationFilter{
			outgoing_filters.ValidateFilter{},
			outgoing_filters.SourceFilter{},
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
func (c Captin) Execute(e models.IncomingEvent) (bool, error) {
	if e.IsValid() != true {
		return false, &ExecutionError{Cause: "invalid incoming event object"}
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
	}).Debug("Ready to dispatch event with destinations")

	// Create dispatcher and dispatch events
	dispatcher := outgoing.NewDispatcherWithDestinations(destinations, c.SenderMapping)
	dispatcher.Dispatch(e, c.store, c.throttler)

	for _, err := range dispatcher.Errors {
		switch dispatcherErr := err.(type) {
		case *outgoing.DispatcherError:
			cLogger.WithFields(log.Fields{
				"target_id":   dispatcherErr.Event.TargetId,
				"target_type": dispatcherErr.Event.TargetType,
			}).Error("Failed to dispatch event")
		default:
			cLogger.WithFields(log.Fields{"error": e}).Error("Unhandled error on dispatcher")
		}
	}

	return true, nil
}
