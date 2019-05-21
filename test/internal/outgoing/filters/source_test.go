package outgoing_filters_test

import (
	"github.com/stretchr/testify/assert"
	"testing"

	helpers "github.com/shoplineapp/captin/internal/helpers"
	models "github.com/shoplineapp/captin/internal/models"
	. "github.com/shoplineapp/captin/internal/outgoing/filters"
)

func TestSourceFilterRunValidate(t *testing.T) {
	event := models.IncomingEvent{Source: "service_1"}
	assert.Equal(t, true, helpers.Tuples(SourceFilter{}.Run(event, models.Destination{Config: models.Configuration{Source: "service_2"}}))[0])
	assert.Equal(t, false, helpers.Tuples(SourceFilter{}.Run(event, models.Destination{Config: models.Configuration{Source: "service_1"}}))[0])
}

func TestSourceFilterApplicable(t *testing.T) {
	event := models.IncomingEvent{}
	assert.Equal(t, true, SourceFilter{}.Applicable(event, models.Destination{Config: models.Configuration{AllowLoopback: false}}))
	assert.Equal(t, false, SourceFilter{}.Applicable(event, models.Destination{Config: models.Configuration{AllowLoopback: true}}))
}
