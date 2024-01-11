package mocks

import (
	"context"

	"github.com/shoplineapp/captin/v2/interfaces"
	"github.com/shoplineapp/captin/v2/models"
	"github.com/stretchr/testify/mock"
)

// DocumentStoreMock - Mock of DocumentStoreInterface
type DocumentStoreMock struct {
	interfaces.DocumentStoreInterface
	mock.Mock
}

// Get - Get value from store, return with remaining time
func (ds *DocumentStoreMock) GetDocument(ctx context.Context, ie interfaces.IncomingEventInterface) map[string]interface{} {
	e := ie.(models.IncomingEvent)
	args := ds.Called(ctx, e)
	return args.Get(0).(map[string]interface{})
}
