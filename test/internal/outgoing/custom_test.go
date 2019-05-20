package outgoing_test

import (
	"testing"

	interfaces "captin/interfaces"
	models "captin/internal/models"
	. "captin/internal/outgoing"
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

func TestCustom_Sift(t *testing.T) {
	// Test if given filter is called by destinations
	event := models.IncomingEvent{}

	filter := new(FilterMock)
	filter.On("Applicable", mock.Anything, mock.Anything).Return(true)
	filter.On("Run", mock.Anything, mock.Anything).Return(true, nil)
	destinations := []models.Destination{
		{Config: models.Configuration{Source: "service_1"}},
		{Config: models.Configuration{Source: "service_2"}},
	}
	sifted := Custom{}.Sift(event, []interfaces.CustomFilter{filter}, destinations)
	filter.AssertNumberOfCalls(t, "Applicable", 2)
	filter.AssertNumberOfCalls(t, "Run", 2)
	assert.Equal(t, len(sifted), 2)

	// Test if filter run is skipped when filter is not applicable
	filter = new(FilterMock)
	filter.On("Applicable", mock.Anything, mock.Anything).Return(false)
	filter.On("Run", mock.Anything, mock.Anything).Return(true, nil)
	destinations = []models.Destination{
		{Config: models.Configuration{Source: "service_1"}},
		{Config: models.Configuration{Source: "service_2"}},
	}
	sifted = Custom{}.Sift(event, []interfaces.CustomFilter{filter}, destinations)
	filter.AssertNumberOfCalls(t, "Applicable", 2)
	filter.AssertNumberOfCalls(t, "Run", 0)
	assert.Equal(t, len(sifted), 2)

	// Test if destinations is filtered when filter run returns false
	filter = new(FilterMock)
	filter.On("Applicable", mock.Anything, mock.Anything).Return(true)
	filter.On("Run", mock.Anything, mock.Anything).Return(false, nil)
	destinations = []models.Destination{
		{Config: models.Configuration{Source: "service_1"}},
	}
	sifted = Custom{}.Sift(event, []interfaces.CustomFilter{filter}, destinations)
	filter.AssertNumberOfCalls(t, "Applicable", 1)
	filter.AssertNumberOfCalls(t, "Run", 1)
	assert.Equal(t, len(sifted), 0)
}
