package throttles

// Throttle - interface for a throttle object
type Throttle interface {
	// Trigger() - Execute throttled resources
	Trigger()

	// Stop() - Stops this throttler
	Stop()

	// Next() - Returns true at most once per period
	Next() bool
}
