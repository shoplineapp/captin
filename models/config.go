package models

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

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
	match := regexp.MustCompile("(\\d+(?:\\.\\d+)?)(s|ms)")
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
		default:
			panic("unrecognized time unit")
		}
	}

	return 0
}
