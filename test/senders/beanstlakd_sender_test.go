package senders_test

import (
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
		haveError  bool
	}{
		"WithIPv4AndPort":             {input: "127.0.0.1:11300", haveError: false},
		"WithIPv6AndPort":             {input: "[0:0:0:0:0:0:0:1]:11300", haveError: false},
		"WithURLAndPort":              {input: "localhost:11300", haveError: false},
		"WithSubdomainAndPort":        {input: "subdomain.localhost:11300", haveError: false},
		"WithIPv4AndWithoutPort":      {input: "127.0.0.1", haveError: true},
		"WithIPv6AndWithoutPort":      {input: "[0:0:0:0:0:0:0:1]", haveError: true},
		"WithURLAndWithoutPort":       {input: "localhost", haveError: true},
		"WithSubdomainAndWithoutPort": {input: "subdomain.localhost", haveError: true},
		"WithEmptyString":             {input: "", haveError: true},
		"WithoutHost":                 {isNilInput: true, haveError: true},
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

			if tc.haveError == false {
				assert.Nil(t, got, "Should not throw error")
			} else {
				assert.NotNil(t, got, "Should throw error")
			}

		})
	}
}

func TestBeanstalkdSender_SendEvent_QueueName(t *testing.T) {
	tests := map[string]struct {
		isNilInput bool
		input      string
		haveError  bool
	}{
		"WithValidName":         {input: "foo", haveError: false},
		"WithSpecialSymbolName": {input: "foo_bar", haveError: false},
		"WithInvalidCharacter":  {input: "foo_!", haveError: true},
		"WithVeryLongName":      {input: strings.Repeat("a", 300), haveError: true},
		"WithEmptyString":       {input: "", haveError: true},
		"WithoutHost":           {isNilInput: true, haveError: true},
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

			if tc.haveError == false {
				assert.Nil(t, got, "Should not throw error")
			} else {
				assert.NotNil(t, got, "Should throw error")
			}

		})
	}
}
