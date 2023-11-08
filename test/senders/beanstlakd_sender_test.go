package senders_test

import (
	"testing"

	models "github.com/shoplineapp/captin/models"

	"github.com/shoplineapp/captin/senders"
	"github.com/stretchr/testify/assert"
)

func TestBeanstalkdSender_SendEvent_Success(t *testing.T) {
	sender := new(senders.BeanstalkdSender)
	result := sender.SendEvent(
		models.IncomingEvent{
			Control: map[string]interface{}{
				"beanstalkd_host": "127.0.0.1:11300",
				"queue_name":      "foo",
			},
		},
		models.Destination{
			Config: models.Configuration{},
		},
	)

	assert.Nil(t, result)
}

func TestBeanstalkdSender_SendEvent_Failed_WithEmptyHost(t *testing.T) {
	sender := new(senders.BeanstalkdSender)
	result := sender.SendEvent(
		models.IncomingEvent{
			Control: map[string]interface{}{
				"queue_name": "foo",
			},
		},
		models.Destination{
			Config: models.Configuration{},
		},
	)

	assert.EqualError(t, result, "UnretryableError: beanstalkd_host is empty", "Should throw UnretryableError")
}

func TestBeanstalkdSender_SendEvent_Failed_WithInvalidHost(t *testing.T) {
	sender := new(senders.BeanstalkdSender)
	result := sender.SendEvent(
		models.IncomingEvent{
			Control: map[string]interface{}{
				"beanstalkd_host": "http://localhost:11300",
				"queue_name":      "foo",
			},
		},
		models.Destination{
			Config: models.Configuration{},
		},
	)
	assert.EqualError(t, result, "UnretryableError: beanstalkd_host is invalid", "Should throw UnretryableError")
}

func TestBeanstalkdSender_SendEvent_Failed_WithEmptyQueueName(t *testing.T) {
	sender := new(senders.BeanstalkdSender)
	result := sender.SendEvent(
		models.IncomingEvent{
			Control: map[string]interface{}{
				"beanstalkd_host": "127.0.0.1:11300",
			},
		},
		models.Destination{
			Config: models.Configuration{},
		},
	)
	assert.EqualError(t, result, "UnretryableError: queue_name is empty", "Should throw UnretryableError")
}
func TestBeanstalkdSender_SendEvent_Failed_WithInvalidQueueName(t *testing.T) {
	sender := new(senders.BeanstalkdSender)
	result := sender.SendEvent(
		models.IncomingEvent{
			Control: map[string]interface{}{
				"beanstalkd_host": "127.0.0.1:11300",
				"queue_name":      "!",
			},
		},
		models.Destination{
			Config: models.Configuration{},
		},
	)
	assert.EqualError(t, result, "UnretryableError: queue_name is invalid", "Should throw UnretryableError")
}
