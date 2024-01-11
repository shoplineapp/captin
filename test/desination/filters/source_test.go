package destination_filters_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	. "github.com/shoplineapp/captin/v2/destinations/filters"
	helpers "github.com/shoplineapp/captin/v2/internal/helpers"
	models "github.com/shoplineapp/captin/v2/models"
)

func TestSourceFilterRunValidate(t *testing.T) {
	event := models.IncomingEvent{Source: "service_1"}
	assert.Equal(t, true, helpers.Tuples(SourceFilter{}.Run(context.Background(), event, models.Destination{Config: models.Configuration{Source: "service_2"}}))[0])
	assert.Equal(t, false, helpers.Tuples(SourceFilter{}.Run(context.Background(), event, models.Destination{Config: models.Configuration{Source: "service_1"}}))[0])
}

func TestSourceFilterApplicable(t *testing.T) {
	event := models.IncomingEvent{}
	assert.Equal(t, true, SourceFilter{}.Applicable(context.Background(), event, models.Destination{Config: models.Configuration{AllowLoopback: false}}))
	assert.Equal(t, false, SourceFilter{}.Applicable(context.Background(), event, models.Destination{Config: models.Configuration{AllowLoopback: true}}))
}
