package senders

import (
	"encoding/json"

	interfaces "github.com/shoplineapp/captin/interfaces"
	models "github.com/shoplineapp/captin/models"
	log "github.com/sirupsen/logrus"

	aws "github.com/aws/aws-sdk-go/aws"
	aws_session "github.com/aws/aws-sdk-go/aws/session"
	aws_sqs "github.com/aws/aws-sdk-go/service/sqs"
	aws_sqsiface "github.com/aws/aws-sdk-go/service/sqs/sqsiface"
)

var sLogger = log.WithFields(log.Fields{"class": "SqsSender"})

// SqsSender - Send Event to AWS SQS
type SqsSender struct {
	interfaces.EventSenderInterface
	Queue SqsSenderQueue
}

// SqsSenderQueue - Extra struct for mocking SQS
type SqsSenderQueue struct {
	Client aws_sqsiface.SQSAPI
}

func NewSqsSender(awsConfig aws.Config) *SqsSender {
	session := aws_session.Must(aws_session.NewSession(&awsConfig))
	return &SqsSender{
		Queue: SqsSenderQueue{Client: aws_sqs.New(session)},
	}
}

// SendEvent - Send incoming event into SQS queue
func (s *SqsSender) SendEvent(e models.IncomingEvent, d models.Destination) error {
	queueURL := d.GetCallbackURL()
	sLogger.WithFields(log.Fields{"queueURL": queueURL}).Debug("Send sqs event")

	payload, jsonErr := json.Marshal(e)
	if jsonErr != nil {
		sLogger.WithFields(log.Fields{"error": jsonErr}).Error("Failed to convert incoming event to json payload")
		return jsonErr
	}

	_, err := s.Queue.Client.SendMessage(&aws_sqs.SendMessageInput{
		MessageBody: aws.String(string(payload)),
		QueueUrl:    &queueURL,
	})

	if err != nil {
		sLogger.WithFields(log.Fields{"error": err, "event": e, "destination": d}).Error("Failed to send event with SQS")
	}

	return err
}
