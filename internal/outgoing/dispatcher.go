package outgoing

import (
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/mohae/deepcopy"
	destination_filters "github.com/shoplineapp/captin/destinations/filters"
	"github.com/shoplineapp/captin/dispatcher"
	captin_errors "github.com/shoplineapp/captin/errors"
	interfaces "github.com/shoplineapp/captin/interfaces"
	documentStores "github.com/shoplineapp/captin/internal/document_stores"
	helpers "github.com/shoplineapp/captin/internal/helpers"
	models "github.com/shoplineapp/captin/models"
	log "github.com/sirupsen/logrus"
)

var dLogger = log.WithFields(log.Fields{"class": "Dispatcher"})
var nullDocumentStore = documentStores.NewNullDocumentStore()

// Dispatcher - Event Dispatcher
type Dispatcher struct {
	destinations   []models.Destination
	senderMapping  map[string]interfaces.EventSenderInterface
	Errors         []interfaces.ErrorInterface
	targetDocument map[string]interface{}
	filters        []destination_filters.DestinationFilterInterface
	middlewares    []destination_filters.DestinationMiddlewareInterface
	errorHandler   interfaces.ErrorHandlerInterface
	delayer        interfaces.DispatchDelayerInterface

	muTargetDocument sync.Mutex
	muErrors         sync.Mutex
}

// NewDispatcherWithDestinations - Create Outgoing event dispatcher with destinations
func NewDispatcherWithDestinations(
	destinations []models.Destination,
	senderMapping map[string]interfaces.EventSenderInterface) *Dispatcher {
	result := Dispatcher{
		destinations:  destinations,
		senderMapping: senderMapping,
		filters:       []destination_filters.DestinationFilterInterface{},
		middlewares:   []destination_filters.DestinationMiddlewareInterface{},
		Errors:        []interfaces.ErrorInterface{},
	}

	return &result
}

// SetFilters - Add filters before dispatch
func (d *Dispatcher) SetFilters(filters []destination_filters.DestinationFilterInterface) {
	d.filters = filters
}

func (d *Dispatcher) SetMiddlewares(middlewares []destination_filters.DestinationMiddlewareInterface) {
	d.middlewares = middlewares
}

func (d *Dispatcher) SetErrorHandler(handler interfaces.ErrorHandlerInterface) {
	d.errorHandler = handler
}

func (d *Dispatcher) SetDelayer(delayer interfaces.DispatchDelayerInterface) {
	d.delayer = delayer
}

func (d *Dispatcher) GetErrors() []interfaces.ErrorInterface {
	d.muErrors.Lock()
	defer d.muErrors.Unlock()
	return d.Errors
}

func (d *Dispatcher) OnError(evt interfaces.IncomingEventInterface, err interfaces.ErrorInterface) {
	d.muErrors.Lock()
	defer d.muErrors.Unlock()
	d.Errors = append(d.Errors, err)

	switch dispatcherErr := err.(type) {
	case *captin_errors.DispatcherError:
		dLogger.WithFields(log.Fields{
			"event":       dispatcherErr.Event,
			"destination": dispatcherErr.Destination,
			"reason":      dispatcherErr.Error(),
		}).Error("Failed to dispatch event")
		d.TriggerErrorHandler(dispatcherErr)
	default:
		dLogger.WithFields(log.Fields{"event": evt, "error": err}).Error("Unhandled error on dispatcher")
	}
}

// Dispatch - Dispatch an event to outgoing webhook
func (d *Dispatcher) Dispatch(
	event interfaces.IncomingEventInterface,
	store interfaces.StoreInterface,
	throttler interfaces.ThrottleInterface,
	documentStoreMappings map[string]interfaces.DocumentStoreInterface,
) interfaces.ErrorInterface {
	responses := make(chan int, len(d.destinations))
	e := event.(models.IncomingEvent)
	for _, destination := range d.destinations {
		config := destination.Config
		canTrigger, timeRemain, err := throttler.CanTrigger(getEventKey(store, e, destination), config.GetThrottleValue())
		documentStore := d.getDocumentStore(destination, documentStoreMappings)

		if err != nil {
			dLogger.WithFields(log.Fields{"event": e.GetTraceInfo(), "destination": destination, "error": err}).Error("Error on getting throttle key")

			// Send without throttling
			go func(e models.IncomingEvent, destination models.Destination, documentStore interfaces.DocumentStoreInterface) {
				d.sendEvent(e, destination, store, documentStore)
				responses <- 1
			}(e, destination, documentStore)
			continue
		}

		if canTrigger {
			go func(e models.IncomingEvent, destination models.Destination, documentStore interfaces.DocumentStoreInterface) {
				d.sendEvent(e, destination, store, documentStore)
				responses <- 1
			}(e, destination, documentStore)
		} else if !config.GetThrottleTrailingDisabled() {
			go func(e models.IncomingEvent, destination models.Destination, documentStore interfaces.DocumentStoreInterface) {
				d.processDelayedEvent(e, timeRemain, destination, store, documentStore)
				responses <- 1
			}(e, destination, documentStore)
		} else {
			dLogger.WithFields(log.Fields{"event": e.GetTraceInfo(), "destination": destination}).Info("Cannot trigger send event")
			responses <- 0
		}
	}
	// Wait for destination completion
	for range d.destinations {
		<-responses
	}
	return nil
}

func (d *Dispatcher) TriggerErrorHandler(err *captin_errors.DispatcherError) {
	if d.errorHandler != nil {
		dispatcher.TrackGoRoutine(func() {
			d.errorHandler.Exec(*err)
		})
	}
}

// Private Functions

func (d *Dispatcher) getDocumentStore(dest models.Destination, documentStoreMappings map[string]interfaces.DocumentStoreInterface) interfaces.DocumentStoreInterface {
	if documentStoreMappings[dest.GetDocumentStore()] != nil {
		return documentStoreMappings[dest.GetDocumentStore()]
	}
	return nullDocumentStore
}

// inject document and sanitize fields in event based on destination
func (d *Dispatcher) customizeEvent(e models.IncomingEvent, destination models.Destination, documentStore interfaces.DocumentStoreInterface) interfaces.IncomingEventInterface {
	customized := e

	customized.TargetDocument = d.customizeDocument(&customized, destination, documentStore)
	customized.Payload = d.customizePayload(customized, destination)
	return customized
}

// inject throttled payloads from store if keep_throttled_payloads is true
func (d *Dispatcher) injectThrottledPayloads(e models.IncomingEvent, destination models.Destination, store interfaces.StoreInterface) interfaces.IncomingEventInterface {
	if destination.Config.GetKeepThrottledPayloads() {
		queueKey := getEventThrottledPayloadsKey(store, e, destination)
		payloadStrings, _, _, _ := store.GetQueue(queueKey)
		store.Remove(queueKey)
		for _, payloadStr := range payloadStrings {
			payload := map[string]interface{}{}
			json.Unmarshal([]byte(payloadStr), &payload)
			e.ThrottledPayloads = append(e.ThrottledPayloads, payload)
		}
	}
	return e
}

// inject throttled documents from store if include_document and keep_throttled_documents is true
func (d *Dispatcher) injectThrottledDocuments(e models.IncomingEvent, destination models.Destination, store interfaces.StoreInterface) interfaces.IncomingEventInterface {
	if destination.Config.GetIncludeDocument() && destination.Config.GetKeepThrottledDocuments() {
		queueKey := getEventThrottledDocumentsKey(store, e, destination)
		documentStrings, _, _, _ := store.GetQueue(queueKey)
		store.Remove(queueKey)
		for _, documentStr := range documentStrings {
			document := map[string]interface{}{}
			json.Unmarshal([]byte(documentStr), &document)
			e.ThrottledDocuments = append(e.ThrottledDocuments, document)
		}
	}
	return e
}

func (d *Dispatcher) customizeDocument(e *models.IncomingEvent, destination models.Destination, documentStore interfaces.DocumentStoreInterface) map[string]interface{} {
	config := destination.Config
	if config.GetIncludeDocument() == false {
		return e.TargetDocument
	}

	d.muTargetDocument.Lock()
	defer d.muTargetDocument.Unlock()

	// memoize document to be used across events for diff. destinations
	if d.targetDocument == nil {
		d.targetDocument = documentStore.GetDocument(*e)
	}

	if len(config.GetIncludeDocumentAttrs()) >= 1 {
		return helpers.IncludeFields(d.targetDocument, config.GetIncludeDocumentAttrs()).(map[string]interface{})
	} else if len(config.GetExcludeDocumentAttrs()) >= 1 {
		return helpers.ExcludeFields(d.targetDocument, config.GetExcludeDocumentAttrs()).(map[string]interface{})
	} else {
		return d.targetDocument
	}
}

func (d *Dispatcher) customizePayload(e models.IncomingEvent, destination interfaces.DestinationInterface) map[string]interface{} {
	config := destination.(models.Destination).Config
	if len(config.GetIncludePayloadAttrs()) >= 1 {
		return helpers.IncludeFields(e.Payload, config.GetIncludePayloadAttrs()).(map[string]interface{})
	} else if len(config.GetExcludePayloadAttrs()) >= 1 {
		return helpers.ExcludeFields(e.Payload, config.GetExcludePayloadAttrs()).(map[string]interface{})
	}

	return e.Payload
}

func (d *Dispatcher) processDelayedEvent(e models.IncomingEvent, timeRemain time.Duration, dest models.Destination, store interfaces.StoreInterface, documentStore interfaces.DocumentStoreInterface) {
	defer func() {
		if err := recover(); err != nil {
			d.OnError(e, &captin_errors.DispatcherError{
				Msg:         err.(error).Error(),
				Destination: dest,
				Event:       e,
			})
		}
	}()

	// Check if store have payload
	dataKey := getEventDataKey(store, e, dest)
	storedData, dataExists, _, storeErr := store.Get(dataKey)
	if storeErr != nil {
		panic(storeErr)
	}

	var storedEvent models.IncomingEvent
	if dataExists {
		storedEvent = models.IncomingEvent{}
		json.Unmarshal([]byte(storedData), &storedEvent)
	}

	if dataExists && getControlTimestamp(storedEvent, 0) > getControlTimestamp(e, uint64(time.Now().UnixNano())) {
		// Skip updating event data as stored data has newer timestamp
		dLogger.WithFields(log.Fields{
			"storedEvent":  storedEvent,
			"event":        e.GetTraceInfo(),
			"eventDataKey": "dataKey",
		}).Debug("Skipping update on event data")
		return
	}

	if dest.Config.GetKeepThrottledPayloads() {
		customizedPayload := d.customizePayload(e, dest)
		queueKey := getEventThrottledPayloadsKey(store, e, dest)
		jsonString, jsonErr := json.Marshal(customizedPayload)
		if jsonErr != nil {
			panic(jsonErr)
		}
		dLogger.WithFields(log.Fields{
			"queueKey":       queueKey,
			"event":          e.GetTraceInfo(),
			"enqueuePayload": jsonString,
		}).Debug("Storing throttled payload")
		store.Enqueue(queueKey, string(jsonString), dest.Config.GetThrottleValue()*2)
	}

	if dest.Config.GetIncludeDocument() && dest.Config.GetKeepThrottledDocuments() {
		customizedDocument := d.customizeDocument(&e, dest, documentStore)
		queueKey := getEventThrottledDocumentsKey(store, e, dest)
		jsonString, jsonErr := json.Marshal(customizedDocument)
		if jsonErr != nil {
			panic(jsonErr)
		}
		dLogger.WithFields(log.Fields{
			"queueKey":        queueKey,
			"event":           e.GetTraceInfo(),
			"enqueueDocument": jsonString,
		}).Debug("Storing throttled document")
		store.Enqueue(queueKey, string(jsonString), dest.Config.GetThrottleValue()*2)
	}

	jsonString, jsonErr := json.Marshal(e)
	if jsonErr != nil {
		panic(jsonErr)
	}

	if dataExists {
		// Update Value
		_, updateErr := store.Update(dataKey, string(jsonString))
		if updateErr != nil {
			panic(updateErr)
		}
	} else {
		// Create Value
		config := dest.Config
		_, saveErr := store.Set(dataKey, string(jsonString), config.GetThrottleValue()*2)
		if saveErr != nil {
			panic(saveErr)
		}

		// Schedule send event later
		dispatcher.TrackAfterFuncJob(timeRemain, func() {
			dLogger.WithFields(log.Fields{"key": dataKey}).Debug("After event callback")
			payload, exists, _, _ := store.Get(dataKey)
			// Key might be deleted by another worker, resulting in data not found
			if !exists {
				dLogger.WithFields(log.Fields{"key": dataKey}).Debug("Event data not found")
				d.OnError(models.IncomingEvent{}, &captin_errors.UnretryableError{Msg: "Event data not found", Event: e, Destination: dest})
				return
			}
			event := models.IncomingEvent{}
			json.Unmarshal([]byte(payload), &event)
			d.sendEvent(event, dest, store, documentStore)

			/* This line is commented out because their is a possibility of race condition,
			when the dataKey is removed by one worker, another worker might try to access it,
			it is better to let the dataKey expire by the ttl */
			// store.Remove(dataKey)
		})
	}
}

func getControlTimestamp(e models.IncomingEvent, defaultValue uint64) uint64 {
	defer func(d uint64) uint64 {
		if err := recover(); err != nil {
			return d
		}
		return 0
	}(defaultValue)

	value := e.Control["ts"]

	// Type assertion from interface
	switch v := value.(type) {
	case int:
		value = uint64(v)
	case string:
		parsed, err := strconv.ParseUint(v, 10, 64)
		if err != nil {
			panic("unable to convert string timestamp")
		}
		value = parsed
	}

	return value.(uint64)
}

func getEventKey(s interfaces.StoreInterface, e interfaces.IncomingEventInterface, d interfaces.DestinationInterface) string {
	return s.DataKey(e, d, "", "")
}

func getEventDataKey(s interfaces.StoreInterface, e interfaces.IncomingEventInterface, d interfaces.DestinationInterface) string {
	return s.DataKey(e, d, "", "-data")
}

func getEventThrottledPayloadsKey(s interfaces.StoreInterface, e models.IncomingEvent, d models.Destination) string {
	return s.DataKey(e, d, "", "-throttled_payloads")
}

func getEventThrottledDocumentsKey(s interfaces.StoreInterface, e models.IncomingEvent, d models.Destination) string {
	return s.DataKey(e, d, "", "-throttled_documents")
}

func (d *Dispatcher) sendEvent(evt models.IncomingEvent, destination models.Destination, store interfaces.StoreInterface, documentStore interfaces.DocumentStoreInterface) {
	config := destination.Config
	callbackLogger := dLogger.WithFields(log.Fields{
		"action":         evt.Key,
		"event":          evt.GetTraceInfo(),
		"hook_name":      config.GetName(),
		"callback_url":   destination.GetCallbackURL(),
		"document_store": destination.GetDocumentStore(),
	})

	defer func() {
		if err := recover(); err != nil {
			callbackLogger.Info(fmt.Sprintf("Event failed sending to %s [%s]", config.GetName(), destination.GetCallbackURL()))
			d.OnError(evt, &captin_errors.DispatcherError{
				Msg:         err.(error).Error(),
				Destination: destination,
				Event:       evt,
			})
		}
		return
	}()

	callbackLogger.Debug("Preprocess payload and document")

	evt = d.customizeEvent(evt, destination, documentStore).(models.IncomingEvent)

	evt = d.injectThrottledPayloads(evt, destination, store).(models.IncomingEvent)

	evt = d.injectThrottledDocuments(evt, destination, store).(models.IncomingEvent)

	callbackLogger.Debug("Final sift on dispatcher")

	sifted := Custom{}.Sift(&evt, []models.Destination{destination}, d.filters, d.middlewares)
	if len(sifted) == 0 {
		callbackLogger.Info("Event interrupted by dispatcher filters")
		return
	}

	callbackLogger.Debug("Ready to send event")

	senderKey := config.GetSender()
	if senderKey == "" {
		senderKey = "http"
	}
	sender, senderExists := d.senderMapping[senderKey]
	if senderExists == false {
		panic(&captin_errors.DispatcherError{
			Msg:         fmt.Sprintf("Sender key %s does not exist", senderKey),
			Destination: destination,
			Event:       evt,
		})
	}

	// Wrap event sender and error handling as closure for reusing in delayer
	_sendEvent := func() {
		defer func() {
			if err := recover(); err != nil {
				var newErr error
				switch err := err.(type) {
				// As the event is invalid, this error is raised so that the event is not retried
				case *captin_errors.UnretryableError:
					newErr = err
				default:
					newErr = &captin_errors.DispatcherError{
						Msg:         err.(error).Error(),
						Destination: destination,
						Event:       evt,
					}
				}
				d.OnError(evt, newErr)
			}
		}()
		// Deep clone a new instance to prevent concurrent iteration and write on json.Marshal
		event := deepcopy.Copy(evt).(models.IncomingEvent)
		err := sender.SendEvent(event, destination)
		if err != nil {
			panic(err)
		}
		callbackLogger.Info(fmt.Sprintf("Event successfully sent to %s [%s]", config.GetName(), destination.GetCallbackURL()))
	}

	if destination.RequireDelay(evt) {
		// Sending message with delay in goroutine, no error will be caught
		callbackLogger.Info(fmt.Sprintf("Event requires delay"))
		if d.delayer != nil {
			// Delayer will usually modify event.Control for delay info, deep clone to prevent concurrent write
			event := deepcopy.Copy(evt).(models.IncomingEvent)
			d.delayer.Execute(event, destination, _sendEvent)
			return
		} else {
			callbackLogger.Warn(fmt.Sprintf("Delayer not found, send event immediately"))
		}
	}

	_sendEvent()
}
