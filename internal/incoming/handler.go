package incoming

import (
	interfaces "github.com/shoplineapp/captin/interfaces"
	models "github.com/shoplineapp/captin/internal/models"
)

type Handler interface {
	SetConfigMapper(configMapper *models.ConfigurationMapper)
	Setup(c interfaces.CaptinInterface)
}
