package outgoing_test

import (
	"testing"

	interfaces "github.com/shoplineapp/captin/interfaces"
	. "github.com/shoplineapp/captin/internal/outgoing"
	models "github.com/shoplineapp/captin/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type FilterMock struct {
	mock.Mock
}

func (f *FilterMock) Run(e models.IncomingEvent, d models.Destination) (bool, error) {
	args := f.Called(e, d)
	return args.Bool(0), args.Error(1)
}

func (f *FilterMock) Applicable(e models.IncomingEvent, d models.Destination) bool {
	args := f.Called(e, d)
	return args.Bool(0)
}

type MiddlewareMock struct {
	mock.Mock
}

func (m *MiddlewareMock) Apply(e models.IncomingEvent, d []models.Destination) (models.IncomingEvent, []models.Destination) {
	m.Called(e, d)
	return e, d
}

func TestCustom_Sift(t *testing.T) {
	// Test if given filter is called by destinations
	event := models.IncomingEvent{}

	filter := new(FilterMock)
	filter.On("Applicable", mock.Anything, mock.Anything).Return(true)
	filter.On("Run", mock.Anything, mock.Anything).Return(true, nil)
	middleware := new(MiddlewareMock)
	middleware.On("Apply", mock.Anything, mock.Anything)
	destinations := []models.Destination{
		{Config: models.Configuration{Source: "service_1"}},
		{Config: models.Configuration{Source: "service_2"}},
	}
	sifted := Custom{}.Sift(event, destinations, []interfaces.DestinationFilter{filter}, []interfaces.DestinationMiddleware{middleware})
	filter.AssertNumberOfCalls(t, "Applicable", 2)
	filter.AssertNumberOfCalls(t, "Run", 2)
	middleware.AssertNumberOfCalls(t, "Apply", 1)
	assert.Equal(t, len(sifted), 2)

	// Test if filter run is skipped when filter is not applicable
	filter = new(FilterMock)
	filter.On("Applicable", mock.Anything, mock.Anything).Return(false)
	filter.On("Run", mock.Anything, mock.Anything).Return(true, nil)
	middleware = new(MiddlewareMock)
	middleware.On("Apply", mock.Anything, mock.Anything)
	destinations = []models.Destination{
		{Config: models.Configuration{Source: "service_1"}},
		{Config: models.Configuration{Source: "service_2"}},
	}
	sifted = Custom{}.Sift(event, destinations, []interfaces.DestinationFilter{filter}, []interfaces.DestinationMiddleware{middleware})
	filter.AssertNumberOfCalls(t, "Applicable", 2)
	filter.AssertNumberOfCalls(t, "Run", 0)
	middleware.AssertNumberOfCalls(t, "Apply", 1)
	assert.Equal(t, len(sifted), 2)

	// Test if destinations is filtered when filter run returns false
	filter = new(FilterMock)
	filter.On("Applicable", mock.Anything, mock.Anything).Return(true)
	filter.On("Run", mock.Anything, mock.Anything).Return(false, nil)
	middleware = new(MiddlewareMock)
	middleware.On("Apply", mock.Anything, mock.Anything)
	destinations = []models.Destination{
		{Config: models.Configuration{Source: "service_1"}},
	}
	sifted = Custom{}.Sift(event, destinations, []interfaces.DestinationFilter{filter}, []interfaces.DestinationMiddleware{middleware})
	filter.AssertNumberOfCalls(t, "Applicable", 1)
	filter.AssertNumberOfCalls(t, "Run", 1)
	middleware.AssertNumberOfCalls(t, "Apply", 1)
	assert.Equal(t, len(sifted), 0)
}
