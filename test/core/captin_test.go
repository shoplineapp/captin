package models_test

import (
	"github.com/stretchr/testify/assert"
	"testing"

	. "github.com/shoplineapp/captin/core"
	interfaces "github.com/shoplineapp/captin/interfaces"
	models "github.com/shoplineapp/captin/models"
	"github.com/stretchr/testify/mock"
)

// EventSenderMock - Mock EventSenderInterface
type EventSenderMock struct {
	interfaces.EventSenderInterface
	mock.Mock
}

func TestNewCaptin(t *testing.T) {
	// When initilizing captin
	// It has a default http sender
	configMapper := models.ConfigurationMapper{}
	captin := NewCaptin(configMapper)
	if captin.SenderMapping["http"] == nil {
		t.Errorf("Expected Captin to have a default http sender")
	}
}

func TestExecute(t *testing.T) {
	// When event is not given or is invalid
	var errors []error

	_, errors = Captin{}.Execute(models.IncomingEvent{})

	if assert.Error(t, errors[0], "invalid incoming event") {
		assert.IsType(t, errors[0], &ExecutionError{})
	}
}

func TestSetSenderMapping(t *testing.T) {
	captin := Captin{}
	mockSender := EventSenderMock{}
	senderMapping := map[string]interfaces.EventSenderInterface{
		"mock": mockSender,
	}
	captin.SetSenderMapping(senderMapping)
	assert.Equal(t, captin.SenderMapping["mock"], mockSender)
}
