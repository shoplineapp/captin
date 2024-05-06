package destination_filters

import (
	"context"

	models "github.com/shoplineapp/captin/v2/models"
)

var _ DestinationFilterInterface = SourceFilter{}

type SourceFilter struct{}

func (f SourceFilter) Run(ctx context.Context, e models.IncomingEvent, d models.Destination) (bool, error) {
	return e.Source != d.Config.GetSource(), nil
}

func (f SourceFilter) Applicable(ctx context.Context, e models.IncomingEvent, d models.Destination) bool {
	return d.Config.GetAllowLoopback() == false
}
