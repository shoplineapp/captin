package mocks

import (
  "github.com/shoplineapp/captin/interfaces"
  "github.com/shoplineapp/captin/models"
  "github.com/stretchr/testify/mock"
)

// DocumentStoreMock - Mock of DocumentStoreInterface
type DocumentStoreMock struct {
  interfaces.DocumentStoreInterface
  mock.Mock
}

// Get - Get value from store, return with remaining time
func (ds *DocumentStoreMock) GetDocument(e models.IncomingEvent) (map[string]interface{}) {
  args := ds.Called(e)
  return args.Get(0).(map[string]interface{})
}

