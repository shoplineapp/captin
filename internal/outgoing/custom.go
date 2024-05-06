package outgoing

import (
	"context"

	destination_filters "github.com/shoplineapp/captin/v2/destinations/filters"
	"github.com/shoplineapp/captin/v2/internal/helpers"
	models "github.com/shoplineapp/captin/v2/models"
	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

var cLogger = log.WithFields(log.Fields{"class": "Custom"})

type Custom struct{}

// Sift - Custom check will filter ineligible destination
func (c Custom) Sift(ctx context.Context, e *models.IncomingEvent, destinations []models.Destination, filters []destination_filters.DestinationFilterInterface, middlewares []destination_filters.DestinationMiddlewareInterface) []models.Destination {
	ctx, span := helpers.Tracer().Start(ctx, "captin.Custom.Sift")
	defer span.End()
	cLogger.WithFields(log.Fields{
		"event":        e,
		"destinations": destinations,
		"filters":      filters,
		"middlewares":  middlewares,
	}).Debug("Custom sift with filters and middlewares")
	sifted := []models.Destination{}
	for _, destination := range destinations {
		eligible := true
		for _, filter := range filters {
			if !filter.Applicable(ctx, *e, destination) {
				continue
			}
			valid, _ := filter.Run(ctx, *e, destination)
			if !valid {
				eligible = false
				break
			}
		}
		if eligible {
			sifted = append(sifted, destination)
		} else {
			span.AddEvent("destination removed", trace.WithAttributes(attribute.String("destination", destination.Config.GetName())))
		}
	}
	for _, m := range middlewares {
		sifted = m.Apply(ctx, e, sifted)
	}

	return sifted
}
