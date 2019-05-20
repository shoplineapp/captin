package incoming

import (
	interfaces "captin/interfaces"
	models "captin/internal/models"
)

type Handler interface {
	SetConfigMapper(configMapper *models.ConfigurationMapper)
	Setup(c interfaces.CaptinInterface)
}
