package dispatcher_delayers

import (
	"fmt"

	"github.com/shoplineapp/captin/dispatcher"
	interfaces "github.com/shoplineapp/captin/interfaces"
	log "github.com/sirupsen/logrus"
)

type GoroutineDelayer struct {
	interfaces.DispatchDelayerInterface
}

var dLogger = log.WithFields(log.Fields{"class": "Goroutine"})

func (d GoroutineDelayer) Execute(evt interfaces.IncomingEventInterface, dest interfaces.DestinationInterface, exec func()) {
	config := dest.GetConfig()

	outstandingDelaySecondsStr := ""
	if evt.GetOutstandingDelaySeconds() > 0 {
		outstandingDelaySecondsStr = fmt.Sprintf("%.0f", evt.GetOutstandingDelaySeconds().Seconds())
	}

	eventLogger := dLogger.WithFields(log.Fields{
		"event":                     evt.GetTraceInfo(),
		"hook_name":                 config.GetName(),
		"hook_delay":                config.GetDelayValue(),
		"outstanding_delay_seconds": outstandingDelaySecondsStr,
	})

	eventLogger.Debug(fmt.Sprintf("Event delayed by GoroutineDelayer"))
	ch := make(chan int, 1)
	go dispatcher.TrackAfterFuncJob(config.GetDelayValue(), func() {
		eventLogger.Info(fmt.Sprintf("Event resumed"))
		exec()
		ch <- 1
	})
	<-ch // waiting for delayed execution
}
