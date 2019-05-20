package outgoing

import (
	models "captin/internal/models"
	. "captin/internal/outgoing/filters"
)

type Custom struct{}

// Sift - Custom check will filter ineligible destination
func (c Custom) Sift(e models.IncomingEvent, filters []Filter, destinations []models.Destination) []models.Destination {
	sifted := []models.Destination{}
	for _, destination := range destinations {
		eligible := true
		for _, filter := range filters {
			config := destination.Config
			if eligible == false || filter.Applicable(e, config) == false {
				continue
			}
			valid, _ := filter.Run(e, config)
			if valid != true {
				eligible = false
			}
		}
		if eligible {
			sifted = append(sifted, destination)
		}
	}
	return sifted
}
