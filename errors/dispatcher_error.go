package errors

import (
	"fmt"

	interfaces "github.com/shoplineapp/captin/v2/interfaces"
	models "github.com/shoplineapp/captin/v2/models"
)

// DispatcherError - Error when send events
type DispatcherError struct {
	interfaces.ErrorInterface

	Msg         string
	Event       models.IncomingEvent
	Destination models.Destination
}

func (e DispatcherError) Error() string {
	return fmt.Sprintf("DispatcherError: %s", e.Msg)
}
