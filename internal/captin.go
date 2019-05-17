package internal

import (
	"fmt"
	models "github.com/shoplineapp/captin/internal/models"
	outgoing "github.com/shoplineapp/captin/internal/outgoing"
)

type ExecutionError struct {
	Cause string
}

func (e *ExecutionError) Error() string {
	return fmt.Sprintf("ExecutionError: caused by %s", e.Cause)
}

type Captin struct {
	ConfigMap models.ConfigurationMapper
}

func (c Captin) Execute(e models.IncomingEvent) (bool, error) {
	if e.IsValid() != true {
		return false, &ExecutionError{Cause: "invalid incoming event object"}
	}

	configs := c.ConfigMap.ConfigsForKey(e.Key)

	destinations := []outgoing.Destination{}

	for _, config := range configs {
		destinations = append(destinations, outgoing.Destination{Config: config})
	}

	destinations = outgoing.Custom{}.Sift(outgoing.CustomFilters(e), destinations)

	// TODO: Pass event and destinations into dispatcher

	// TODO: return dispatcher instance (?)

	return true, nil
}
