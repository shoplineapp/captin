package senders_test

import (
	"encoding/json"
	"errors"
	"testing"

	aws "github.com/aws/aws-sdk-go/aws"
	aws_sqs "github.com/aws/aws-sdk-go/service/sqs"
	aws_sqsiface "github.com/aws/aws-sdk-go/service/sqs/sqsiface"
	models "github.com/shoplineapp/captin/models"
	. "github.com/shoplineapp/captin/senders"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type sqsMock struct {
	aws_sqsiface.SQSAPI
	mock.Mock

	SentMessages []aws_sqs.SendMessageInput
}

func (s *sqsMock) SendMessage(input *aws_sqs.SendMessageInput) (*aws_sqs.SendMessageOutput, error) {
	if s.SentMessages == nil {
		s.SentMessages = []aws_sqs.SendMessageInput{}
	}
	s.SentMessages = append(s.SentMessages, *input)
	_ = s.Called(input)

	payload := map[string]interface{}{}
	json.Unmarshal([]byte(*input.MessageBody), &payload)

	// Throw error manually
	if payload != nil && payload["control"] != nil && payload["control"].(map[string]interface{})["result"] == "failed" {
		return nil, errors.New("SQSError: some error")
	}
	return nil, nil
}

func TestSqsSender_SendEvent_Success(t *testing.T) {
	awsConfig := aws.Config{Region: aws.String("ap-southeast-1")}
	sender := NewSqsSender(awsConfig)

	sqs := new(sqsMock)
	sqs.On("SendMessage", mock.Anything).Return(nil)

	sender.Queue = SqsSenderQueue{Client: sqs}
	result := sender.SendEvent(
		models.IncomingEvent{},
		models.Destination{
			models.Configuration{},
		},
	)

	assert.Nil(t, result)
	sqs.AssertNumberOfCalls(t, "SendMessage", 1)
}

func TestSqsSender_SendEvent_Failed(t *testing.T) {
	awsConfig := aws.Config{Region: aws.String("ap-southeast-1")}
	sender := NewSqsSender(awsConfig)

	sqs := new(sqsMock)
	sqs.On("SendMessage", mock.Anything).Return(nil)

	sender.Queue = SqsSenderQueue{Client: sqs}
	result := sender.SendEvent(
		models.IncomingEvent{Control: map[string]interface{}{"result": "failed"}},
		models.Destination{
			models.Configuration{Name: "failed"},
		},
	)

	assert.Error(t, result, "some error")
	sqs.AssertNumberOfCalls(t, "SendMessage", 1)
}
