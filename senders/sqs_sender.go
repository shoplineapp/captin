package senders

import (
	"context"

	interfaces "github.com/shoplineapp/captin/v2/interfaces"
	"github.com/shoplineapp/captin/v2/internal/helpers"
	models "github.com/shoplineapp/captin/v2/models"
	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"

	aws "github.com/aws/aws-sdk-go/aws"
	aws_credentials "github.com/aws/aws-sdk-go/aws/credentials"
	aws_session "github.com/aws/aws-sdk-go/aws/session"
	aws_sqs "github.com/aws/aws-sdk-go/service/sqs"
	aws_sqsiface "github.com/aws/aws-sdk-go/service/sqs/sqsiface"
)

var sLogger = log.WithFields(log.Fields{"class": "SqsSender"})

var _ interfaces.EventSenderInterface = &SqsSender{}

// SqsSender - Send Event to AWS SQS
type SqsSender struct {
	DefaultClient        aws_sqsiface.SQSAPI
	DestinationClientMap map[string]aws_sqsiface.SQSAPI
}

func NewSqsSender(defaultAwsConfig aws.Config) *SqsSender {
	defaultSession := aws_session.Must(aws_session.NewSession(&defaultAwsConfig))
	return &SqsSender{
		DefaultClient:        aws_sqs.New(defaultSession),
		DestinationClientMap: map[string]aws_sqsiface.SQSAPI{},
	}
}

// SendEvent - Send incoming event into SQS queue
func (s *SqsSender) SendEvent(ctx context.Context, ev interfaces.IncomingEventInterface, dv interfaces.DestinationInterface) (err error) {
	ctx, span := helpers.Tracer().Start(ctx, "captin.SqsSender.SendEvent")
	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
		span.End()
	}()

	e := ev.(models.IncomingEvent)
	d := dv.(models.Destination)

	queueURL := d.GetCallbackURL()
	span.SetAttributes(attribute.String("queueURL", queueURL))
	sLogger.WithFields(log.Fields{"queueURL": queueURL}).Debug("Send sqs event")

	e.DistributedTracingInfo.InjectContext(ctx)
	payload, jsonErr := e.ToJson()
	if jsonErr != nil {
		sLogger.WithFields(log.Fields{"error": jsonErr}).Error("Failed to convert incoming event to json payload")
		return jsonErr
	}

	_, err = s.GetClient(dv).SendMessageWithContext(ctx, &aws_sqs.SendMessageInput{
		MessageBody: aws.String(string(payload)),
		QueueUrl:    &queueURL,
	})

	if err != nil {
		sLogger.WithFields(log.Fields{"error": err, "event": e, "destination": d}).Error("Failed to send event with SQS")
	}

	return err
}

func (s *SqsSender) GetClient(dv interfaces.DestinationInterface) aws_sqsiface.SQSAPI {
	d := dv.(models.Destination)
	destName := d.Config.GetName()

	if dv.GetSqsSenderConfig("USE_CUSTOM_CONFIG") == "true" {
		_, queueInitialized := s.DestinationClientMap[destName]
		if !queueInitialized {
			awsConfig := aws.Config{}

			if dv.GetSqsSenderConfig("AWS_ENDPOINT") != "" {
				awsConfig.Endpoint = aws.String(dv.GetSqsSenderConfig("AWS_ENDPOINT"))
			}

			if dv.GetSqsSenderConfig("AWS_REGION") != "" {
				awsConfig.Region = aws.String(dv.GetSqsSenderConfig("AWS_REGION"))
			}

			if dv.GetSqsSenderConfig("AWS_ACCESS_KEY_ID") != "" && dv.GetSqsSenderConfig("AWS_SECRET_ACCESS_KEY") != "" {
				awsConfig.Credentials = aws_credentials.NewStaticCredentials(dv.GetSqsSenderConfig("AWS_ACCESS_KEY_ID"), dv.GetSqsSenderConfig("AWS_SECRET_ACCESS_KEY"), "")
			}

			session := aws_session.Must(aws_session.NewSession(&awsConfig))
			s.DestinationClientMap[destName] = aws_sqs.New(session)
		}

		return s.DestinationClientMap[destName]
	}

	return s.DefaultClient
}
