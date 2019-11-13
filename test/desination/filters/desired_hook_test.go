package destination_filters_test

import (
	"github.com/stretchr/testify/assert"
	"testing"

	. "github.com/shoplineapp/captin/destinations/filters"
	helpers "github.com/shoplineapp/captin/internal/helpers"
	models "github.com/shoplineapp/captin/models"
)

func TestDesiredHookFilterRunValidate(t *testing.T) {
	event := models.IncomingEvent{
		Control: map[string]interface{}{
			"desired_hooks": []string{"desired"},
		},
	}
	assert.Equal(t, true, helpers.Tuples(DesiredHookFilter{}.Run(event, models.Destination{Config: models.Configuration{Name: "desired"}}))[0])
	assert.Equal(t, false, helpers.Tuples(DesiredHookFilter{}.Run(event, models.Destination{Config: models.Configuration{Name: "not_desired"}}))[0])

	// When event is unmarshal from JSON, hooks will be type []interface{}
	event = models.IncomingEvent{
		Control: map[string]interface{}{
			"desired_hooks": []interface{}{"desired"},
		},
	}
	assert.Equal(t, true, helpers.Tuples(DesiredHookFilter{}.Run(event, models.Destination{Config: models.Configuration{Name: "desired"}}))[0])
	assert.Equal(t, false, helpers.Tuples(DesiredHookFilter{}.Run(event, models.Destination{Config: models.Configuration{Name: "not_desired"}}))[0])

	// When event control contains multiple hooks
	event = models.IncomingEvent{
		Control: map[string]interface{}{
			"desired_hooks": []interface{}{
				"hook-1",
				"hook-2",
			},
		},
	}
	assert.Equal(t, true, helpers.Tuples(DesiredHookFilter{}.Run(event, models.Destination{Config: models.Configuration{Name: "hook-1"}}))[0])
	assert.Equal(t, true, helpers.Tuples(DesiredHookFilter{}.Run(event, models.Destination{Config: models.Configuration{Name: "hook-2"}}))[0])
	assert.Equal(t, false, helpers.Tuples(DesiredHookFilter{}.Run(event, models.Destination{Config: models.Configuration{Name: "not_desired"}}))[0])
}

func TestDesiredHookFilterApplicable(t *testing.T) {
	destination := models.Destination{Config: models.Configuration{}}

	event := models.IncomingEvent{Control: map[string]interface{}{
		"desired_hooks": []interface{}{"desired"},
	}}
	assert.Equal(t, true, DesiredHookFilter{}.Applicable(event, destination))

	event = models.IncomingEvent{}
	assert.Equal(t, false, DesiredHookFilter{}.Applicable(event, destination))
}
