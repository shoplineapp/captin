package incoming_tests

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"testing"

	internal "github.com/shoplineapp/captin/internal"
	incoming "github.com/shoplineapp/captin/internal/incoming"
	models "github.com/shoplineapp/captin/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type CaptinMock struct {
	internal.Captin
	mock.Mock
}

func (f *CaptinMock) Execute(c models.IncomingEvent) (bool, error) {
	args := f.Called(c)
	return args.Bool(0), args.Error(1)
}

func TestHttpEventHandler_SetRoutes(t *testing.T) {
	router := gin.Default()

	configMapper := models.NewConfigurationMapper([]models.Configuration{})

	c := internal.Captin{ConfigMap: *configMapper}

	handler := incoming.HttpEventHandler{}
	handler.Setup(&c)
	handler.SetRoutes(router)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/non-exist", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, 404, w.Code)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	jsonStr := []byte(`{"event_key":"model.action","source":"service_one","payload":{"_id":"xxxxx"}}`)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/events", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	assert.Equal(t, 201, w.Code)
}
