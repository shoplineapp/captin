package incoming_tests

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/shoplineapp/captin/v2/incoming"
	interfaces "github.com/shoplineapp/captin/v2/interfaces"
	models "github.com/shoplineapp/captin/v2/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type captinMock struct {
	interfaces.CaptinInterface
	mock.Mock
}

func (f *captinMock) Execute(c models.IncomingEvent) (bool, []error) {
	args := f.Called(c)
	errors := args.Error(1)
	if errors == nil {
		errors = []error{}
	}
	return args.Bool(0), errors
}

func TestHttpEventHandler_SetRoutes(t *testing.T) {
	router := gin.Default()

	captin := new(captinMock)
	captin.On("Execute", mock.Anything).Return(true, []errors{})

	handler := HttpEventHandler{}
	handler.Setup(captin)
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

	captin.AssertNumberOfCalls(t, "Execute", 1)
}
