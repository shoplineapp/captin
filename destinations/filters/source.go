package destination_filters

import (
	interfaces "github.com/shoplineapp/captin/interfaces"
	models "github.com/shoplineapp/captin/models"
)

type SourceFilter struct {
	interfaces.DestinationFilter
}

func (f SourceFilter) Run(e models.IncomingEvent, d models.Destination) (bool, error) {
	return e.Source != d.Config.Source, nil
}

func (f SourceFilter) Applicable(e models.IncomingEvent, d models.Destination) bool {
	return d.Config.AllowLoopback == false
}
