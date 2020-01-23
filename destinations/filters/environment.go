package destination_filters

import (
  interfaces "github.com/shoplineapp/captin/interfaces"
  models "github.com/shoplineapp/captin/models"
  "fmt"
  log "github.com/sirupsen/logrus"
  "os"
  "strings"
)

var eLogger = log.WithFields(log.Fields{"class": "EnvironmentFilter"})

type EnvironmentFilter struct {
  interfaces.DestinationFilter
}

// Destination needs to be enabled by ENV Variable {Config Name}_ENABLED, e.g, WAPOS_SYNC_ENABLED
func (f EnvironmentFilter) Run(e models.IncomingEvent, d models.Destination) (bool, error) {
	variableName := fmt.Sprintf("HOOK_%s_ENABLED", strings.ToUpper(d.Config.Name))
	isEnabled := os.Getenv(variableName) != "false"

	if isEnabled == false {
		eLogger.WithFields(log.Fields{"variableName": variableName}).Debug("Hook disabled by ENV. Destination ignored.")
	}

  return isEnabled, nil
}

func (f EnvironmentFilter) Applicable(e models.IncomingEvent, d models.Destination) bool {
  return true
}
