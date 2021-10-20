package models

import (
	interfaces "github.com/shoplineapp/captin/interfaces"
	"os"
	"fmt"
	"time"
	"strings"
	"strconv"
)

// Destination - Event dispatch destination
type Destination struct {
	interfaces.DestinationInterface

	Config interfaces.ConfigurationInterface
	callbackUrl string
}

var DEFAULT_RETRY_BACKOFF_SECONDS int64 = 10

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

func (d Destination) GetRetryBackoffSeconds(evt interfaces.IncomingEventInterface) int64 {
	globalRetryBackoffSeconds := os.Getenv("APP_GLOBAL_RETRY_BACKOFF_SECONDS")
	backoffConfig := trimArray(d.Config.GetRetryBackoff())
	if len(backoffConfig) <= 0 && globalRetryBackoffSeconds != "" {
		backoffConfig = trimArray(strings.Split(globalRetryBackoffSeconds, ","))
	}

	if len(backoffConfig) == 0 {
		return DEFAULT_RETRY_BACKOFF_SECONDS
	}

	control := evt.GetControl()
	retryCount, _ := control["retry_count"].(float64)
	if (len(backoffConfig) > 0 && len(backoffConfig) <= int(retryCount)) {
		lastConfig := backoffConfig[(len(backoffConfig) - 1):]
		seconds, pErr := strconv.ParseInt(lastConfig[0], 10, 64)
		if pErr != nil {
			return DEFAULT_RETRY_BACKOFF_SECONDS
		}
		return seconds
	}

	seconds, err := strconv.ParseInt(backoffConfig[int(retryCount)], 10, 64)
	if err != nil {
		return DEFAULT_RETRY_BACKOFF_SECONDS
	}
	return seconds
}


func trimArray(arr []string) []string {
	var r []string
	for _, str := range arr {
		if str != "" {
			r = append(r, str)
		}
	}
	return r
}
