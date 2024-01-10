package models_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	. "github.com/shoplineapp/captin/v2/core"
	captin_errors "github.com/shoplineapp/captin/v2/errors"
	interfaces "github.com/shoplineapp/captin/v2/interfaces"
	models "github.com/shoplineapp/captin/v2/models"
	"github.com/stretchr/testify/mock"
)

// EventSenderMock - Mock EventSenderInterface
type EventSenderMock struct {
	interfaces.EventSenderInterface
	mock.Mock
}

// DocumentStoreMock - Mock DocumentStoreInterface
type DocumentStoreMock struct {
	interfaces.DocumentStoreInterface
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
	var errors []interfaces.ErrorInterface
	captin := Captin{}

	_, errors = captin.Execute(context.Background(), models.IncomingEvent{})

	if assert.Error(t, errors[0], "invalid incoming event") {
		assert.IsType(t, errors[0], &captin_errors.ExecutionError{})
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

func TestSetDocumentStoreMapping(t *testing.T) {
	captin := Captin{}
	mockStore := DocumentStoreMock{}
	storeMapping := map[string]interfaces.DocumentStoreInterface{
		"mock": mockStore,
	}
	captin.SetDocumentStoreMapping(storeMapping)
	assert.Equal(t, captin.DocumentStoreMapping["mock"], mockStore)
}
