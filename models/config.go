package models

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/shoplineapp/captin/v2/interfaces"
)

var _ interfaces.ConfigurationInterface = &Configuration{}

// Configuration - Webhook Configuration Model
type Configuration struct {
	ConfigID                 string            `json:"id"`
	CallbackURL              string            `json:"callback_url"`
	Validate                 string            `json:"validate"`
	Actions                  []string          `json:"actions"`
	Source                   string            `json:"source"`
	Throttle                 string            `json:"throttle"`
	Delay                    string            `json:"delay"`
	ThrottleTrailingDisabled bool              `json:"throttle_trailing_disabled"`
	KeepThrottledPayloads    bool              `json:"keep_throttled_payloads"`
	KeepThrottledDocuments   bool              `json:"keep_throttled_documents"`
	IncludeDocument          bool              `json:"include_document"`
	Name                     string            `json:"name"`
	AllowLoopback            bool              `json:"allow_loopback"`
	Sender                   string            `json:"sender"`
	DocumentStore            string            `json:"document_store"`
	RetryBackoff             string            `json:"retry_backoff"`
	IncludeDocumentAttrs     []string          `json:"include_document_attrs"`
	ExcludeDocumentAttrs     []string          `json:"exclude_document_attrs"`
	IncludePayloadAttrs      []string          `json:"include_payload_attrs"`
	ExcludePayloadAttrs      []string          `json:"exclude_payload_attrs"`
	Extras                   map[string]string `json:"extras"`
}

func (c Configuration) GetByEnv(key string) (string, string) {
	envKey := fmt.Sprintf("HOOK_%s_%s", strings.ToUpper(c.Name), strings.ToUpper(key))
	return envKey, os.Getenv(envKey)
}

// GetThrottleValue - Get Throttle Value in millisecond
func (c Configuration) GetThrottleValue() time.Duration {
	return c.GetTimeValueMillis(c.Throttle)
}

// GetDelayValue - Get delay time in millisecond
func (c Configuration) GetDelayValue() time.Duration {
	return c.GetTimeValueMillis(c.Delay)
}

// GetTimeValueMillis - Get millisecond from time value string
func (c Configuration) GetTimeValueMillis(timeValue string) time.Duration {
	match := regexp.MustCompile("(\\d+(?:\\.\\d+)?)(s|ms|m|h)")
	res := match.FindAllStringSubmatch(timeValue, -1)

	for i := range res {
		value, err := strconv.Atoi(res[i][1])

		if err != nil {
			panic(err)
		}

		unit := res[i][2]

		switch unit {
		case "ms":
			return time.Millisecond * time.Duration(value)
		case "s":
			return time.Second * time.Duration(value)
		case "m":
			return time.Minute * time.Duration(value)
		case "h":
			return time.Hour * time.Duration(value)
		default:
			panic("unrecognized time unit")
		}
	}

	return 0
}

func (c Configuration) GetActions() []string {
	return c.Actions
}

func (c Configuration) GetConfigID() string {
	return c.ConfigID
}

func (c Configuration) GetCallbackURL() string {
	return c.CallbackURL
}

func (c Configuration) GetValidate() string {
	return c.Validate
}

func (c Configuration) GetSource() string {
	return c.Source
}

func (c Configuration) GetThrottle() string {
	return c.Throttle
}

func (c Configuration) GetDelay() string {
	return c.Delay
}

func (c Configuration) GetThrottleTrailingDisabled() bool {
	return c.ThrottleTrailingDisabled
}

func (c Configuration) GetKeepThrottledPayloads() bool {
	return c.KeepThrottledPayloads
}

func (c Configuration) GetKeepThrottledDocuments() bool {
	return c.KeepThrottledDocuments
}

func (c Configuration) GetIncludeDocument() bool {
	return c.IncludeDocument
}

func (c Configuration) GetName() string {
	return c.Name
}

func (c Configuration) GetAllowLoopback() bool {
	return c.AllowLoopback
}

func (c Configuration) GetSender() string {
	return c.Sender
}

func (c Configuration) GetDocumentStore() string {
	return c.DocumentStore
}

func (c Configuration) GetRetryBackoff() []string {
	return strings.Split(c.RetryBackoff, ",")
}

func (c Configuration) GetIncludeDocumentAttrs() []string {
	return c.IncludeDocumentAttrs
}

func (c Configuration) GetExcludeDocumentAttrs() []string {
	return c.ExcludeDocumentAttrs
}

func (c Configuration) GetIncludePayloadAttrs() []string {
	return c.IncludePayloadAttrs
}

func (c Configuration) GetExcludePayloadAttrs() []string {
	return c.ExcludePayloadAttrs
}

func (c Configuration) GetExtras() map[string]string {
	return c.Extras
}
