package destination_filters_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	. "github.com/shoplineapp/captin/v2/destinations/filters"
	helpers "github.com/shoplineapp/captin/v2/internal/helpers"
	models "github.com/shoplineapp/captin/v2/models"
)

func TestValidateFilterRunValidate(t *testing.T) {
	payload := map[string]interface{}{
		"_id":  "xxxxxx",
		"type": "line",
	}
	event := models.IncomingEvent{Payload: payload}
	assert.Equal(t, true, helpers.Tuples(ValidateFilter{}.Run(context.Background(), event, models.Destination{Config: models.Configuration{Validate: "document.type == 'line'"}}))[0])
	assert.Equal(t, false, helpers.Tuples(ValidateFilter{}.Run(context.Background(), event, models.Destination{Config: models.Configuration{Validate: "document.type != 'line'"}}))[0])

	// Error validate check
	assert.Equal(t, false, helpers.Tuples(ValidateFilter{}.Run(context.Background(), event, models.Destination{Config: models.Configuration{Validate: "document.noneExist"}}))[0])
	assert.Equal(t, false, helpers.Tuples(ValidateFilter{}.Run(context.Background(), event, models.Destination{Config: models.Configuration{Validate: "invalidCall()"}}))[0])
	assert.Equal(t, false, helpers.Tuples(ValidateFilter{}.Run(context.Background(), event, models.Destination{Config: models.Configuration{Validate: "asd"}}))[0])

	// Calling with nil payload should not affect test result
	event = models.IncomingEvent{Payload: nil}
	assert.Equal(t, true, helpers.Tuples(ValidateFilter{}.Run(context.Background(), event, models.Destination{Config: models.Configuration{Validate: "true"}}))[0])
	assert.Equal(t, false, helpers.Tuples(ValidateFilter{}.Run(context.Background(), event, models.Destination{Config: models.Configuration{Validate: "false"}}))[0])
}

func TestValidateFilterApplicable(t *testing.T) {
	event := models.IncomingEvent{}
	assert.Equal(t, true, ValidateFilter{}.Applicable(context.Background(), event, models.Destination{Config: models.Configuration{Validate: "true"}}))
	assert.Equal(t, false, ValidateFilter{}.Applicable(context.Background(), event, models.Destination{Config: models.Configuration{}}))
}
