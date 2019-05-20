package incoming

import (
	models "captin/internal/models"
)

type Handler interface {
	SetConfigMapper(configMapper *models.ConfigurationMapper)
}
