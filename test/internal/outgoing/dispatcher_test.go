package outgoing_test

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	interfaces "github.com/shoplineapp/captin/interfaces"
	outgoing "github.com/shoplineapp/captin/internal/outgoing"
	models "github.com/shoplineapp/captin/models"
	mocks "github.com/shoplineapp/captin/test/mocks"

	"github.com/stretchr/testify/mock"
)

func setup(path string) (*mocks.StoreMock, *mocks.SenderMock, *outgoing.Dispatcher, *mocks.ThrottleMock) {
	store := new(mocks.StoreMock)
	sender := new(mocks.SenderMock)
	senderMapping := map[string]interfaces.EventSenderInterface{
		"mock": sender,
	}
	throttler := new(mocks.ThrottleMock)

	data, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	configs := []models.Configuration{}
	json.Unmarshal(data, &configs)
	destinations := []models.Destination{}
	for _, config := range configs {
		destinations = append(destinations, models.Destination{Config: config})
	}

	dispatcher := outgoing.NewDispatcherWithDestinations(destinations, senderMapping)

	return store, sender, dispatcher, throttler
}

func TestDispatchEvents_Error(t *testing.T) {
	store, sender, dispatcher, throttler := setup("fixtures/config.json")

	sender.On("SendEvent", mock.Anything, mock.Anything).Return(errors.New("Mock Error"))
	store.On("Get", mock.Anything).Return("", false, time.Duration(0), nil)
	store.On("Set", mock.Anything, mock.Anything, mock.Anything).Return(true, nil)
	throttler.On("CanTrigger", mock.Anything, mock.Anything).Return(true, time.Duration(0), nil)

	dispatcher.Dispatch(models.IncomingEvent{
		Key:        "product.update",
		Source:     "core",
		Payload:    map[string]interface{}{"field1": 1},
		TargetType: "Product",
		TargetId:   "product_id",
	}, store, throttler)

	dispatcher.Dispatch(models.IncomingEvent{
		Key:        "product.update",
		Source:     "core",
		Payload:    map[string]interface{}{"field1": 2},
		TargetType: "Product",
		TargetId:   "product_id_2",
	}, store, throttler)

	time.Sleep(50 * time.Millisecond)

	sender.AssertNumberOfCalls(t, "SendEvent", 6)
	assert.Equal(t, 6, len(dispatcher.Errors))
}

func TestDispatchEvents(t *testing.T) {
	store, sender, dispatcher, throttler := setup("fixtures/config.json")

	sender.On("SendEvent", mock.Anything, mock.Anything).Return(nil)
	store.On("Get", mock.Anything).Return("", false, time.Duration(0), nil)
	store.On("Set", mock.Anything, mock.Anything, mock.Anything).Return(true, nil)
	throttler.On("CanTrigger", mock.Anything, mock.Anything).Return(true, time.Duration(0), nil)

	dispatcher.Dispatch(models.IncomingEvent{
		Key:        "product.update",
		Source:     "core",
		Payload:    map[string]interface{}{"field1": 1},
		TargetType: "Product",
		TargetId:   "product_id",
	}, store, throttler)

	dispatcher.Dispatch(models.IncomingEvent{
		Key:        "product.update",
		Source:     "core",
		Payload:    map[string]interface{}{"field1": 2},
		TargetType: "Product",
		TargetId:   "product_id_2",
	}, store, throttler)

	time.Sleep(50 * time.Millisecond)

	sender.AssertNumberOfCalls(t, "SendEvent", 6)
}

func TestDispatchEvents_Throttled_DelaySend(t *testing.T) {
	store, sender, dispatcher, throttler := setup("fixtures/config.single.json")

	sender.On("SendEvent", mock.Anything, mock.Anything).Return(nil)
	store.On("Get", mock.Anything).Return("", false, time.Duration(0), nil)
	store.On("Set", mock.Anything, mock.Anything, mock.Anything).Return(true, nil)
	store.On("Remove", mock.Anything).Return(true, nil)
	throttler.On("CanTrigger", mock.Anything, mock.Anything).Return(false, 500*time.Millisecond, nil)

	dispatcher.Dispatch(models.IncomingEvent{
		Key:        "product.update",
		Source:     "core",
		Payload:    map[string]interface{}{"field1": 1},
		TargetType: "Product",
		TargetId:   "product_id",
	}, store, throttler)

	sender.AssertNumberOfCalls(t, "SendEvent", 0)

	throttleID := "product.update.service_one.product_id-data"
	throttlePeriod := time.Millisecond * 500 * 2
	store.AssertCalled(t, "Get", throttleID)
	store.AssertCalled(t, "Set", throttleID, `{"event_key":"product.update","source":"core","payload":{"field1":1},"control":null,"target_type":"Product","target_id":"product_id"}`, throttlePeriod)

	time.Sleep(600 * time.Millisecond)

	sender.AssertNumberOfCalls(t, "SendEvent", 1)
	store.AssertCalled(t, "Remove", throttleID)
}

func TestDispatchEvents_Throttled_SkipUpdatingValue(t *testing.T) {
	// throttleTimestamp := time.Now().Unix()
	store, sender, dispatcher, throttler := setup("fixtures/config.single.json")
	throttleID := "product.update.service_one.product_id-data"
	throttlePeriod := time.Millisecond * 700 * 2

	sender.On("SendEvent", mock.Anything, mock.Anything).Return(nil)
	store.On("Get", throttleID).Return(`{"event_key":"product.update","source":"core","payload":{"field1":1},"control":{"ts":"99999999999999"},"target_type":"Product","target_id":"product_id"}`, true, throttlePeriod, nil)
	store.On("Set", mock.Anything, mock.Anything, mock.Anything).Return(true, nil)
	store.On("Update", mock.Anything, mock.Anything).Return(true, nil)
	throttler.On("CanTrigger", mock.Anything, mock.Anything).Return(false, time.Millisecond*700, nil)

	dispatcher.Dispatch(models.IncomingEvent{
		Key:        "product.update",
		Source:     "core",
		Payload:    map[string]interface{}{"field1": 2},
		Control:    map[string]interface{}{"ts": 1572830980},
		TargetType: "Product",
		TargetId:   "product_id",
	}, store, throttler)

	sender.AssertNumberOfCalls(t, "SendEvent", 0)
	store.AssertNumberOfCalls(t, "Update", 0)
	store.AssertCalled(t, "Get", throttleID)
}

func TestDispatchEvents_Throttled_UpdatePayload(t *testing.T) {
	store, sender, dispatcher, throttler := setup("fixtures/config.single.json")
	throttleID := "product.update.service_one.product_id-data"
	throttlePeriod := time.Millisecond * 700 * 2

	sender.On("SendEvent", mock.Anything, mock.Anything).Return(nil)
	store.On("Get", throttleID).Return(`{"event_key":"product.update","source":"core","payload":{"field1":1},"control":null,"target_type":"Product","target_id":"product_id"}`, true, throttlePeriod, nil)
	store.On("Set", mock.Anything, mock.Anything, mock.Anything).Return(true, nil)
	store.On("Update", mock.Anything, mock.Anything).Return(true, nil)
	store.On("Remove", mock.Anything).Return(true, nil)
	throttler.On("CanTrigger", mock.Anything, mock.Anything).Return(false, time.Millisecond*700, nil)

	dispatcher.Dispatch(models.IncomingEvent{
		Key:        "product.update",
		Source:     "core",
		Payload:    map[string]interface{}{"field1": 2},
		TargetType: "Product",
		TargetId:   "product_id",
	}, store, throttler)

	sender.AssertNumberOfCalls(t, "SendEvent", 0)
	store.AssertCalled(t, "Get", throttleID)
	store.AssertCalled(t, "Update", throttleID, `{"event_key":"product.update","source":"core","payload":{"field1":2},"control":null,"target_type":"Product","target_id":"product_id"}`)
}
