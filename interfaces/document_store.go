package interfaces

// StoreInterface - Store for throttle events
type DocumentStoreInterface interface {
  // GetDocument - Get value from store, return the document map
  GetDocument(e IncomingEventInterface) (map[string]interface{})
}
