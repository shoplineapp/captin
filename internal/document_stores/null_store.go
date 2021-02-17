package document_stores

import (
  interfaces "github.com/shoplineapp/captin/interfaces"
)

// NullDocumentStore - Null data store
type NullDocumentStore struct {
  interfaces.DocumentStoreInterface
}

// NewNullDocumentStore - Create new NullDocumentStore
func NewNullDocumentStore() *NullDocumentStore {
  return &NullDocumentStore{}
}

// Get - Get value from store, return with remaining time
func (ms NullDocumentStore) GetDocument(e interfaces.IncomingEventInterface) (map[string]interface{}) {
  return map[string]interface{}{}
}

