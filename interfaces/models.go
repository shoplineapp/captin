package interfaces

import (
	models "github.com/shoplineapp/captin/models"
)

// DestinationMiddleware - Interface for third-party application to add extra handling on destinations
type DestinationMiddleware interface {
	Apply(e models.IncomingEvent, d []models.Destination) (models.IncomingEvent, []models.Destination)
}

// DestinationFilter - Interface for third-party application to filter destination by event
type DestinationFilter interface {
	Run(e models.IncomingEvent, c models.Destination) (bool, error)
	Applicable(e models.IncomingEvent, c models.Destination) bool
}

// IncomingHandler - Interface for creating handler to trigger captin execute
type IncomingHandler interface {
	SetConfigMapper(configMapper *ConfigMapperInterface)
	Setup(c CaptinInterface)
}

// ConfigMapperInterface - Interface for config mapper
type ConfigMapperInterface interface {
	ConfigsForKey(eventKey string) []models.Configuration
}
