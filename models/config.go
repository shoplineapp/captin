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
	ConfigID                 string   `json:"id"`
	CallbackURL              string   `json:"callback_url"`
	Validate                 string   `json:"validate"`
	Actions                  []string `json:"actions"`
	Source                   string   `json:"source"`
	Throttle                 string   `json:"throttle"`
	ThrottleTrailingDisabled bool     `json:"throttle_trailing_disabled"`
	IncludeDocument          bool     `json:"include_document"`
	Name                     string   `json:"name"`
	AllowLoopback            bool     `json:"allow_loopback"`
	Sender                   string   `json:"sender"`
	IncludeDocumentAttrs     []string   `json:"include_document_attrs"`
	ExcludeDocumentAttrs     []string   `json:"exclude_document_attrs"`
	IncludePayloadAttrs      []string   `json:"include_payload_attrs"`
	ExcludePayloadAttrs      []string   `json:"exclude_payload_attrs"`
}

func (c Configuration) GetByEnv(key string) (string, string) {
	envKey := fmt.Sprintf("HOOK_%s_%s", strings.ToUpper(c.Name), strings.ToUpper(key))
	return envKey, os.Getenv(envKey)
}

// GetThrottleValue - Get Throttle Value in millisecond
func (c Configuration) GetThrottleValue() time.Duration {
	match := regexp.MustCompile("(\\d+(?:\\.\\d+)?)(s|ms)")
	res := match.FindAllStringSubmatch(c.Throttle, -1)

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
