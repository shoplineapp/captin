package incoming

import (
	internal "captin/internal"
	models "captin/internal/models"
)

type Handler interface {
	SetConfigMapper(configMapper *models.ConfigurationMapper)
	Setup(c internal.CaptinInterface)
}
