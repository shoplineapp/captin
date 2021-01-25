package outgoing_test

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	interfaces "github.com/shoplineapp/captin/interfaces"
	outgoing "github.com/shoplineapp/captin/internal/outgoing"
	models "github.com/shoplineapp/captin/models"
	mocks "github.com/shoplineapp/captin/test/mocks"

	"github.com/stretchr/testify/mock"
)

func setup(path string) (*mocks.StoreMock, map[string]interfaces.DocumentStoreInterface, *mocks.SenderMock, *outgoing.Dispatcher, *mocks.ThrottleMock) {
	store := new(mocks.StoreMock)
	documentstores := map[string]interfaces.DocumentStoreInterface{
		"default": new(mocks.DocumentStoreMock),
	}
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

	return store, documentstores, sender, dispatcher, throttler
}

func hasNoDocument(e models.IncomingEvent) bool {
	return e.TargetDocument == nil
}

func hasDocument(e models.IncomingEvent) bool {
	return hasNoDocument(e) != true
}

func TestDispatchEvents_Error(t *testing.T) {
	store, documentStores, sender, dispatcher, throttler := setup("fixtures/config.json")

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
	}, store, throttler, documentStores)

	dispatcher.Dispatch(models.IncomingEvent{
		Key:        "product.update",
		Source:     "core",
		Payload:    map[string]interface{}{"field1": 2},
		TargetType: "Product",
		TargetId:   "product_id_2",
	}, store, throttler, documentStores)

	time.Sleep(50 * time.Millisecond)

	sender.AssertNumberOfCalls(t, "SendEvent", 6)
	assert.Equal(t, 6, len(dispatcher.Errors))
}

func TestDispatchEvents(t *testing.T) {
	store, documentStores, sender, dispatcher, throttler := setup("fixtures/config.json")

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
	}, store, throttler, documentStores)

	dispatcher.Dispatch(models.IncomingEvent{
		Key:        "product.update",
		Source:     "core",
		Payload:    map[string]interface{}{"field1": 2},
		TargetType: "Product",
		TargetId:   "product_id_2",
	}, store, throttler, documentStores)

	time.Sleep(50 * time.Millisecond)

	sender.AssertNumberOfCalls(t, "SendEvent", 6)
}

func TestDispatchEvents_Throttled_DelaySend(t *testing.T) {
	store, documentStores, sender, dispatcher, throttler := setup("fixtures/config.single.json")

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
	}, store, throttler, documentStores)

	sender.AssertNumberOfCalls(t, "SendEvent", 0)

	throttleID := "product.update.service_one.product_id-data"
	throttlePeriod := time.Millisecond * 500 * 2
	store.AssertCalled(t, "Get", throttleID)
	store.AssertCalled(t, "Set", throttleID, `{"TraceId":"","event_key":"product.update","source":"core","payload":{"field1":1},"control":null,"target_type":"Product","target_id":"product_id"}`, throttlePeriod)

	time.Sleep(600 * time.Millisecond)

	sender.AssertNumberOfCalls(t, "SendEvent", 1)
	store.AssertCalled(t, "Remove", throttleID)
}

func TestDispatchEvents_Throttled_SkipUpdatingValue(t *testing.T) {
	store, documentStores, sender, dispatcher, throttler := setup("fixtures/config.single.json")
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
		Control:    map[string]interface{}{"ts": uint(time.Now().UnixNano())},
		TargetType: "Product",
		TargetId:   "product_id",
	}, store, throttler, documentStores)

	// Expecting update of store is skipped due to older timestamp from new incoming event
	sender.AssertNumberOfCalls(t, "SendEvent", 0)
	store.AssertNumberOfCalls(t, "Update", 0)
	store.AssertCalled(t, "Get", throttleID)
}

func TestDispatchEvents_Throttled_UpdatePayload(t *testing.T) {
	store, documentStores, sender, dispatcher, throttler := setup("fixtures/config.single.json")
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
	}, store, throttler, documentStores)

	sender.AssertNumberOfCalls(t, "SendEvent", 0)
	store.AssertCalled(t, "Get", throttleID)
	store.AssertCalled(t, "Update", throttleID, `{"TraceId":"","event_key":"product.update","source":"core","payload":{"field1":2},"control":null,"target_type":"Product","target_id":"product_id"}`)
}

func TestDispatchEvents_With_Document(t *testing.T) {
	store, documentStores, sender, dispatcher, throttler := setup("fixtures/config.include_document.json")
	mockDocumentStore := new(mocks.DocumentStoreMock)
	documentStores["default"] = mockDocumentStore

	sender.On("SendEvent", mock.MatchedBy(hasDocument), mock.Anything).Return(nil)
	sender.On("SendEvent", mock.MatchedBy(hasNoDocument), mock.Anything).Return(nil)
	store.On("Get", mock.Anything).Return("", false, time.Duration(0), nil)
	store.On("Set", mock.Anything, mock.Anything, mock.Anything).Return(true, nil)
	mockDocumentStore.On("GetDocument", mock.Anything).Return(map[string]interface{}{"foo": "bar"})

	throttler.On("CanTrigger", mock.Anything, mock.Anything).Return(true, time.Duration(0), nil)

	dispatcher.Dispatch(models.IncomingEvent{
		Key:        "product.update",
		Source:     "core",
		Payload:    map[string]interface{}{"field1": 1},
		TargetType: "Product",
		TargetId:   "product_id",
	}, store, throttler, documentStores)

	time.Sleep(50 * time.Millisecond)

	sender.AssertExpectations(t)
}

func TestDispatchEvents_With_Include_Document_Attrs(t *testing.T) {
	store, documentStores, sender, dispatcher, throttler := setup("fixtures/config.include_document_attrs.json")
	mockDocumentStore := new(mocks.DocumentStoreMock)
	documentStores["default"] = mockDocumentStore

	sender.On("SendEvent", mock.MatchedBy(func(e models.IncomingEvent) bool {
		return reflect.DeepEqual(e.TargetDocument, map[string]interface{}{"foo": map[string]interface{}{"bar": "yo"}})
	}), mock.Anything).Return(nil)
	store.On("Get", mock.Anything).Return("", false, time.Duration(0), nil)
	store.On("Set", mock.Anything, mock.Anything, mock.Anything).Return(true, nil)
	mockDocumentStore.On("GetDocument", mock.Anything).Return(map[string]interface{}{"foo": map[string]interface{}{"bar": "yo", "a": "b"}, "foo2": "bar2"})

	throttler.On("CanTrigger", mock.Anything, mock.Anything).Return(true, time.Duration(0), nil)

	dispatcher.Dispatch(models.IncomingEvent{
		Key:        "product.update",
		Source:     "core",
		Payload:    map[string]interface{}{"field1": 1},
		TargetType: "Product",
		TargetId:   "product_id",
	}, store, throttler, documentStores)

	time.Sleep(50 * time.Millisecond)

	sender.AssertExpectations(t)
}

func TestDispatchEvents_With_Exclude_Document_Attrs(t *testing.T) {
	store, documentStores, sender, dispatcher, throttler := setup("fixtures/config.exclude_document_attrs.json")
	mockDocumentStore := new(mocks.DocumentStoreMock)
	documentStores["default"] = mockDocumentStore

	sender.On("SendEvent", mock.MatchedBy(func(e models.IncomingEvent) bool {
		return reflect.DeepEqual(e.TargetDocument, map[string]interface{}{"foo": map[string]interface{}{"a": "b"}, "foo2": "bar2"})
	}), mock.Anything).Return(nil)
	store.On("Get", mock.Anything).Return("", false, time.Duration(0), nil)
	store.On("Set", mock.Anything, mock.Anything, mock.Anything).Return(true, nil)
	mockDocumentStore.On("GetDocument", mock.Anything).Return(map[string]interface{}{"foo": map[string]interface{}{"bar": "yo", "a": "b"}, "foo2": "bar2"})

	throttler.On("CanTrigger", mock.Anything, mock.Anything).Return(true, time.Duration(0), nil)

	dispatcher.Dispatch(models.IncomingEvent{
		Key:        "product.update",
		Source:     "core",
		Payload:    map[string]interface{}{"field1": 1},
		TargetType: "Product",
		TargetId:   "product_id",
	}, store, throttler, documentStores)

	time.Sleep(50 * time.Millisecond)

	sender.AssertExpectations(t)
}

func TestDispatchEvents_With_Include_Payload_Attrs(t *testing.T) {
	store, documentStores, sender, dispatcher, throttler := setup("fixtures/config.include_payload_attrs.json")

	sender.On("SendEvent", mock.MatchedBy(func(e models.IncomingEvent) bool {
		return reflect.DeepEqual(e.Payload, map[string]interface{}{"field1": 1})
	}), mock.Anything).Return(nil)
	store.On("Get", mock.Anything).Return("", false, time.Duration(0), nil)
	store.On("Set", mock.Anything, mock.Anything, mock.Anything).Return(true, nil)

	throttler.On("CanTrigger", mock.Anything, mock.Anything).Return(true, time.Duration(0), nil)

	dispatcher.Dispatch(models.IncomingEvent{
		Key:        "product.update",
		Source:     "core",
		Payload:    map[string]interface{}{"field1": 1, "field2": 2},
		TargetType: "Product",
		TargetId:   "product_id",
	}, store, throttler, documentStores)

	time.Sleep(50 * time.Millisecond)

	sender.AssertExpectations(t)
}

func TestDispatchEvents_With_Exclude_Payload_Attrs(t *testing.T) {
	store, documentStores, sender, dispatcher, throttler := setup("fixtures/config.exclude_payload_attrs.json")

	sender.On("SendEvent", mock.MatchedBy(func(e models.IncomingEvent) bool {
		return reflect.DeepEqual(e.Payload, map[string]interface{}{"field2": 2})
	}), mock.Anything).Return(nil)
	store.On("Get", mock.Anything).Return("", false, time.Duration(0), nil)
	store.On("Set", mock.Anything, mock.Anything, mock.Anything).Return(true, nil)

	throttler.On("CanTrigger", mock.Anything, mock.Anything).Return(true, time.Duration(0), nil)

	dispatcher.Dispatch(models.IncomingEvent{
		Key:        "product.update",
		Source:     "core",
		Payload:    map[string]interface{}{"field1": 1, "field2": 2},
		TargetType: "Product",
		TargetId:   "product_id",
	}, store, throttler, documentStores)

	time.Sleep(50 * time.Millisecond)

	sender.AssertExpectations(t)
}

func TestDispatchEvents_WithSpecificDocumentStore(t *testing.T) {
	store, documentStores, sender, dispatcher, throttler := setup("fixtures/config.specific_document_store.json")
	defaultDocumentStore := new(mocks.DocumentStoreMock)
	anotherDocumentStore := new(mocks.DocumentStoreMock)
	documentStores["specific"] = anotherDocumentStore
	documentStores["default"] = defaultDocumentStore

	sender.On("SendEvent", mock.MatchedBy(hasDocument), mock.Anything).Return(nil)
	// sender.On("SendEvent", mock.MatchedBy(hasNoDocument), mock.Anything).Return(nil)
	store.On("Get", mock.Anything).Return("", false, time.Duration(0), nil)
	store.On("Set", mock.Anything, mock.Anything, mock.Anything).Return(true, nil)
	defaultDocumentStore.On("GetDocument", mock.Anything).Return(map[string]interface{}{"foo": "bar"})
	anotherDocumentStore.On("GetDocument", mock.Anything).Return(map[string]interface{}{"foo": "bar"})

	throttler.On("CanTrigger", mock.Anything, mock.Anything).Return(true, time.Duration(0), nil)

	dispatcher.Dispatch(models.IncomingEvent{
		Key:        "product.update",
		Source:     "core",
		Payload:    map[string]interface{}{"field1": 1},
		TargetType: "Product",
		TargetId:   "product_id",
	}, store, throttler, documentStores)

	time.Sleep(50 * time.Millisecond)
	sender.AssertExpectations(t)
}

func TestDispatchEvents_Throttled_Without_TrailingEdge(t *testing.T) {
	store, documentStores, sender, dispatcher, throttler := setup("fixtures/config.throttle.disable_trailing.json")

	sender.On("SendEvent", mock.Anything, mock.Anything).Return(nil)
	store.On("Get", mock.Anything).Return("", false, time.Duration(0), nil)
	store.On("Set", mock.Anything, mock.Anything, mock.Anything).Return(true, nil)
	store.On("Remove", mock.Anything).Return(true, nil)

	throttler.On("CanTrigger", mock.Anything, mock.Anything).Return(true, 500*time.Millisecond, nil).Once()
	throttler.On("CanTrigger", mock.Anything, mock.Anything).Return(false, 500*time.Millisecond, nil).Twice()

	dispatcher.Dispatch(models.IncomingEvent{
		Key:        "product.update",
		Source:     "core",
		Payload:    map[string]interface{}{"field1": 1},
		TargetType: "Product",
		TargetId:   "product_id",
	}, store, throttler, documentStores)

	dispatcher.Dispatch(models.IncomingEvent{
		Key:        "product.update",
		Source:     "core",
		Payload:    map[string]interface{}{"field1": 1},
		TargetType: "Product",
		TargetId:   "product_id",
	}, store, throttler, documentStores)

	dispatcher.Dispatch(models.IncomingEvent{
		Key:        "product.update",
		Source:     "core",
		Payload:    map[string]interface{}{"field1": 1},
		TargetType: "Product",
		TargetId:   "product_id",
	}, store, throttler, documentStores)

	time.Sleep(1000 * time.Millisecond)
	sender.AssertNumberOfCalls(t, "SendEvent", 1)
}
