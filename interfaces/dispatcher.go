package interfaces

// DispatchDelayerInterface - class for deciding how event message is delayed
type DispatchDelayerInterface interface {
  Execute(e IncomingEventInterface, d DestinationInterface, exec func()) ()
}
