package outgoing

import (
	"encoding/json"
	"fmt"
	"time"

	interfaces "github.com/shoplineapp/captin/interfaces"
	models "github.com/shoplineapp/captin/models"
)

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
	destinations []models.Destination
	sender       interfaces.EventSenderInterface
	Errors       []error
}

// NewDispatcherWithDestinations - Create Outgoing event dispatcher with destinations
func NewDispatcherWithDestinations(
	destinations []models.Destination,
	sender interfaces.EventSenderInterface) *Dispatcher {
	result := Dispatcher{
		destinations: destinations,
		sender:       sender,
		Errors:       []error{},
	}

	return &result
}

// Dispatch - Dispatch an event to outgoing webhook
func (d *Dispatcher) Dispatch(
	e models.IncomingEvent,
	store interfaces.StoreInterface,
	throttler interfaces.ThrottleInterface) error {

	for _, destination := range d.destinations {
		canTrigger, timeRemain, err := throttler.CanTrigger(getEventKey(e, destination), destination.Config.GetThrottleValue())

		if err != nil {
			fmt.Println("[Dispatcher] Error: ", err)
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
	dataKey := getEventDataKey(e, dest)
	_, ok, _, storeErr := store.Get(dataKey)
	if storeErr != nil {
		panic(storeErr)
	}

	jsonString, jsonErr := json.Marshal(e)
	if jsonErr != nil {
		panic(jsonErr)
	}

	if ok {
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

func getEventKey(e models.IncomingEvent, d models.Destination) string {
	return fmt.Sprintf("%s.%s.%s", e.Key, d.Config.Name, e.TargetId)
}

func getEventDataKey(e models.IncomingEvent, d models.Destination) string {
	return fmt.Sprintf("%s.%s.%s-data", e.Key, d.Config.Name, e.TargetId)
}

func (d *Dispatcher) sendEvent(evt models.IncomingEvent, destination models.Destination) {
	err := d.sender.SendEvent(evt, destination)
	if err != nil {
		d.Errors = append(d.Errors, &DispatcherError{
			msg:         err.Error(),
			Destination: destination,
			Event:       evt,
		})
	}
}

func (d *Dispatcher) sendAfterEvent(key string, store interfaces.StoreInterface, dest models.Destination) func() {
	fmt.Println("[Dispatcher]: sendAfterEvent: ", key)
	payload, _, _, _ := store.Get(key)
	event := models.IncomingEvent{}
	json.Unmarshal([]byte(payload), &event)
	return func() {
		d.sendEvent(event, dest)
		store.Remove(key)
	}
}
