package outgoing_filters_test

import (
	"github.com/stretchr/testify/assert"
	"testing"

	helpers "captin/internal/helpers"
	models "captin/internal/models"
	. "captin/internal/outgoing/filters"
)

func TestValidateFilterRunValidate(t *testing.T) {
	payload := map[string]interface{}{
		"_id":  "xxxxxx",
		"type": "line",
	}
	event := models.IncomingEvent{Payload: payload}
	assert.Equal(t, true, helpers.Tuples(ValidateFilter{}.Run(event, models.Destination{Config: models.Configuration{Validate: "document.type == 'line'"}}))[0])
	assert.Equal(t, false, helpers.Tuples(ValidateFilter{}.Run(event, models.Destination{Config: models.Configuration{Validate: "document.type != 'line'"}}))[0])

	// Error validate check
	assert.Equal(t, false, helpers.Tuples(ValidateFilter{}.Run(event, models.Destination{Config: models.Configuration{Validate: "document.noneExist"}}))[0])
	assert.Equal(t, false, helpers.Tuples(ValidateFilter{}.Run(event, models.Destination{Config: models.Configuration{Validate: "invalidCall()"}}))[0])
	assert.Equal(t, false, helpers.Tuples(ValidateFilter{}.Run(event, models.Destination{Config: models.Configuration{Validate: "asd"}}))[0])

	// Calling with nil payload should not affect test result
	event = models.IncomingEvent{Payload: nil}
	assert.Equal(t, true, helpers.Tuples(ValidateFilter{}.Run(event, models.Destination{Config: models.Configuration{Validate: "true"}}))[0])
	assert.Equal(t, false, helpers.Tuples(ValidateFilter{}.Run(event, models.Destination{Config: models.Configuration{Validate: "false"}}))[0])
}

func TestValidateFilterApplicable(t *testing.T) {
	event := models.IncomingEvent{}
	assert.Equal(t, true, ValidateFilter{}.Applicable(event, models.Destination{Config: models.Configuration{Validate: "true"}}))
	assert.Equal(t, false, ValidateFilter{}.Applicable(event, models.Destination{Config: models.Configuration{}}))
}
