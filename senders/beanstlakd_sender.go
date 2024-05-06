package senders

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"regexp"
	"strings"
	"time"

	beanstalk "github.com/beanstalkd/go-beanstalk"
	statsd "github.com/joeycumines/statsd"
	captin_errors "github.com/shoplineapp/captin/v2/errors"
	interfaces "github.com/shoplineapp/captin/v2/interfaces"
	"github.com/shoplineapp/captin/v2/internal/helpers"
	models "github.com/shoplineapp/captin/v2/models"
	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

var bLogger = log.WithFields(log.Fields{"class": "BeanstalkdSender"})

// Characters allowed in go-beanstalkd
// Source: https://github.com/beanstalkd/go-beanstalk/blob/master/name.go
const allowedCharacters = `^[A-Za-z0-9\\\-\+\/\;\.\$\_\(\)]{1,199}$`

var _ interfaces.EventSenderInterface = &BeanstalkdSender{}

// BeanstalkdSender - Send Event to beanstalkd
type BeanstalkdSender struct {
	StatsdClient *statsd.Client
}

// SendEvent - #BeanstalkdSender SendEvent
func (c *BeanstalkdSender) SendEvent(ctx context.Context, ev interfaces.IncomingEventInterface, dv interfaces.DestinationInterface) (err error) {
	ctx, span := helpers.Tracer().Start(ctx, "captin.BeanstalkdSender.SendEvent")
	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
		span.End()
	}()

	e := ev.(models.IncomingEvent)
	d := dv.(models.Destination)

	if e.Control == nil {
		bLogger.Error("Event control is empty")

		if c.StatsdClient != nil {
			c.StatsdClient.Increment(fmt.Sprintf("hook.sender.beanstalkd.error,metricname=%s,hook=%s,code=EmptyControl", d.Config.GetName(), d.Config.GetName()))
		}

		return &captin_errors.UnretryableError{Msg: "Event control is empty", Event: e}
	}

	beanstalkdHost := e.Control["beanstalkd_host"]
	if beanstalkdHost == nil || beanstalkdHost == "" {
		bLogger.WithFields(log.Fields{
			"Event": e,
		}).Error("beanstalkd_host is empty")

		if c.StatsdClient != nil {
			c.StatsdClient.Increment(fmt.Sprintf("hook.sender.beanstalkd.error,metricname=%s,hook=%s,code=EmptyBeanstalkdHostName", d.Config.GetName(), d.Config.GetName()))
		}

		//return &captin_errors.UnretryableError{Msg: "beanstalkd_host is empty", Event: e}
	}

	beanstalkdHostStr := beanstalkdHost.(string)
	if isValidBeanstalkdHost(beanstalkdHostStr) == false {
		bLogger.WithFields(log.Fields{
			"Event": e,
		}).Error("beanstalkd_host is invalid")

		if c.StatsdClient != nil {
			c.StatsdClient.Increment(fmt.Sprintf("hook.sender.beanstalkd.error,metricname=%s,hook=%s,code=InvalidBeanstalkdHostName", d.Config.GetName(), d.Config.GetName()))
		}

		//return &captin_errors.UnretryableError{Msg: "beanstalkd_host is invalid", Event: e}
	}
	span.SetAttributes(attribute.String("beanstalkdHost", beanstalkdHostStr))

	conn, err := beanstalk.Dial("tcp", beanstalkdHostStr)
	if err != nil {
		bLogger.WithFields(log.Fields{
			"error": err,
		}).Error("Beanstalk create connection failed.")
		if c.StatsdClient != nil {
			c.StatsdClient.Increment(fmt.Sprintf("hook.sender.beanstalkd.error,metricname=%s,hook=%s,code=CreateConnectionFailed", d.Config.GetName(), d.Config.GetName()))
		}
		return err
	}
	defer conn.Close()

	beanstalkdQueueName := e.Control["queue_name"]
	if beanstalkdQueueName == nil || beanstalkdQueueName == "" {
		bLogger.WithFields(log.Fields{
			"Event": e,
		}).Error("queue_name for beanstalkd sender is empty")

		if c.StatsdClient != nil {
			c.StatsdClient.Increment(fmt.Sprintf("hook.sender.beanstalkd.error,metricname=%s,hook=%s,code=EmptyBeanstalkdQueueName", d.Config.GetName(), d.Config.GetName()))
		}

		//return &captin_errors.UnretryableError{Msg: "queue_name for beanstalkd sender is empty", Event: e}
	}

	beanstalkdQueueNameStr := beanstalkdQueueName.(string)
	span.SetAttributes(attribute.String("beanstalkdQueueName", beanstalkdQueueNameStr))
	isValidBeanstalkdQueueNameStr, err := regexp.MatchString(allowedCharacters, beanstalkdQueueNameStr)
	if err != nil || !isValidBeanstalkdQueueNameStr {
		bLogger.WithFields(log.Fields{
			"Event": e,
		}).Error("queue_name for beanstalkd sender is invalid")

		if c.StatsdClient != nil {
			c.StatsdClient.Increment(fmt.Sprintf("hook.sender.beanstalkd.error,metricname=%s,hook=%s,code=InvalidBeanstalkdQueueName", d.Config.GetName(), d.Config.GetName()))
		}

		//return &captin_errors.UnretryableError{Msg: "queue_name for beanstalkd sender is invalid", Event: e}
	}

	conn.Tube = beanstalk.Tube{Conn: conn, Name: beanstalkdQueueNameStr}

	e.DistributedTracingInfo.InjectContext(ctx)
	jobBody, err := json.Marshal(e.Payload)
	if err != nil {
		bLogger.WithFields(log.Fields{
			"error": err,
		}).Error("Beanstalkd job payload format invalid.")
		return err
	}

	pri := uint32(65536)
	var delay time.Duration
	ttr := time.Duration(time.Minute)

	if e.Control["priority"] != nil {
		pri = e.Control["priority"].(uint32)
	}

	if e.Control["delay"] != nil {
		delay, _ = time.ParseDuration(e.Control["delay"].(string))
	}

	if e.Control["ttr"] != nil {
		ttr, _ = time.ParseDuration(e.Control["ttr"].(string))
	}

	id, err := conn.Put(jobBody, pri, time.Duration(delay), time.Duration(ttr))
	if err != nil {
		bLogger.WithFields(log.Fields{
			"error": err,
		}).Error("Beanstalk client put job failed.")
		if c.StatsdClient != nil {
			c.StatsdClient.Increment(fmt.Sprintf("hook.sender.beanstalkd.error,metricname=%s,hook=%s,code=PutJobFailed", d.Config.GetName(), d.Config.GetName()))
		}
		return err
	}

	bLogger.WithFields(log.Fields{
		"id":      id,
		"pri":     pri,
		"delay":   delay,
		"ttr":     ttr,
		"jobBody": string(jobBody),
	}).Info("Enqueue job.")

	return nil
}

func isValidBeanstalkdHost(addr string) bool {
	host, _, err := net.SplitHostPort(addr)

	// Check is contain port
	if err != nil {
		return false
	}

	ip := net.ParseIP(host)

	if ip == nil {
		return !(strings.HasPrefix(host, "http://") || strings.HasPrefix(host, "https://"))
	}
	return true
}
