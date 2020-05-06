package outgoing

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	captin_errors "github.com/shoplineapp/captin/errors"
	interfaces "github.com/shoplineapp/captin/interfaces"
	models "github.com/shoplineapp/captin/models"
	log "github.com/sirupsen/logrus"
)

var dLogger = log.WithFields(log.Fields{"class": "Dispatcher"})

// Dispatcher - Event Dispatcher
type Dispatcher struct {
	destinations   []models.Destination
	senderMapping  map[string]interfaces.EventSenderInterface
	Errors         []captin_errors.ErrorInterface
	targetDocument map[string]interface{}
}

// NewDispatcherWithDestinations - Create Outgoing event dispatcher with destinations
func NewDispatcherWithDestinations(
	destinations []models.Destination,
	senderMapping map[string]interfaces.EventSenderInterface) *Dispatcher {
	result := Dispatcher{
		destinations:  destinations,
		senderMapping: senderMapping,
		Errors:        []captin_errors.ErrorInterface{},
	}

	return &result
}

// Dispatch - Dispatch an event to outgoing webhook
func (d *Dispatcher) Dispatch(
	e models.IncomingEvent,
	store interfaces.StoreInterface,
	throttler interfaces.ThrottleInterface,
	documentStore interfaces.DocumentStoreInterface,
) captin_errors.ErrorInterface {

	for _, destination := range d.destinations {
		canTrigger, timeRemain, err := throttler.CanTrigger(getEventKey(store, e, destination), destination.Config.GetThrottleValue())

		if err != nil {
			dLogger.WithFields(log.Fields{"event": e, "destination": destination}).Error("Failed to dispatch event")
			d.Errors = append(d.Errors, err)

			// Send without throttling
			go d.sendEvent(e, destination, documentStore)
			continue
		}

		if canTrigger {
			go d.sendEvent(e, destination, documentStore)
		} else if destination.Config.ThrottleTrailing {
			d.processDelayedEvent(e, timeRemain, destination, store, documentStore)
		}
	}
	return nil
}

// Private Functions
func (d *Dispatcher) cloneEventWithDocument(e models.IncomingEvent, destination models.Destination, documentStore interfaces.DocumentStoreInterface) models.IncomingEvent {
	if destination.Config.IncludeDocument == false {
		return e;
	}

	// memoize document to be used across events for diff. destinations
	if d.targetDocument == nil {
		d.targetDocument = documentStore.GetDocument(e)
	}

	clone := e
	clone.TargetDocument = d.targetDocument

	return clone
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

	jsonString, jsonErr := json.Marshal(e)
	if jsonErr != nil {
		panic(jsonErr)
	}

	if dataExists {
		storedEvent := models.IncomingEvent{}
		json.Unmarshal([]byte(storedData), &storedEvent)
		if getControlTimestamp(storedEvent, 0) > getControlTimestamp(e, uint64(time.Now().UnixNano())) {
			// Skip updating event data as stored data has newer timestamp
			dLogger.WithFields(log.Fields{
				"storedEvent":  storedEvent,
				"event":        e,
				"eventDataKey": "dataKey",
			}).Debug("Skipping update on event data")
			return
		}
	}

	if dataExists {
		// Update Value
		_, updateErr := store.Update(dataKey, string(jsonString))
		if updateErr != nil {
			panic(updateErr)
		}
	} else {
		// Create Value
		_, saveErr := store.Set(dataKey, string(jsonString), dest.Config.GetThrottleValue()*2)
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

func getEventKey(s interfaces.StoreInterface, e models.IncomingEvent, d models.Destination) string {
	return s.DataKey(e, d, "", "")
}

func getEventDataKey(s interfaces.StoreInterface, e models.IncomingEvent, d models.Destination) string {
	return s.DataKey(e, d, "", "-data")
}

func (d *Dispatcher) sendEvent(evt models.IncomingEvent, destination models.Destination, documentStore interfaces.DocumentStoreInterface) {
	callbackLogger := dLogger.WithFields(log.Fields{
		"callback_url": destination.Config.CallbackURL,
	})
	callbackLogger.Debug("Ready to send event")

	evtWithDoc := d.cloneEventWithDocument(evt, destination, documentStore)

	senderKey := destination.Config.Sender
	if senderKey == "" {
		senderKey = "http"
	}
	sender, senderExists := d.senderMapping[senderKey]
	if senderExists == false {
		d.Errors = append(d.Errors, &captin_errors.DispatcherError{
			Msg:         fmt.Sprintf("Sender key %s does not exist", senderKey),
			Destination: destination,
			Event:       evtWithDoc,
		})
		return
	}

	err := sender.SendEvent(evtWithDoc, destination)
	if err != nil {
		d.Errors = append(d.Errors, &captin_errors.DispatcherError{
			Msg:         err.Error(),
			Destination: destination,
			Event:       evtWithDoc,
		})
		return
	}
	callbackLogger.Info(fmt.Sprintf("Event successfully sent to %s", destination.Config.CallbackURL))
}

func (d *Dispatcher) sendAfterEvent(key string, store interfaces.StoreInterface, dest models.Destination, documentStore interfaces.DocumentStoreInterface) func() {
	dLogger.WithFields(log.Fields{"key": key}).Debug("After event callback")
	payload, _, _, _ := store.Get(key)
	event := models.IncomingEvent{}
	json.Unmarshal([]byte(payload), &event)
	return func() {
		d.sendEvent(event, dest, documentStore)
		store.Remove(key)
	}
}
