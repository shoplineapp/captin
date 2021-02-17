package destination_filters

import (
	models "github.com/shoplineapp/captin/models"
	log "github.com/sirupsen/logrus"
)

var eLogger = log.WithFields(log.Fields{"class": "EnvironmentFilter"})

type EnvironmentFilter struct {
	DestinationFilterInterface
}

// Destination needs to be enabled by ENV Variable {Config Name}_ENABLED, e.g, WAPOS_SYNC_ENABLED
func (f EnvironmentFilter) Run(e models.IncomingEvent, d models.Destination) (bool, error) {
	variableName, value := d.Config.GetByEnv("enabled")
	isEnabled := value != "false"

	if isEnabled == false {
		eLogger.WithFields(log.Fields{"variableName": variableName}).Debug("Hook disabled by ENV. Destination ignored.")
	}

	return isEnabled, nil
}

func (f EnvironmentFilter) Applicable(e models.IncomingEvent, d models.Destination) bool {
	return true
}
