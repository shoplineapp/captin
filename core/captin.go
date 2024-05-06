package core

import (
	"context"
	"fmt"

	destination_filters "github.com/shoplineapp/captin/v2/destinations/filters"
	d "github.com/shoplineapp/captin/v2/dispatcher"
	interfaces "github.com/shoplineapp/captin/v2/interfaces"
	outgoing "github.com/shoplineapp/captin/v2/internal/outgoing"
	models "github.com/shoplineapp/captin/v2/models"
	senders "github.com/shoplineapp/captin/v2/senders"

	captin_errors "github.com/shoplineapp/captin/v2/errors"
	documentStores "github.com/shoplineapp/captin/v2/internal/document_stores"
	stores "github.com/shoplineapp/captin/v2/internal/stores"
	throttles "github.com/shoplineapp/captin/v2/internal/throttles"
	log "github.com/sirupsen/logrus"
)

var cLogger = log.WithFields(log.Fields{"class": "Captin"})

var STATUS_READY = "ready"
var STATUS_RUNNING = "running"

// Captin - Captin instance
type Captin struct {
	Status               string
	ConfigMap            interfaces.ConfigMapperInterface
	filters              []destination_filters.DestinationFilterInterface
	middlewares          []destination_filters.DestinationMiddlewareInterface
	dispatchFilters      []destination_filters.DestinationFilterInterface
	dispatchMiddlewares  []destination_filters.DestinationMiddlewareInterface
	dispatchErrorHandler interfaces.ErrorHandlerInterface
	dispatchDelayer      interfaces.DispatchDelayerInterface
	SenderMapping        map[string]interfaces.EventSenderInterface
	store                interfaces.StoreInterface
	DocumentStoreMapping map[string]interfaces.DocumentStoreInterface
	throttler            interfaces.ThrottleInterface
}

// NewCaptin - Create Captin instance with default http senders and time throttler
func NewCaptin(configMap interfaces.ConfigMapperInterface) *Captin {
	store := stores.NewMemoryStore()
	senderMapping := map[string]interfaces.EventSenderInterface{
		"http":       &senders.HTTPEventSender{},
		"beanstalkd": &senders.BeanstalkdSender{},
	}
	c := Captin{
		Status:    STATUS_READY,
		ConfigMap: configMap,
		filters: []destination_filters.DestinationFilterInterface{
			destination_filters.ValidateFilter{},
			destination_filters.SourceFilter{},
			destination_filters.DesiredHookFilter{},
			destination_filters.EnvironmentFilter{},
		},
		SenderMapping: senderMapping,
		store:         store,
		DocumentStoreMapping: map[string]interfaces.DocumentStoreInterface{
			"default": documentStores.NewNullDocumentStore(),
		},
		throttler: throttles.NewThrottler(store),
	}
	return &c
}

// SetStore - Set store
func (c *Captin) SetStore(store interfaces.StoreInterface) {
	c.store = store
	c.throttler = throttles.NewThrottler(store)
}

// SetDocumentStoreMapping - Set store where event targets are being stored
func (c *Captin) SetDocumentStoreMapping(mappings map[string]interfaces.DocumentStoreInterface) {
	c.DocumentStoreMapping = mappings
}

// SetThrottler - Set throttle
func (c *Captin) SetThrottler(throttle interfaces.ThrottleInterface) {
	c.throttler = throttle
}

// SetDestinationFilters - Set filters
func (c *Captin) SetDestinationFilters(filters []destination_filters.DestinationFilterInterface) {
	c.filters = filters
}

// SetDestinationFilters - Set filters
func (c *Captin) SetDispatchFilters(filters []destination_filters.DestinationFilterInterface) {
	c.dispatchFilters = filters
}

// SetDestinationMiddlewares - Set middlewares
func (c *Captin) SetDestinationMiddlewares(middlewares []destination_filters.DestinationMiddlewareInterface) {
	c.middlewares = middlewares
}

func (c *Captin) SetDispatchMiddlewares(middlewares []destination_filters.DestinationMiddlewareInterface) {
	c.dispatchMiddlewares = middlewares
}

func (c *Captin) SetDispatchErrorHandler(handler interfaces.ErrorHandlerInterface) {
	c.dispatchErrorHandler = handler
}

func (c *Captin) SetDispatchDelayer(delayer interfaces.DispatchDelayerInterface) {
	c.dispatchDelayer = delayer
}

func (c *Captin) SetSenderMapping(senderMapping map[string]interfaces.EventSenderInterface) {
	c.SenderMapping = senderMapping
}

func (c Captin) IsRunning() bool {
	return c.Status == STATUS_RUNNING || d.PendingJobCount() > 0
}

// Execute - Execute for events
func (c *Captin) Execute(ctx context.Context, ie interfaces.IncomingEventInterface) (bool, []interfaces.ErrorInterface) {
	c.Status = STATUS_RUNNING

	e := ie.(models.IncomingEvent)
	if e.IsValid() != true {
		return false, []interfaces.ErrorInterface{&captin_errors.ExecutionError{Cause: "invalid incoming event object"}}
	}

	configs := c.ConfigMap.ConfigsForKey(e.Key)

	destinations := []models.Destination{}
	for _, config := range configs {
		destinations = append(destinations, models.Destination{Config: config})
	}

	destinations = outgoing.Custom{}.Sift(ctx, &e, destinations, c.filters, c.middlewares)
	cLogger.WithFields(log.Fields{
		"event":        e,
		"destinations": destinations,
	}).Info("Ready to dispatch event with destinations")

	// Create dispatcher and dispatch events
	dispatcher := outgoing.NewDispatcherWithDestinations(destinations, c.SenderMapping)
	dispatcher.SetFilters(c.dispatchFilters)
	dispatcher.SetMiddlewares(c.dispatchMiddlewares)
	dispatcher.SetErrorHandler(c.dispatchErrorHandler)
	dispatcher.SetDelayer(c.dispatchDelayer)
	dispatcher.Dispatch(ctx, e, c.store, c.throttler, c.DocumentStoreMapping)

	errors := dispatcher.GetErrors()

	cLogger.Debug(fmt.Sprintf("Captin event executed, %d destinations, %d failed, %d pending", len(destinations), len(errors), d.PendingJobCount()))

	c.Status = STATUS_READY
	return true, errors
}
