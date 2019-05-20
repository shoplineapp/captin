package senders

import (
	"fmt"

	models "captin/internal/models"
)

// ConsoleEventSender - Present Event in console
type ConsoleEventSender struct{}

// SendEvent - #ConsoleEventSender SendEvent
func (c *ConsoleEventSender) SendEvent(e models.IncomingEvent, config models.Configuration) error {
	fmt.Println("Configuration: \t\t", config.Name)
	fmt.Println("Process Event ID: \t", e.TargetId)
	fmt.Println("Process Event Type: \t", e.TargetType)
	return nil
}
