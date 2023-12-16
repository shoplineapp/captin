package senders_test

import (
	"github.com/shoplineapp/captin/errors"
	"strings"
	"testing"

	models "github.com/shoplineapp/captin/models"

	"github.com/shoplineapp/captin/senders"
	"github.com/stretchr/testify/assert"
)

func TestBeanstalkdSender_SendEvent_BeanstalkdHost(t *testing.T) {
	tests := map[string]struct {
		isNilInput bool
		input      string
		want       error
	}{
		"WithIPv4AndPort":             {input: "127.0.0.1:11300", want: nil},
		"WithIPv6AndPort":             {input: "[0:0:0:0:0:0:0:1]:11300", want: nil},
		"WithURLAndPort":              {input: "localhost:11300", want: nil},
		"WithSubdomainAndPort":        {input: "subdomain.localhost:11300", want: nil},
		"WithIPv4AndWithoutPort":      {input: "127.0.0.1", want: errors.UnretryableError{Msg: "beanstalkd_host is invalid"}},
		"WithIPv6AndWithoutPort":      {input: "[0:0:0:0:0:0:0:1]", want: errors.UnretryableError{Msg: "beanstalkd_host is invalid"}},
		"WithURLAndWithoutPort":       {input: "localhost", want: errors.UnretryableError{Msg: "beanstalkd_host is invalid"}},
		"WithSubdomainAndWithoutPort": {input: "subdomain.localhost", want: errors.UnretryableError{Msg: "beanstalkd_host is invalid"}},
		"WithEmptyString":             {input: "", want: errors.UnretryableError{Msg: "beanstalkd_host is empty"}},
		"WithoutHost":                 {isNilInput: true, want: errors.UnretryableError{Msg: "beanstalkd_host is empty"}},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			sender := new(senders.BeanstalkdSender)

			var beanstalkdHost string
			if !tc.isNilInput {
				beanstalkdHost = tc.input
			}

			got := sender.SendEvent(
				models.IncomingEvent{
					Control: map[string]interface{}{
						"beanstalkd_host": beanstalkdHost,
						"queue_name":      "foo",
					},
				},
				models.Destination{
					Config: models.Configuration{},
				},
			)

			if tc.want == nil {
				assert.Nil(t, got, "Should not throw error")

			} else {
				assert.EqualError(t, got, tc.want.Error(), "Should throw UnretryableError")
			}

		})
	}
}

func TestBeanstalkdSender_SendEvent_QueueName(t *testing.T) {
	tests := map[string]struct {
		isNilInput bool
		input      string
		want       error
	}{
		"WithValidName":         {input: "foo", want: nil},
		"WithSpecialSymbolName": {input: "foo_bar", want: nil},
		"WithInvalidCharacter":  {input: "foo_!", want: errors.UnretryableError{Msg: "queue_name for beanstalkd sender is invalid"}},
		"WithVeryLongName":      {input: strings.Repeat("a", 300), want: errors.UnretryableError{Msg: "queue_name for beanstalkd sender is invalid"}},
		"WithEmptyString":       {input: "", want: errors.UnretryableError{Msg: "queue_name for beanstalkd sender is empty"}},
		"WithoutHost":           {isNilInput: true, want: errors.UnretryableError{Msg: "queue_name for beanstalkd sender is empty"}},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			sender := new(senders.BeanstalkdSender)

			var queueName string
			if !tc.isNilInput {
				queueName = tc.input
			}

			got := sender.SendEvent(
				models.IncomingEvent{
					Control: map[string]interface{}{
						"beanstalkd_host": "127.0.0.1:11300",
						"queue_name":      queueName,
					},
				},
				models.Destination{
					Config: models.Configuration{},
				},
			)

			if tc.want == nil {
				assert.Nil(t, got, "Should not throw error")

			} else {
				assert.EqualError(t, got, tc.want.Error(), "Should throw UnretryableError")
			}

		})
	}
}
