package models

import (
	interfaces "github.com/shoplineapp/captin/interfaces"
	"fmt"
	"time"
)

// Destination - Event dispatch destination
type Destination struct {
	interfaces.DestinationInterface

	Config interfaces.ConfigurationInterface
	callbackUrl string
}

func (d Destination) GetConfig() interfaces.ConfigurationInterface {
	return d.Config
}

func (d *Destination) SetCallbackURL(url string) {
	d.callbackUrl = url
}

func (d Destination) GetCallbackURL() string {
	_, value := d.Config.GetByEnv("callback_url")
	if len(value) > 0 {
		return value
	}
	if len(d.callbackUrl) > 0 {
		return d.callbackUrl
	}
	return d.Config.GetCallbackURL()
}

func (d Destination) GetSqsSenderConfig(key string) string {
	_, value := d.Config.GetByEnv(fmt.Sprintf("SQS_SENDER_%s", key))
	return value
}

func (d Destination) GetDocumentStore() string {
	_, value := d.Config.GetByEnv("document_store")
	if len(value) == 0 {
		if len(d.Config.GetDocumentStore()) == 0 {
			return "default"
		}
		return d.Config.GetDocumentStore()
	}
	return value
}

func (d Destination) RequireDelay(evt interfaces.IncomingEventInterface) bool {
	if (d.Config.GetDelayValue() <= time.Duration(0) ||
            evt.GetOutstandingDelaySeconds() == time.Duration(0)) {
		return false
	}

	return true
}
