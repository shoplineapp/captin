package throttles

// Throttle - interface for a throttle object
// Throttle event flow:
// Mutex Lock:		1
// Payload Store:	p p p p (Expired * 2)
// Events:				x x x x -
// Throttle: 			  t     1
type Throttle interface {
	// Trigger() - Execute throttled resources
	Trigger()

	// Stop() - Stops this throttler
	Stop()

	// Next() - Returns true at most once per period
	Next() bool
}
