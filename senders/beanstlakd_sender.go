package senders

import (
	"encoding/json"
	"fmt"
	"time"

	beanstalk "github.com/beanstalkd/go-beanstalk"
	statsd "github.com/joeycumines/statsd"
	captin_errors "github.com/shoplineapp/captin/errors"
	interfaces "github.com/shoplineapp/captin/interfaces"
	models "github.com/shoplineapp/captin/models"
	log "github.com/sirupsen/logrus"
)

var bLogger = log.WithFields(log.Fields{"class": "BeanstalkdSender"})

// BeanstalkdSender - Send Event to beanstalkd
type BeanstalkdSender struct {
	interfaces.EventSenderInterface
	StatsdClient *statsd.Client
}

// SendEvent - #BeanstalkdSender SendEvent
func (c *BeanstalkdSender) SendEvent(ev interfaces.IncomingEventInterface, dv interfaces.DestinationInterface) error {
	e := ev.(models.IncomingEvent)
	d := dv.(models.Destination)

	if e.Control == nil {
		bLogger.Error("Event control is empty")
		return &captin_errors.UnretryableError{Msg: "Event control is empty", Event: e}
	}

	conn, err := beanstalk.Dial("tcp", e.Control["beanstalkd_host"].(string))
	if err != nil {
		bLogger.WithFields(log.Fields{
			"error": err,
		}).Error("Beanstalk create connection failed.")
		if c.StatsdClient != nil {
			c.StatsdClient.Increment(fmt.Sprintf("hook.sender.beanstalkd.error,metricname=%s,hook=%s,code=CreateConnectionFailed", d.Config.GetName(), d.Config.GetName()))
		}
		return err
	}

	conn.Tube = beanstalk.Tube{Conn: conn, Name: e.Control["queue_name"].(string)}

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

	defer conn.Close()
	return nil
}
