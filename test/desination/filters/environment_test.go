package destination_filters_test

import (
  "github.com/stretchr/testify/assert"
  "testing"
  "os"

  . "github.com/shoplineapp/captin/destinations/filters"
  helpers "github.com/shoplineapp/captin/internal/helpers"
  models "github.com/shoplineapp/captin/models"
)

func TestEnvironmentFilterRunValidate(t *testing.T) {
  defer os.Unsetenv("HOOK_SERVICE_1_ENABLED")
  event := models.IncomingEvent{}
  assert.Equal(t, true, helpers.Tuples(EnvironmentFilter{}.Run(event, models.Destination{Config: models.Configuration{Name: "service_1"}}))[0])
  os.Setenv("HOOK_SERVICE_1_ENABLED", "false")
  assert.Equal(t, false, helpers.Tuples(EnvironmentFilter{}.Run(event, models.Destination{Config: models.Configuration{Name: "service_1"}}))[0])
  os.Setenv("HOOK_SERVICE_1_ENABLED", "true")
  assert.Equal(t, true, helpers.Tuples(EnvironmentFilter{}.Run(event, models.Destination{Config: models.Configuration{Name: "service_1"}}))[0])
}

func TestEnvironmentFilterApplicable(t *testing.T) {
  event := models.IncomingEvent{}
  assert.Equal(t, true, EnvironmentFilter{}.Applicable(event, models.Destination{Config: models.Configuration{}}))
}
