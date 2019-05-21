package senders

import (
	"fmt"

	models "github.com/shoplineapp/captin/models"
)

// ConsoleEventSender - Present Event in console
type ConsoleEventSender struct{}

// SendEvent - #ConsoleEventSender SendEvent
func (c *ConsoleEventSender) SendEvent(e models.IncomingEvent, d models.Destination) error {
	fmt.Println("Configuration: \t\t", d.Config.Name)
	fmt.Println("Process Event ID: \t", e.TargetId)
	fmt.Println("Process Event Type: \t", e.TargetType)
	return nil
}
