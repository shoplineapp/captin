package destination_filters

import (
	"context"

	models "github.com/shoplineapp/captin/v2/models"
)

// DestinationMiddleware - Interface for third-party application to add extra handling on destinations
type DestinationMiddlewareInterface interface {
	Apply(ctx context.Context, e *models.IncomingEvent, d []models.Destination) []models.Destination
}

// DestinationFilter - Interface for third-party application to filter destination by event
type DestinationFilterInterface interface {
	Run(ctx context.Context, e models.IncomingEvent, c models.Destination) (bool, error)
	Applicable(ctx context.Context, e models.IncomingEvent, c models.Destination) bool
}
