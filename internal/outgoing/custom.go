package outgoing

import (
	interfaces "github.com/shoplineapp/captin/interfaces"
	models "github.com/shoplineapp/captin/internal/models"
)

type Custom struct{}

// Sift - Custom check will filter ineligible destination
func (c Custom) Sift(e models.IncomingEvent, destinations []models.Destination, filters []interfaces.CustomFilter, middlewares []interfaces.CustomMiddleware) []models.Destination {
	sifted := []models.Destination{}
	for _, destination := range destinations {
		eligible := true
		for _, filter := range filters {
			if eligible == false || filter.Applicable(e, destination) == false {
				continue
			}
			valid, _ := filter.Run(e, destination)
			if valid != true {
				eligible = false
			}
		}
		if eligible {
			sifted = append(sifted, destination)
		}
	}
	for _, m := range middlewares {
		e, sifted = m.Apply(e, sifted)
	}
	return sifted
}
