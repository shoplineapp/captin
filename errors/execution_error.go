package errors

import (
	"fmt"
)

// ExecutionError - Error on executing events
type ExecutionError struct {
	ErrorInterface
	Cause string
}

func (e ExecutionError) Error() string {
	return fmt.Sprintf("ExecutionError: caused by %s", e.Cause)
}
