package outgoing_filters

import (
	models "github.com/shoplineapp/captin/internal/models"
)

type Filter interface {
	// Event models.IncomingEvent
	Run(c models.Configuration) (bool, error)
	Applicable(c models.Configuration) bool
}
