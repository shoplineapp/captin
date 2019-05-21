package outgoing_filters

import (
	models "github.com/shoplineapp/captin/models"
)

type Filter interface {
	Run(e models.IncomingEvent, c models.Configuration) (bool, error)
	Applicable(e models.IncomingEvent, c models.Configuration) bool
}
