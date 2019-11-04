package outgoing

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	interfaces "github.com/shoplineapp/captin/interfaces"
	models "github.com/shoplineapp/captin/models"
	log "github.com/sirupsen/logrus"
)

var dLogger = log.WithFields(log.Fields{"class": "Dispatcher"})

// DispatcherError - Error when send events
type DispatcherError struct {
	msg         string
	Event       models.IncomingEvent
	Destination models.Destination
}

func (e *DispatcherError) Error() string {
	return e.msg
}

// Dispatcher - Event Dispatcher
type Dispatcher struct {
	destinations  []models.Destination
	senderMapping map[string]interfaces.EventSenderInterface
	Errors        []error
}

// NewDispatcherWithDestinations - Create Outgoing event dispatcher with destinations
func NewDispatcherWithDestinations(
	destinations []models.Destination,
	senderMapping map[string]interfaces.EventSenderInterface) *Dispatcher {
	result := Dispatcher{
		destinations:  destinations,
		senderMapping: senderMapping,
		Errors:        []error{},
	}

	return &result
}

// Dispatch - Dispatch an event to outgoing webhook
func (d *Dispatcher) Dispatch(
	e models.IncomingEvent,
	store interfaces.StoreInterface,
	throttler interfaces.ThrottleInterface) error {

	for _, destination := range d.destinations {
		canTrigger, timeRemain, err := throttler.CanTrigger(getEventKey(store, e, destination), destination.Config.GetThrottleValue())

		if err != nil {
			dLogger.WithFields(log.Fields{"event": e, "destination": destination}).Error("Failed to dispatch event")
			d.Errors = append(d.Errors, err)

			// Send without throttling
			go d.sendEvent(e, destination)
			continue
		}

		if canTrigger {
			go d.sendEvent(e, destination)
		} else {
			d.processDelayedEvent(e, timeRemain, destination, store)
		}
	}
	return nil
}

// Private Functions

func (d *Dispatcher) processDelayedEvent(e models.IncomingEvent, timeRemain time.Duration, dest models.Destination, store interfaces.StoreInterface) {
	defer func() {
		if err := recover(); err != nil {
			d.Errors = append(d.Errors, &DispatcherError{
				msg:         err.(error).Error(),
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
		if getControlTimestamp(storedEvent, 0) > getControlTimestamp(e, time.Now().Unix()) {
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
		time.AfterFunc(timeRemain, d.sendAfterEvent(dataKey, store, dest))
	}
}

func getControlTimestamp(e models.IncomingEvent, defaultValue int64) int64 {
	defer func(d int64) int64 {
		if err := recover(); err != nil {
			return d
		}
		return 0
	}(defaultValue)

	value := e.Control["ts"]

	// Type assertion from interface
	switch v := value.(type) {
	case int:
		value = int64(v)
	case string:
		parsed, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			panic("unable to convert string timestamp")
		}
		value = parsed
	}

	return value.(int64)
}

func getEventKey(s interfaces.StoreInterface, e models.IncomingEvent, d models.Destination) string {
	return s.DataKey(e, d, "", "")
}

func getEventDataKey(s interfaces.StoreInterface, e models.IncomingEvent, d models.Destination) string {
	return s.DataKey(e, d, "", "-data")
}

func (d *Dispatcher) sendEvent(evt models.IncomingEvent, destination models.Destination) {
	callbackLogger := dLogger.WithFields(log.Fields{
		"callback_url": destination.Config.CallbackURL,
	})
	callbackLogger.Debug("Ready to send event")

	senderKey := destination.Config.Sender
	if senderKey == "" {
		senderKey = "http"
	}
	sender, senderExists := d.senderMapping[senderKey]
	if senderExists == false {
		d.Errors = append(d.Errors, &DispatcherError{
			msg:         fmt.Sprintf("Sender key %s does not exist", senderKey),
			Destination: destination,
			Event:       evt,
		})
		return
	}

	err := sender.SendEvent(evt, destination)
	if err != nil {
		d.Errors = append(d.Errors, &DispatcherError{
			msg:         err.Error(),
			Destination: destination,
			Event:       evt,
		})
		return
	}
	callbackLogger.Info(fmt.Sprintf("Event successfully sent to %s", destination.Config.CallbackURL))
}

func (d *Dispatcher) sendAfterEvent(key string, store interfaces.StoreInterface, dest models.Destination) func() {
	dLogger.WithFields(log.Fields{"key": key}).Debug("After event callback")
	payload, _, _, _ := store.Get(key)
	event := models.IncomingEvent{}
	json.Unmarshal([]byte(payload), &event)
	return func() {
		d.sendEvent(event, dest)
		store.Remove(key)
	}
}
