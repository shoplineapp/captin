package outgoing

import (
	models "github.com/shoplineapp/captin/internal/models"
	. "github.com/shoplineapp/captin/internal/outgoing/filters"
)

type Custom struct{}

// Sift - Custom check will filter ineligible destination
func (c Custom) Sift(filters []Filter, destinations []Destination) []Destination {
	sifted := []Destination{}
	for _, destination := range destinations {
		eligible := true
		for _, filter := range filters {
			config := destination.Config
			if eligible == false || filter.Applicable(config) == false {
				continue
			}
			valid, _ := filter.Run(config)
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

func CustomFilters(e models.IncomingEvent) []Filter {
	return []Filter{
		ValidateFilter{Event: e},
		SourceFilter{Event: e},
	}
}
