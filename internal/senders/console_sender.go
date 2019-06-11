package senders

import (
	models "github.com/shoplineapp/captin/models"
	log "github.com/sirupsen/logrus"
)

var cLogger = log.WithFields(log.Fields{"class": "ConsoleEventSender"})

// ConsoleEventSender - Present Event in console
type ConsoleEventSender struct{}

// SendEvent - #ConsoleEventSender SendEvent
func (c *ConsoleEventSender) SendEvent(e models.IncomingEvent, d models.Destination) error {
	cLogger.WithFields(log.Fields{
		"config_name": d.Config.Name,
		"target_id":   e.TargetId,
		"target_type": e.TargetType,
	}).Debug("Event sent")
	return nil
}
