package interfaces

import "context"

// DispatchDelayerInterface - class for deciding how event message is delayed
type DispatchDelayerInterface interface {
	Execute(ctx context.Context, e IncomingEventInterface, d DestinationInterface, exec func())
}
