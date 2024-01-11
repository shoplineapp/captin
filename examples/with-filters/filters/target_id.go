package filters

import (
	interfaces "github.com/shoplineapp/captin/v2/interfaces"
	models "github.com/shoplineapp/captin/v2/models"
)

type TargetIdFilter struct {
	interfaces.DestinationFilter
}

// Run - Callback url filter will filter out callback_url starting without https
func (f TargetIdFilter) Run(e models.IncomingEvent, d models.Destination) (bool, error) {
	return e.TargetId == "A", nil
}

func (f TargetIdFilter) Applicable(e models.IncomingEvent, d models.Destination) bool {
	return true
}
