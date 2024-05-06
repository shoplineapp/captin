package senders_test

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"testing"

	aws "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	aws_sqs "github.com/aws/aws-sdk-go/service/sqs"
	aws_sqsiface "github.com/aws/aws-sdk-go/service/sqs/sqsiface"
	models "github.com/shoplineapp/captin/v2/models"
	. "github.com/shoplineapp/captin/v2/senders"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type sqsMock struct {
	aws_sqsiface.SQSAPI
	mock.Mock

	SentMessages []aws_sqs.SendMessageInput
}

func (s *sqsMock) SendMessageWithContext(ctx context.Context, input *aws_sqs.SendMessageInput, _ ...request.Option) (*aws_sqs.SendMessageOutput, error) {
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
	sqs.On("SendMessageWithContext", mock.Anything, mock.Anything).Return(nil)

	sender.DefaultClient = sqs
	result := sender.SendEvent(
		context.Background(),
		models.IncomingEvent{},
		models.Destination{
			Config: models.Configuration{},
		},
	)

	assert.Nil(t, result)
	sqs.AssertNumberOfCalls(t, "SendMessageWithContext", 1)
}

func TestSqsSender_SendEvent_Failed(t *testing.T) {
	awsConfig := aws.Config{Region: aws.String("ap-southeast-1")}
	sender := NewSqsSender(awsConfig)

	sqs := new(sqsMock)
	sqs.On("SendMessageWithContext", mock.Anything, mock.Anything).Return(nil)

	sender.DefaultClient = sqs
	result := sender.SendEvent(
		context.Background(),
		models.IncomingEvent{Control: map[string]interface{}{"result": "failed"}},
		models.Destination{
			Config: models.Configuration{Name: "failed"},
		},
	)

	assert.Error(t, result, "some error")
	sqs.AssertNumberOfCalls(t, "SendMessageWithContext", 1)
}

func TestSqsSender_GetClient_UseAccessKey_WithCorrectAwsConfig(t *testing.T) {
	awsConfig := aws.Config{Region: aws.String("ap-southeast-1")}
	sender := NewSqsSender(awsConfig)

	sqs := new(sqsMock)
	sqs.On("SendMessageWithContext", mock.Anything, mock.Anything).Return(nil)

	os.Setenv("HOOK_TEST_DESTINATION_CALLBACK_URL", "https://sqs.ap-southeast-1.amazonaws.com/000000000000/queue")
	os.Setenv("HOOK_TEST_DESTINATION_SQS_SENDER_USE_CUSTOM_CONFIG", "true")
	os.Setenv("HOOK_TEST_DESTINATION_SQS_SENDER_AWS_ENDPOINT", "http://localhost:4566")
	os.Setenv("HOOK_TEST_DESTINATION_SQS_SENDER_AWS_REGION", "ap-southeast-1")
	os.Setenv("HOOK_TEST_DESTINATION_SQS_SENDER_AWS_ACCESS_KEY_ID", "MY_ACCESS_KEY_ID")
	os.Setenv("HOOK_TEST_DESTINATION_SQS_SENDER_AWS_SECRET_ACCESS_KEY", "MY_SECRET_ACCESS_KEY")

	client := sender.GetClient(
		models.Destination{
			Config: models.Configuration{
				Name: "test_destination",
			},
		},
	)

	sqsClient, _ := (client).(*aws_sqs.SQS)
	credentials, _ := sqsClient.Config.Credentials.Get()

	assert.Equal(t, *sqsClient.Config.Region, "ap-southeast-1")
	assert.Equal(t, *sqsClient.Config.Endpoint, "http://localhost:4566")
	assert.Equal(t, credentials.AccessKeyID, "MY_ACCESS_KEY_ID")
	assert.Equal(t, credentials.SecretAccessKey, "MY_SECRET_ACCESS_KEY")
}

func TestSqsSender_SendEvent_UseAccessKey_Success(t *testing.T) {
	os.Setenv("HOOK_TEST_DESTINATION_SQS_SENDER_USE_CUSTOM_CONFIG", "true")

	awsConfig := aws.Config{Region: aws.String("ap-southeast-1")}
	sender := NewSqsSender(awsConfig)

	sqs := new(sqsMock)
	sqs.On("SendMessageWithContext", mock.Anything, mock.Anything).Return(nil)

	sender.DestinationClientMap["test_destination"] = sqs

	result := sender.SendEvent(
		context.Background(),
		models.IncomingEvent{},
		models.Destination{
			Config: models.Configuration{
				Name: "test_destination",
			},
		},
	)

	assert.Nil(t, result)
	sqs.AssertNumberOfCalls(t, "SendMessageWithContext", 1)
}
