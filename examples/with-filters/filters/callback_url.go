package filters

import (
	interfaces "github.com/shoplineapp/captin/interfaces"
)

type CallbackUrlFilter struct {
	interfaces.CustomFilter
}

func (f SourceFilter) Run(e models.IncomingEvent, d models.Destination) (bool, error) {
	return e.Source != d.Config.Source, nil
}

func (f SourceFilter) Applicable(e models.IncomingEvent, d models.Destination) bool {
	return d.Config.AllowLoopback == false
}
