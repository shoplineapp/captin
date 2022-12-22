package outgoing_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"reflect"
	"testing"
	"time"
	"unsafe"

	delayers "github.com/shoplineapp/captin/dispatcher/delayers"
	captin_errors "github.com/shoplineapp/captin/errors"
	interfaces "github.com/shoplineapp/captin/interfaces"
	outgoing "github.com/shoplineapp/captin/internal/outgoing"
	stores "github.com/shoplineapp/captin/internal/stores"
	models "github.com/shoplineapp/captin/models"
	mocks "github.com/shoplineapp/captin/test/mocks"
	"github.com/stretchr/testify/assert"
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

	sender.AssertNumberOfCalls(t, "SendEvent", 6)
	assert.Equal(t, 6, len(dispatcher.GetErrors()))
}

func TestDispatchEvents_SendEvent_WithNotDispatcherError(t *testing.T) {
	store, documentStores, sender, dispatcher, throttler := setup("fixtures/config.json")

	sender.On("SendEvent", mock.Anything, mock.Anything).Panic("Non-dispatcher panic")
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

	sender.AssertNumberOfCalls(t, "SendEvent", 6)
	assert.Equal(t, 6, len(dispatcher.GetErrors()))
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

	sender.AssertNumberOfCalls(t, "SendEvent", 6)
}

func TestDispatchEvents_Throttled_DelaySend(t *testing.T) {
	store, documentStores, sender, dispatcher, throttler := setup("fixtures/config.single.json")

	sender.On("SendEvent", mock.Anything, mock.Anything).Return(nil)
	store.On("Get", mock.Anything).Return("", false, time.Duration(0), nil).Once()
	store.On("Get", mock.Anything).Return("", true, time.Duration(0), nil)
	store.On("Set", mock.Anything, mock.Anything, mock.Anything).Return(true, nil)
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
	store.AssertCalled(t, "Set", throttleID, `{"IncomingEventInterface":null,"TraceId":"","event_key":"product.update","source":"core","payload":{"field1":1},"control":null,"target_type":"Product","target_id":"product_id"}`, throttlePeriod)

	time.Sleep(600 * time.Millisecond)

	sender.AssertNumberOfCalls(t, "SendEvent", 1)
}

func TestDispatchEvents_Delayer_Send(t *testing.T) {
	store, documentStores, sender, dispatcher, throttler := setup("fixtures/config.delay.json")
	dispatcher.SetDelayer(&delayers.GoroutineDelayer{})
	sender.On("SendEvent", mock.Anything, mock.Anything).Return(nil)
	throttler.On("CanTrigger", mock.Anything, mock.Anything).Return(false, 1*time.Millisecond, nil)
	store.On("Get", mock.Anything).Return("", false, time.Duration(0), nil).Twice()
	store.On("Get", mock.Anything).Return("", true, time.Duration(0), nil)
	store.On("Set", mock.Anything, mock.Anything, mock.Anything).Return(true, nil)

	dispatcher.Dispatch(models.IncomingEvent{
		Key:        "product.update",
		Source:     "core",
		Payload:    map[string]interface{}{"field1": 1},
		TargetType: "Product",
		TargetId:   "product_id",
		Control:    map[string]interface{}{"existing_data": "test"},
	}, store, throttler, documentStores)

	time.Sleep(100 * time.Millisecond)

	sender.AssertNumberOfCalls(t, "SendEvent", 2)
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
	store.AssertCalled(t, "Update", throttleID, `{"IncomingEventInterface":null,"TraceId":"","event_key":"product.update","source":"core","payload":{"field1":2},"control":null,"target_type":"Product","target_id":"product_id"}`)
}

func TestDispatchEvents_Throttled_KeepThrottledPayloads(t *testing.T) {
	_, documentStores, sender, dispatcher, throttler := setup("fixtures/config.keep_throttled_payloads.json")
	store := stores.NewMemoryStore()

	throttler.On("CanTrigger", mock.Anything, mock.Anything).Return(false, 500*time.Millisecond, nil)
	sender.On("SendEvent", mock.Anything, mock.Anything).Return(nil)

	throttleID := "product.update.service_one.product_id-data"
	throttlePayloadsID := "product.update.service_one.product_id-throttled_payloads"

	dispatcher.Dispatch(models.IncomingEvent{
		Key:        "product.update",
		Source:     "core",
		Payload:    map[string]interface{}{"field1": 1},
		TargetType: "Product",
		TargetId:   "product_id",
	}, store, throttler, documentStores)

	throttledPayloads, _, _, _ := store.GetQueue(throttlePayloadsID)
	assert.Equal(t, len(throttledPayloads), 1)

	sender.AssertNumberOfCalls(t, "SendEvent", 0)

	time.Sleep(600 * time.Millisecond)

	sender.AssertCalled(t, "SendEvent", mock.MatchedBy(func(e models.IncomingEvent) bool {
		return fmt.Sprint(e.ThrottledPayloads) == fmt.Sprint([]map[string]interface{}{
			{"field1": 1},
		})
	}), mock.Anything)

	_, storedEventExists, _, _ := store.Get(throttleID)
	throttledPayloads, _, _, _ = store.GetQueue(throttlePayloadsID)
	assert.Equal(t, true, storedEventExists)
	assert.Equal(t, 0, len(throttledPayloads))
}

func TestDispatchEvents_Throttled_KeepThrottledPayloads_Multiple(t *testing.T) {
	_, documentStores, sender, dispatcher, throttler := setup("fixtures/config.keep_throttled_payloads.json")
	store := stores.NewMemoryStore()

	throttler.On("CanTrigger", mock.Anything, mock.Anything).Return(false, 500*time.Millisecond, nil)
	sender.On("SendEvent", mock.Anything, mock.Anything).Return(nil)

	throttleID := "product.update.service_one.product_id-data"
	throttlePayloadsID := "product.update.service_one.product_id-throttled_payloads"

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
		TargetId:   "product_id",
	}, store, throttler, documentStores)

	throttledPayloads, _, _, _ := store.GetQueue(throttlePayloadsID)
	assert.Equal(t, 2, len(throttledPayloads))

	sender.AssertNumberOfCalls(t, "SendEvent", 0)

	time.Sleep(600 * time.Millisecond)

	// it should SendEvent with 2 throttled payloads
	sender.AssertCalled(t, "SendEvent", mock.MatchedBy(func(e models.IncomingEvent) bool {
		return fmt.Sprint(e.ThrottledPayloads) == fmt.Sprint([]map[string]interface{}{
			{"field1": 1},
			{"field1": 2},
		})
	}), mock.Anything)

	// it should SendEvent with last payload
	sender.AssertCalled(t, "SendEvent", mock.MatchedBy(func(e models.IncomingEvent) bool {
		return fmt.Sprint(e.Payload) == fmt.Sprint(map[string]interface{}{"field1": 2})
	}), mock.Anything)

	// it should clean up after sendEvent
	_, storedEventExists, _, _ := store.Get(throttleID)
	throttledPayloads, _, _, _ = store.GetQueue(throttlePayloadsID)
	assert.Equal(t, true, storedEventExists)
	assert.Equal(t, 0, len(throttledPayloads))
}

func TestDispatchEvents_Throttled_KeepThrottledDocuments(t *testing.T) {
	_, documentStores, sender, dispatcher, throttler := setup("fixtures/config.keep_throttled_documents.json")
	store := stores.NewMemoryStore()

	mockDocumentStore := new(mocks.DocumentStoreMock)
	documentStores["default"] = mockDocumentStore
	mockDocumentStore.On("GetDocument", mock.Anything).Return(map[string]interface{}{"foo": "bar"})

	throttler.On("CanTrigger", mock.Anything, mock.Anything).Return(false, 500*time.Millisecond, nil)
	sender.On("SendEvent", mock.Anything, mock.Anything).Return(nil)

	throttleID := "product.update.service_one.product_id-data"
	throttleDocumentsID := "product.update.service_one.product_id-throttled_documents"

	dispatcher.Dispatch(models.IncomingEvent{
		Key:        "product.update",
		Source:     "core",
		TargetType: "Product",
		TargetId:   "product_id",
	}, store, throttler, documentStores)

	throttledDocuments, _, _, _ := store.GetQueue(throttleDocumentsID)
	assert.Equal(t, 1, len(throttledDocuments))

	sender.AssertNumberOfCalls(t, "SendEvent", 0)

	time.Sleep(600 * time.Millisecond)

	sender.AssertCalled(t, "SendEvent", mock.MatchedBy(func(e models.IncomingEvent) bool {
		return fmt.Sprint(e.ThrottledDocuments) == fmt.Sprint([]map[string]interface{}{
			{"foo": "bar"},
		})
	}), mock.Anything)

	// it should clean up after sendEvent
	_, storedEventExists, _, _ := store.Get(throttleID)
	throttledDocuments, _, _, _ = store.GetQueue(throttleDocumentsID)
	assert.Equal(t, true, storedEventExists)
	assert.Equal(t, 0, len(throttledDocuments))
}

func TestDispatchEvents_Throttled_KeepThrottledDocuments_Multiple(t *testing.T) {
	_, documentStores, sender, dispatcher, throttler := setup("fixtures/config.keep_throttled_documents.json")
	store := stores.NewMemoryStore()

	mockDocumentStore := new(mocks.DocumentStoreMock)
	documentStores["default"] = mockDocumentStore
	mockDocumentStore.On("GetDocument", mock.Anything).Return(map[string]interface{}{"foo": "bar1"}).Once()
	mockDocumentStore.On("GetDocument", mock.Anything).Return(map[string]interface{}{"foo": "bar2"})

	throttler.On("CanTrigger", mock.Anything, mock.Anything).Return(false, 500*time.Millisecond, nil)
	sender.On("SendEvent", mock.Anything, mock.Anything).Return(nil)

	throttleID := "product.update.service_one.product_id-data"
	throttleDocumentsID := "product.update.service_one.product_id-throttled_documents"

	dispatcher.Dispatch(models.IncomingEvent{
		Key:        "product.update",
		Source:     "core",
		TargetType: "Product",
		TargetId:   "product_id",
	}, store, throttler, documentStores)

	// clear cache in dispatcher.targetDocument private field
	// src: https://gist.github.com/CyberLight/1da35b4e0093bc12302f
	ptrToTargetDocument := (*interface{})(unsafe.Pointer(reflect.Indirect(reflect.ValueOf(dispatcher)).FieldByName("targetDocument").UnsafeAddr()))
	*ptrToTargetDocument = nil

	dispatcher.Dispatch(models.IncomingEvent{
		Key:        "product.update",
		Source:     "core",
		TargetType: "Product",
		TargetId:   "product_id",
	}, store, throttler, documentStores)

	throttledDocuments, _, _, _ := store.GetQueue(throttleDocumentsID)
	assert.Equal(t, 2, len(throttledDocuments))

	sender.AssertNumberOfCalls(t, "SendEvent", 0)

	time.Sleep(600 * time.Millisecond)

	sender.AssertCalled(t, "SendEvent", mock.MatchedBy(func(e models.IncomingEvent) bool {
		return fmt.Sprint(e.ThrottledDocuments) == fmt.Sprint([]map[string]interface{}{
			{"foo": "bar1"},
			{"foo": "bar2"},
		})
	}), mock.Anything)

	// it should clean up after sendEvent
	_, storedEventExists, _, _ := store.Get(throttleID)
	throttledDocuments, _, _, _ = store.GetQueue(throttleDocumentsID)
	assert.Equal(t, true, storedEventExists)
	assert.Equal(t, 0, len(throttledDocuments))
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

func TestDispatchEvents_OnError(t *testing.T) {

	_, _, _, dispatcher, _ := setup("fixtures/config.single.json")

	dispatcher.OnError(
		models.IncomingEvent{
			Key:        "product.update",
			Source:     "core",
			Payload:    map[string]interface{}{"field1": 1},
			TargetType: "Product",
			TargetId:   "product_id",
		}, errors.New("error"),
	)

	assert.Equal(t, 1, len(dispatcher.Errors))
}

func TestDispatchEvents_DispatchErrorTriggerOnError(t *testing.T) {

	store, documentStores, sender, dispatcher, throttler := setup("fixtures/config.single.json")

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

	time.Sleep(50 * time.Millisecond)

	sender.AssertExpectations(t)
	// error is only append to errors list when OnError is called
	assert.Equal(t, 1, len(dispatcher.GetErrors()))
}

func TestDispatchEvents_DelaySend_NotExist(t *testing.T) {
	store, documentStores, sender, dispatcher, throttler := setup("fixtures/config.single.json")

	sender.On("SendEvent", mock.Anything, mock.Anything).Return(nil)
	store.On("Get", mock.Anything).Return("", false, time.Duration(0), nil)
	store.On("Set", mock.Anything, mock.Anything, mock.Anything).Return(true, nil)
	throttler.On("CanTrigger", mock.Anything, mock.Anything).Return(false, time.Duration(0), nil)

	dispatcher.Dispatch(models.IncomingEvent{
		Key:        "product.update",
		Source:     "core-api",
		Payload:    map[string]interface{}{"field1": 1},
		TargetType: "Product",
		TargetId:   "product_id",
	}, store, throttler, documentStores)

	time.Sleep(50 * time.Millisecond)

	assert.Equal(t, 1, len(dispatcher.GetErrors()))
	assert.IsType(t, &captin_errors.UnretryableError{}, dispatcher.GetErrors()[0])
	sender.AssertNumberOfCalls(t, "SendEvent", 0)
}
