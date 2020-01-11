package interfaces

import (
  models "github.com/shoplineapp/captin/models"
)

// StoreInterface - Store for throttle events
type DocumentStoreInterface interface {
  // GetDocument - Get value from store, return the document map
  GetDocument(e models.IncomingEvent) (map[string]interface{})
}
