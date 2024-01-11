package destination_filters_test

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	. "github.com/shoplineapp/captin/v2/destinations/filters"
	helpers "github.com/shoplineapp/captin/v2/internal/helpers"
	models "github.com/shoplineapp/captin/v2/models"
)

func TestEnvironmentFilterRunValidate(t *testing.T) {
	defer os.Unsetenv("HOOK_SERVICE_1_ENABLED")
	event := models.IncomingEvent{}
	assert.Equal(t, true, helpers.Tuples(EnvironmentFilter{}.Run(context.Background(), event, models.Destination{Config: models.Configuration{Name: "service_1"}}))[0])
	os.Setenv("HOOK_SERVICE_1_ENABLED", "false")
	assert.Equal(t, false, helpers.Tuples(EnvironmentFilter{}.Run(context.Background(), event, models.Destination{Config: models.Configuration{Name: "service_1"}}))[0])
	os.Setenv("HOOK_SERVICE_1_ENABLED", "true")
	assert.Equal(t, true, helpers.Tuples(EnvironmentFilter{}.Run(context.Background(), event, models.Destination{Config: models.Configuration{Name: "service_1"}}))[0])
}

func TestEnvironmentFilterApplicable(t *testing.T) {
	event := models.IncomingEvent{}
	assert.Equal(t, true, EnvironmentFilter{}.Applicable(context.Background(), event, models.Destination{Config: models.Configuration{}}))
}
