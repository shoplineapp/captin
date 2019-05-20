package outgoing

import (
	interfaces "captin/interfaces"
	models "captin/internal/models"
)

type Custom struct{}

// Sift - Custom check will filter ineligible destination
func (c Custom) Sift(e models.IncomingEvent, filters []interfaces.CustomFilter, destinations []models.Destination) []models.Destination {
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
	return sifted
}
