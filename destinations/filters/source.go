package destination_filters

import (
	"context"

	models "github.com/shoplineapp/captin/v2/models"
)

type SourceFilter struct {
	DestinationFilterInterface
}

func (f SourceFilter) Run(ctx context.Context, e models.IncomingEvent, d models.Destination) (bool, error) {
	return e.Source != d.Config.GetSource(), nil
}

func (f SourceFilter) Applicable(ctx context.Context, e models.IncomingEvent, d models.Destination) bool {
	return d.Config.GetAllowLoopback() == false
}
