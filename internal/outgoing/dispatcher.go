package outgoing

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	captin_errors "github.com/shoplineapp/captin/errors"
	interfaces "github.com/shoplineapp/captin/interfaces"
	destination_filters "github.com/shoplineapp/captin/destinations/filters"
	models "github.com/shoplineapp/captin/models"
	documentStores "github.com/shoplineapp/captin/internal/document_stores"
	helpers "github.com/shoplineapp/captin/internal/helpers"
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
			responses <- 0
			d.processDelayedEvent(e, timeRemain, destination, store, documentStore)
		} else {
			responses <- 0
		}
	}
	// Wait for destination completion
	for _, _ = range d.destinations {
		<-responses
	}
	return nil
}

func (d Dispatcher) TriggerErrorHandler(err *captin_errors.DispatcherError) {
	if d.errorHandler != nil {
		go d.errorHandler.Exec(*err)
	}
}

// Private Functions

func (d Dispatcher) getDocumentStore(dest models.Destination, documentStoreMappings map[string]interfaces.DocumentStoreInterface) interfaces.DocumentStoreInterface {
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

func (d *Dispatcher) customizePayload(e models.IncomingEvent, destination interfaces.DestinationInterface) map[string]interface{}  {
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
			d.Errors = append(d.Errors, &captin_errors.DispatcherError{
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
			"queueKey":  queueKey,
			"event":        e.GetTraceInfo(),
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
			"queueKey":  queueKey,
			"event":        e.GetTraceInfo(),
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
		time.AfterFunc(timeRemain, d.sendAfterEvent(dataKey, store, dest, documentStore))
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
			d.Errors = append(d.Errors, &captin_errors.DispatcherError{
				Msg:         fmt.Sprintf("%+v", err),
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
		return
	}

	if config.GetDelayValue() != time.Duration(0) {
		// Sending message with delay in goroutine, no error will be caught
		callbackLogger.Info(fmt.Sprintf("Event delayed with %s", config.GetDelay()))
		ch := make(chan int, 1)
		go time.AfterFunc(config.GetDelayValue(), func() {
			delayedErr := sender.SendEvent(evt, destination)
			if delayedErr != nil {
				callbackLogger.WithFields(log.Fields{"error": delayedErr}).Error(fmt.Sprintf("Delayed event failed with error on %s [%s]", config.GetName(), destination.GetCallbackURL()))
				d.TriggerErrorHandler(&captin_errors.DispatcherError{
					Msg:         delayedErr.Error(),
					Destination: destination,
					Event:       evt,
				})
			} else {
				callbackLogger.Info(fmt.Sprintf("Event successfully sent to %s [%s]", config.GetName(), destination.GetCallbackURL()))
			}
			ch <- 1
		})
		<-ch // waiting for delayed execution
		return
	}

	err := sender.SendEvent(evt, destination)
	if err != nil {
		panic(&captin_errors.DispatcherError{
			Msg:         err.Error(),
			Destination: destination,
			Event:       evt,
		})
		return
	}

	callbackLogger.Info(fmt.Sprintf("Event successfully sent to %s [%s]", config.GetName(), destination.GetCallbackURL()))
}

func (d *Dispatcher) sendAfterEvent(key string, store interfaces.StoreInterface, dest models.Destination, documentStore interfaces.DocumentStoreInterface) func() {
	return func() {
		dLogger.WithFields(log.Fields{"key": key}).Debug("After event callback")
		payload, _, _, _ := store.Get(key)
		event := models.IncomingEvent{}
		json.Unmarshal([]byte(payload), &event)
		d.sendEvent(event, dest, store, documentStore)
		store.Remove(key)
	}
}
