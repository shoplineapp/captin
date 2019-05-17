package outgoing_filters_test

import (
	"github.com/stretchr/testify/assert"
	"testing"

	helpers "github.com/shoplineapp/captin/internal/helpers"
	models "github.com/shoplineapp/captin/internal/models"
	. "github.com/shoplineapp/captin/internal/outgoing/filters"
)

func TestValidateFilterRunValidate(t *testing.T) {
	payload := map[string]interface{}{
		"_id":  "xxxxxx",
		"type": "line",
	}
	assert.Equal(t, true, helpers.Tuples(ValidateFilter{Event: models.IncomingEvent{Payload: payload}}.Run(models.Configuration{Validate: "document.type == 'line'"}))[0])
	assert.Equal(t, false, helpers.Tuples(ValidateFilter{Event: models.IncomingEvent{Payload: payload}}.Run(models.Configuration{Validate: "document.type != 'line'"}))[0])

	// Error validate check
	assert.Equal(t, false, helpers.Tuples(ValidateFilter{Event: models.IncomingEvent{Payload: payload}}.Run(models.Configuration{Validate: "document.noneExist"}))[0])
	assert.Equal(t, false, helpers.Tuples(ValidateFilter{Event: models.IncomingEvent{Payload: payload}}.Run(models.Configuration{Validate: "invalidCall()"}))[0])
	assert.Equal(t, false, helpers.Tuples(ValidateFilter{Event: models.IncomingEvent{Payload: payload}}.Run(models.Configuration{Validate: "asd"}))[0])

	// Calling with nil payload should not affect test result
	assert.Equal(t, true, helpers.Tuples(ValidateFilter{Event: models.IncomingEvent{Payload: nil}}.Run(models.Configuration{Validate: "true"}))[0])
	assert.Equal(t, false, helpers.Tuples(ValidateFilter{Event: models.IncomingEvent{Payload: nil}}.Run(models.Configuration{Validate: "false"}))[0])
}

func TestValidateFilterApplicable(t *testing.T) {
	assert.Equal(t, true, ValidateFilter{}.Applicable(models.Configuration{Validate: "true"}))
	assert.Equal(t, false, ValidateFilter{}.Applicable(models.Configuration{}))
}
