package outgoing_test

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	outgoing "github.com/shoplineapp/captin/internal/outgoing"
	models "github.com/shoplineapp/captin/models"
	"github.com/stretchr/testify/mock"
)

type senderMock struct {
	mock.Mock
}

func (s *senderMock) SendEvent(e models.IncomingEvent, d models.Destination) error {
	args := s.Called(e, d)
	return args.Error(0)
}

func TestDispatchEvents_Error(t *testing.T) {
	data, err := ioutil.ReadFile("fixtures/config.json")
	if err != nil {
		panic(err)
	}
	configs := []models.Configuration{}
	json.Unmarshal(data, &configs)
	destinations := []models.Destination{}
	for _, config := range configs {
		destinations = append(destinations, models.Destination{Config: config})
	}

	sender := new(senderMock)

	sender.On("SendEvent", mock.Anything, mock.Anything).Return(errors.New("Mock Error"))

	dispatcher := outgoing.NewDispatcherWithDestinations(destinations, sender)

	dispatcher.Dispatch(models.IncomingEvent{
		Key:        "product.update",
		Source:     "core",
		Payload:    map[string]interface{}{"field1": 1},
		TargetType: "Product",
		TargetId:   "product_id",
	})

	dispatcher.Dispatch(models.IncomingEvent{
		Key:        "product.update",
		Source:     "core",
		Payload:    map[string]interface{}{"field1": 2},
		TargetType: "Product",
		TargetId:   "product_id_2",
	})

	time.Sleep(50 * time.Millisecond)

	sender.AssertNumberOfCalls(t, "SendEvent", 6)
	assert.Equal(t, 6, len(dispatcher.Errors))
}

func TestDispatchEvents(t *testing.T) {
	data, err := ioutil.ReadFile("fixtures/config.json")
	if err != nil {
		panic(err)
	}
	configs := []models.Configuration{}
	json.Unmarshal(data, &configs)
	destinations := []models.Destination{}
	for _, config := range configs {
		destinations = append(destinations, models.Destination{Config: config})
	}

	sender := new(senderMock)

	sender.On("SendEvent", mock.Anything, mock.Anything).Return(nil)

	dispatcher := outgoing.NewDispatcherWithDestinations(destinations, sender)

	dispatcher.Dispatch(models.IncomingEvent{
		Key:        "product.update",
		Source:     "core",
		Payload:    map[string]interface{}{"field1": 1},
		TargetType: "Product",
		TargetId:   "product_id",
	})

	dispatcher.Dispatch(models.IncomingEvent{
		Key:        "product.update",
		Source:     "core",
		Payload:    map[string]interface{}{"field1": 2},
		TargetType: "Product",
		TargetId:   "product_id_2",
	})

	time.Sleep(50 * time.Millisecond)

	sender.AssertNumberOfCalls(t, "SendEvent", 6)
}
