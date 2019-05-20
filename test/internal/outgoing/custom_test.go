package outgoing_test

import (
	"testing"

	models "github.com/shoplineapp/captin/internal/models"
	. "github.com/shoplineapp/captin/internal/outgoing"
	outgoing_filters "github.com/shoplineapp/captin/internal/outgoing/filters"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type FilterMock struct {
	mock.Mock
}

func (f *FilterMock) Run(c models.Configuration) (bool, error) {
	args := f.Called(c)
	return args.Bool(0), args.Error(1)
}

func (f *FilterMock) Applicable(c models.Configuration) bool {
	args := f.Called(c)
	return args.Bool(0)
}

func TestCustom_Sift(t *testing.T) {
	// Test if given filter is called by destinations
	filter := new(FilterMock)
	filter.On("Applicable", mock.Anything).Return(true)
	filter.On("Run", mock.Anything).Return(true, nil)
	destinations := []models.Destination{
		{Config: models.Configuration{Source: "service_1"}},
		{Config: models.Configuration{Source: "service_2"}},
	}
	sifted := Custom{}.Sift([]outgoing_filters.Filter{filter}, destinations)
	filter.AssertNumberOfCalls(t, "Applicable", 2)
	filter.AssertNumberOfCalls(t, "Run", 2)
	assert.Equal(t, len(sifted), 2)

	// Test if filter run is skipped when filter is not applicable
	filter = new(FilterMock)
	filter.On("Applicable", mock.Anything).Return(false)
	filter.On("Run", mock.Anything).Return(true, nil)
	destinations = []models.Destination{
		{Config: models.Configuration{Source: "service_1"}},
		{Config: models.Configuration{Source: "service_2"}},
	}
	sifted = Custom{}.Sift([]outgoing_filters.Filter{filter}, destinations)
	filter.AssertNumberOfCalls(t, "Applicable", 2)
	filter.AssertNumberOfCalls(t, "Run", 0)
	assert.Equal(t, len(sifted), 2)

	// Test if destinations is filtered when filter run returns false
	filter = new(FilterMock)
	filter.On("Applicable", mock.Anything).Return(true)
	filter.On("Run", mock.Anything).Return(false, nil)
	destinations = []models.Destination{
		{Config: models.Configuration{Source: "service_1"}},
	}
	sifted = Custom{}.Sift([]outgoing_filters.Filter{filter}, destinations)
	filter.AssertNumberOfCalls(t, "Applicable", 1)
	filter.AssertNumberOfCalls(t, "Run", 1)
	assert.Equal(t, len(sifted), 0)
}

func TestCustom_CustomFilters(t *testing.T) {
	filters := CustomFilters(models.IncomingEvent{})
	for _, filter := range filters {
		// Test if filter implements base filter interface
		_, ok := filter.(outgoing_filters.Filter)
		assert.Equal(t, ok, true)
	}
}
