package destination_filters

import (
	models "github.com/shoplineapp/captin/models"
)

type SourceFilter struct {
	DestinationFilterInterface
}

func (f SourceFilter) Run(e models.IncomingEvent, d models.Destination) (bool, error) {
	return e.Source != d.Config.GetSource(), nil
}

func (f SourceFilter) Applicable(e models.IncomingEvent, d models.Destination) bool {
	return d.Config.GetAllowLoopback() == false
}
