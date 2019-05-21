package incoming

import (
	core "github.com/shoplineapp/captin/core"
	interfaces "github.com/shoplineapp/captin/interfaces"
)

type Handler interface {
	SetConfigMapper(configMapper *core.ConfigurationMapper)
	Setup(c interfaces.CaptinInterface)
}
