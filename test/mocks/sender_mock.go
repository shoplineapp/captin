package mocks

import (
	"context"

	"github.com/shoplineapp/captin/v2/interfaces"
	"github.com/shoplineapp/captin/v2/models"
	"github.com/stretchr/testify/mock"
)

// SenderMock - Mock of SenderInterface
type SenderMock struct {
	mock.Mock
	interfaces.EventSenderInterface
}

// SendEvent - Send an event
func (s *SenderMock) SendEvent(ctx context.Context, ie interfaces.IncomingEventInterface, id interfaces.DestinationInterface) error {
	e := ie.(models.IncomingEvent)
	d := id.(models.Destination)
	args := s.Called(ctx, e, d)
	return args.Error(0)
}
