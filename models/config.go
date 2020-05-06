package models

import (
	"regexp"
	"strconv"
	"time"
)

// Configuration - Webhook Configuration Model
type Configuration struct {
	ConfigID         string   `json:"id"`
	CallbackURL      string   `json:"callback_url"`
	Validate         string   `json:"validate"`
	Actions          []string `json:"actions"`
	Source           string   `json:"source"`
	Throttle         string   `json:"throttle"`
	ThrottleTrailing bool     `json:"throttle_trailing"`
	IncludeDocument  bool     `json:"include_document"`
	Name             string   `json:"name"`
	AllowLoopback    bool     `json:"allow_loopback"`
	Sender           string   `json:"sender"`
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
