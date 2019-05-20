package outgoing_filters

import (
	models "captin/internal/models"
)

type SourceFilter struct {
	Filter
	Event models.IncomingEvent
}

func (f SourceFilter) Run(c models.Configuration) (bool, error) {
	return f.Event.Source != c.Source, nil
}

func (f SourceFilter) Applicable(c models.Configuration) bool {
	return c.AllowLoopback == false
}
