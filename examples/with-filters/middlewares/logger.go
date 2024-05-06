package middlewares

import (
	"fmt"
	models "github.com/shoplineapp/captin/v2/models"
)

type LoggerMiddleware struct{}

func (m LoggerMiddleware) Apply(e *models.IncomingEvent, d []models.Destination) (models.IncomingEvent, []models.Destination) {
	fmt.Printf("= Handling event with %d destinations\n", len(d))
	return e, d
}
