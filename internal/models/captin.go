package models

import (
	"fmt"
)

type ExecutionError struct {
	Cause string
}

func (e *ExecutionError) Error() string {
	return fmt.Sprintf("ExecutionError: caused by %s", e.Cause)
}

type Captin struct {
	ConfigMap ConfigurationMapper
}

func (c Captin) Execute(e IncomingEvent) (bool, error) {
	if e.IsValid() != true {
		return false, &ExecutionError{Cause: "invalid incoming event object"}
	}

	_ = c.ConfigMap.ConfigsForKey(e.Key)

	// TODO: Pass event and configs into custom to filter out distinations

	// TODO: Pass event and distinations into dispatcher

	// TODO: return dispatcher instance (?)

	return true, nil
}
