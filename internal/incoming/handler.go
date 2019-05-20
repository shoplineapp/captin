package incoming

import (
	models "github.com/shoplineapp/captin/internal/models"
)

type Handler interface {
	SetConfigMapper(configMapper *models.ConfigurationMapper)
}
