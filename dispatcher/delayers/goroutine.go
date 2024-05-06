package dispatcher_delayers

import (
	"context"
	"fmt"

	"time"

	"github.com/shoplineapp/captin/v2/dispatcher"
	"github.com/shoplineapp/captin/v2/interfaces"
	"github.com/shoplineapp/captin/v2/models"
	log "github.com/sirupsen/logrus"
)

var _ interfaces.DispatchDelayerInterface = GoroutineDelayer{}

type GoroutineDelayer struct{}

var dLogger = log.WithFields(log.Fields{"class": "Goroutine"})

func (d GoroutineDelayer) Execute(ctx context.Context, evt interfaces.IncomingEventInterface, dest interfaces.DestinationInterface, exec func()) {
	event := d.TapDelayedEvent(evt.(models.IncomingEvent), dest.(models.Destination))
	config := dest.GetConfig()

	delay, outstanding := d.GetDelayAndOutstandingSeconds(event, dest.(models.Destination))
	eventLogger := dLogger.WithFields(log.Fields{
		"event":                     event,
		"hook_name":                 config.GetName(),
		"hook_delay":                config.GetDelayValue(),
		"event_delay":               delay,
		"outstanding_delay_seconds": outstanding,
	})

	eventLogger.Debug("Event delayed by GoroutineDelayer")
	ch := make(chan int, 1)
	dispatcher.TrackAfterFuncJob(config.GetDelayValue(), func() {
		eventLogger.Info("Event resumed")
		exec()
		ch <- 1
	})
	<-ch // waiting for delayed execution
}

func (d GoroutineDelayer) TapDelayedEvent(evt models.IncomingEvent, dest models.Destination) models.IncomingEvent {
	if evt.Control == nil {
		evt.Control = map[string]interface{}{}
	}
	_, outstanding := d.GetDelayAndOutstandingSeconds(evt, dest)
	evt.Control["outstanding_delay_seconds"] = fmt.Sprintf("%.0f", outstanding)
	evt.Control["desired_hooks"] = []string{dest.Config.GetName()}

	// Unset target document loaded from dispatcher to prevent exceed of delayed message payload limit
	evt.TargetDocument = map[string]interface{}{}
	return evt
}

func (d GoroutineDelayer) GetDelayAndOutstandingSeconds(evt models.IncomingEvent, dest models.Destination) (float64, float64) {
	config := dest.Config
	delay := float64(config.GetDelayValue() / time.Second)
	outstanding := float64(evt.GetOutstandingDelaySeconds() / time.Second)
	if outstanding < 0 {
		outstanding = delay
	}
	return float64(delay), float64(outstanding)
}
