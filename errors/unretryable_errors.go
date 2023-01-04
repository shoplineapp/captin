package errors

import (
	"fmt"

	interfaces "github.com/shoplineapp/captin/interfaces"
	models "github.com/shoplineapp/captin/models"
)

type UnretryableError struct {
	interfaces.ErrorInterface

	Msg         string
	Event       models.IncomingEvent
	Destination models.Destination
}

func (e UnretryableError) Error() string {
	return fmt.Sprintf("UnretryableError: %s", e.Msg)
}
