package interfaces

import "context"

// StoreInterface - Store for throttle events
type DocumentStoreInterface interface {
	// GetDocument - Get value from store, return the document map
	GetDocument(ctx context.Context, e IncomingEventInterface) map[string]interface{}
}
