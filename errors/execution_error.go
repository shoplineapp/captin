package errors

import (
	"fmt"
	interfaces "github.com/shoplineapp/captin/interfaces"
)

// ExecutionError - Error on executing events
type ExecutionError struct {
	interfaces.ErrorInterface

	Cause string
}

func (e ExecutionError) Error() string {
	return fmt.Sprintf("ExecutionError: caused by %s", e.Cause)
}
