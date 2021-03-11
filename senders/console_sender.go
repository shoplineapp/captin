package senders

import (
	interfaces "github.com/shoplineapp/captin/interfaces"
	models "github.com/shoplineapp/captin/models"
	log "github.com/sirupsen/logrus"
)

var cLogger = log.WithFields(log.Fields{"class": "ConsoleEventSender"})

// ConsoleEventSender - Present Event in console
type ConsoleEventSender struct{}

// SendEvent - #ConsoleEventSender SendEvent
func (c *ConsoleEventSender) SendEvent(ev interfaces.IncomingEventInterface, dv interfaces.DestinationInterface) error {
	e := ev.(models.IncomingEvent)
	d := dv.(models.Destination)

	cLogger.WithFields(log.Fields{
		"config_name": d.Config.GetName(),
		"target_id":   e.TargetId,
		"target_type": e.TargetType,
		"target_document": e.TargetDocument,
	}).Debug("Event sent")
	return nil
}
