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

	config := c.ConfigMap.ConfigsForKey(e.Key)

	// _ = []outgoing.Destination{}

	// for _, config := range configs {
	// 	append(destinations, &outgoing.Destination{config})
	// }

	// TODO: Pass event and configs into custom to filter out destinations

	// TODO: Pass event and destinations into dispatcher

	// Create dispatcher and dispatch events
	dispatcher := outgoing.NewDispatcherWithConfig(config)
	dispatcher.Dispatch(e)

	return true, nil
}
