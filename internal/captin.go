package internal

import (
	"fmt"

	models "captin/internal/models"
	outgoing "captin/internal/outgoing"
	senders "captin/internal/senders"
)

type ExecutionError struct {
	Cause string
}

func (e *ExecutionError) Error() string {
	return fmt.Sprintf("ExecutionError: caused by %s", e.Cause)
}

type Captin struct {
	ConfigMap models.ConfigurationMapper
}

func (c Captin) Execute(e models.IncomingEvent) (bool, error) {
	if e.IsValid() != true {
		return false, &ExecutionError{Cause: "invalid incoming event object"}
	}

	config := c.ConfigMap.ConfigsForKey(e.Key)

	// _ = []outgoing.Destination{}

	// for _, config := range configs {
	// 	append(destinations, &outgoing.Destination{config})
	// }

	// TODO: Pass event and configs into custom to filter out destinations

	// TODO: Pass event and destinations into dispatcher

	// Create dispatcher and dispatch events
	sender := senders.HTTPEventSender{}
	dispatcher := outgoing.NewDispatcherWithConfig(config, &sender)
	dispatcher.Dispatch(e)

	for _, err := range dispatcher.Errors {
		switch dispatcherErr := err.(type) {
		case *outgoing.DispatcherError:
			fmt.Println("[Dispatcher] Error on event: ", dispatcherErr.Event.TargetId)
			fmt.Println("[Dispatcher] Error on event type: ", dispatcherErr.Event.TargetType)
		default:
			fmt.Println(e)
		}
	}

	return true, nil
}