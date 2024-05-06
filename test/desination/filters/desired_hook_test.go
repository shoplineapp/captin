package destination_filters_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	. "github.com/shoplineapp/captin/v2/destinations/filters"
	helpers "github.com/shoplineapp/captin/v2/internal/helpers"
	models "github.com/shoplineapp/captin/v2/models"
)

func TestDesiredHookFilterRunValidate(t *testing.T) {
	event := models.IncomingEvent{
		Control: map[string]interface{}{
			"desired_hooks": []string{"desired"},
		},
	}
	assert.Equal(t, true, helpers.Tuples(DesiredHookFilter{}.Run(context.Background(), event, models.Destination{Config: models.Configuration{Name: "desired"}}))[0])
	assert.Equal(t, false, helpers.Tuples(DesiredHookFilter{}.Run(context.Background(), event, models.Destination{Config: models.Configuration{Name: "not_desired"}}))[0])

	// When event is unmarshal from JSON, hooks will be type []interface{}
	event = models.IncomingEvent{
		Control: map[string]interface{}{
			"desired_hooks": []interface{}{"desired"},
		},
	}
	assert.Equal(t, true, helpers.Tuples(DesiredHookFilter{}.Run(context.Background(), event, models.Destination{Config: models.Configuration{Name: "desired"}}))[0])
	assert.Equal(t, false, helpers.Tuples(DesiredHookFilter{}.Run(context.Background(), event, models.Destination{Config: models.Configuration{Name: "not_desired"}}))[0])

	// When event control contains multiple hooks
	event = models.IncomingEvent{
		Control: map[string]interface{}{
			"desired_hooks": []interface{}{
				"hook-1",
				"hook-2",
			},
		},
	}
	assert.Equal(t, true, helpers.Tuples(DesiredHookFilter{}.Run(context.Background(), event, models.Destination{Config: models.Configuration{Name: "hook-1"}}))[0])
	assert.Equal(t, true, helpers.Tuples(DesiredHookFilter{}.Run(context.Background(), event, models.Destination{Config: models.Configuration{Name: "hook-2"}}))[0])
	assert.Equal(t, false, helpers.Tuples(DesiredHookFilter{}.Run(context.Background(), event, models.Destination{Config: models.Configuration{Name: "not_desired"}}))[0])
}

func TestDesiredHookFilterApplicable(t *testing.T) {
	destination := models.Destination{Config: models.Configuration{}}

	event := models.IncomingEvent{Control: map[string]interface{}{
		"desired_hooks": []interface{}{"desired"},
	}}
	assert.Equal(t, true, DesiredHookFilter{}.Applicable(context.Background(), event, destination))

	event = models.IncomingEvent{}
	assert.Equal(t, false, DesiredHookFilter{}.Applicable(context.Background(), event, destination))
}
