package outgoing

import (
	interfaces "github.com/shoplineapp/captin/interfaces"
	models "github.com/shoplineapp/captin/models"
	log "github.com/sirupsen/logrus"
)

var cLogger = log.WithFields(log.Fields{"class": "Custom"})

type Custom struct{}

// Sift - Custom check will filter ineligible destination
func (c Custom) Sift(e *models.IncomingEvent, destinations []models.Destination, filters []interfaces.DestinationFilter, middlewares []interfaces.DestinationMiddleware) []models.Destination {
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
			if eligible == false || filter.Applicable(*e, destination) == false {
				continue
			}
			valid, _ := filter.Run(*e, destination)
			if valid != true {
				eligible = false
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
