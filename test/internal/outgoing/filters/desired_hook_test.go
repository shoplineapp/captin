package outgoing_filters_test

import (
	"github.com/stretchr/testify/assert"
	"testing"

	helpers "github.com/shoplineapp/captin/internal/helpers"
	. "github.com/shoplineapp/captin/internal/outgoing/filters"
	models "github.com/shoplineapp/captin/models"
)

func TestDesiredHookFilterRunValidate(t *testing.T) {
	event := models.IncomingEvent{Control: map[string]interface{}{
		"desired_hooks": []string{"desired"},
	}}
	assert.Equal(t, true, helpers.Tuples(DesiredHookFilter{}.Run(event, models.Destination{Config: models.Configuration{Name: "desired"}}))[0])
	assert.Equal(t, false, helpers.Tuples(DesiredHookFilter{}.Run(event, models.Destination{Config: models.Configuration{Name: "not_desired"}}))[0])
}

func TestDesiredHookFilterApplicable(t *testing.T) {
	destination := models.Destination{Config: models.Configuration{}}

	event := models.IncomingEvent{Control: map[string]interface{}{
		"desired_hooks": []string{"desired"},
	}}
	assert.Equal(t, true, DesiredHookFilter{}.Applicable(event, destination))

	event = models.IncomingEvent{}
	assert.Equal(t, false, DesiredHookFilter{}.Applicable(event, destination))
}
