package outgoing_filters

import (
	models "captin/internal/models"
)

type SourceFilter struct {
	Filter
}

func (f SourceFilter) Run(e models.IncomingEvent, c models.Configuration) (bool, error) {
	return e.Source != c.Source, nil
}

func (f SourceFilter) Applicable(e models.IncomingEvent, c models.Configuration) bool {
	return c.AllowLoopback == false
}
