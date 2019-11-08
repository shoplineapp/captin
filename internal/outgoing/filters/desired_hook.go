package outgoing_filters

import (
	interfaces "github.com/shoplineapp/captin/interfaces"
	models "github.com/shoplineapp/captin/models"
)

func ispresent(str string, list []string) bool {
	for _, v := range list {
		if v == str {
			return true
		}
	}
	return false
}

// DesiredHookFilter - Filter destination if given event has desired destination
type DesiredHookFilter struct {
	interfaces.DestinationFilter
}

// Run - Get desired hooks in control and filter out exclusion
func (f DesiredHookFilter) Run(e models.IncomingEvent, d models.Destination) (bool, error) {
	return ispresent(d.Config.Name, e.Control["desired_hooks"].([]string)), nil
}

// Applicable - Check if desired hooks is present
func (f DesiredHookFilter) Applicable(e models.IncomingEvent, d models.Destination) bool {
	return e.Control["desired_hooks"] != nil
}
