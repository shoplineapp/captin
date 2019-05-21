package filters

import (
	interfaces "github.com/shoplineapp/captin/interfaces"
	models "github.com/shoplineapp/captin/models"
	"strings"
)

type CallbackUrlFilter struct {
	interfaces.CustomFilter
}

// Run - Callback url filter will filter out callback_url starting without https
func (f CallbackUrlFilter) Run(e models.IncomingEvent, d models.Destination) (bool, error) {
	return strings.HasPrefix("d.Config.CallbackURL", "https"), nil
}

func (f CallbackUrlFilter) Applicable(e models.IncomingEvent, d models.Destination) bool {
	return true
}
