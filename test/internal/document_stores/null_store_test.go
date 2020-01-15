package document_stores_test

import (
  "testing"

  document_stores "github.com/shoplineapp/captin/internal/document_stores"
  models "github.com/shoplineapp/captin/models"
  "github.com/stretchr/testify/assert"

)

func TestGetDocument(t *testing.T) {
  ns := document_stores.NewNullDocumentStore()
  e := models.IncomingEvent{}
  assert.Equal(t, ns.GetDocument(e), map[string]interface{}{})
}
