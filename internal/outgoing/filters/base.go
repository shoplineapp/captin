package outgoing_filters

import (
	models "captin/internal/models"
)

type Filter interface {
	Run(e models.IncomingEvent, c models.Configuration) (bool, error)
	Applicable(e models.IncomingEvent, c models.Configuration) bool
}
