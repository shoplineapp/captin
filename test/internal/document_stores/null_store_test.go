package document_stores_test

import (
	"context"
	"testing"

	document_stores "github.com/shoplineapp/captin/v2/internal/document_stores"
	models "github.com/shoplineapp/captin/v2/models"
	"github.com/stretchr/testify/assert"
)

func TestGetDocument(t *testing.T) {
	ns := document_stores.NewNullDocumentStore()
	e := models.IncomingEvent{}
	assert.Equal(t, ns.GetDocument(context.Background(), e), map[string]interface{}{})
}
