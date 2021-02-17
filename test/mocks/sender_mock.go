package mocks

import (
	"github.com/shoplineapp/captin/interfaces"
	"github.com/shoplineapp/captin/models"
	"github.com/stretchr/testify/mock"
)

// SenderMock - Mock of SenderInterface
type SenderMock struct {
	mock.Mock
	interfaces.EventSenderInterface
}

// SendEvent - Send an event
func (s *SenderMock) SendEvent(ie interfaces.IncomingEventInterface, id interfaces.DestinationInterface) error {
	e := ie.(models.IncomingEvent)
	d := id.(models.Destination)
	args := s.Called(e, d)
	return args.Error(0)
}
