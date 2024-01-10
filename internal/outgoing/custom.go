package outgoing

import (
	"context"

	destination_filters "github.com/shoplineapp/captin/v2/destinations/filters"
	log "github.com/sirupsen/logrus"
)

var cLogger = log.WithFields(log.Fields{"class": "Custom"})

type Custom struct{}

// Sift - Custom check will filter ineligible destination
func (c Custom) Sift(ctx context.Context, e *models.IncomingEvent, destinations []models.Destination, filters []destination_filters.DestinationFilterInterface, middlewares []destination_filters.DestinationMiddlewareInterface) []models.Destination {
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
		}
	}
	for _, m := range middlewares {
		sifted = m.Apply(e, sifted)
	}

	return sifted
}
