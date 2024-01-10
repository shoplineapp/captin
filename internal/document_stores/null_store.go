package document_stores

import (
	"context"

	interfaces "github.com/shoplineapp/captin/v2/interfaces"
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
func (ms NullDocumentStore) GetDocument(ctx context.Context, e interfaces.IncomingEventInterface) map[string]interface{} {
	return map[string]interface{}{}
}
