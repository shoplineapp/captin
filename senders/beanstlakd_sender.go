package senders

import (
	"encoding/json"
	"fmt"
	"net"
	"regexp"
	"strings"
	"time"

	beanstalk "github.com/beanstalkd/go-beanstalk"
	statsd "github.com/joeycumines/statsd"
	captin_errors "github.com/shoplineapp/captin/errors"
	interfaces "github.com/shoplineapp/captin/interfaces"
	models "github.com/shoplineapp/captin/models"
	log "github.com/sirupsen/logrus"
)

var bLogger = log.WithFields(log.Fields{"class": "BeanstalkdSender"})

// Characters allowed in beanstalkd
const AllowedCharacters = `[A-Za-z0-9\\\-\+\/\;\.\$\_\(\)]+`

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

	beanstalkdHost := e.Control["beanstalkd_host"]
	if beanstalkdHost == nil {
		return &captin_errors.UnretryableError{Msg: "beanstalkd_host is empty", Event: e}
	}

	beanstalkdHostStr := beanstalkdHost.(string)
	if isValidBeanstlakdHost(beanstalkdHostStr) == false {
		return &captin_errors.UnretryableError{Msg: "beanstalkd_host is invalid", Event: e}
	}

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

	queueName := e.Control["queue_name"]
	if queueName == nil {
		return &captin_errors.UnretryableError{Msg: "queue_name is empty", Event: e}
	}

	queueNameStr := queueName.(string)
	isValidQueueName, err := regexp.MatchString(AllowedCharacters, queueNameStr)
	if err != nil || !isValidQueueName {
		return &captin_errors.UnretryableError{Msg: "queue_name is invalid", Event: e}
	}

	conn.Tube = beanstalk.Tube{Conn: conn, Name: queueNameStr}

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

func isValidBeanstlakdHost(addr string) bool {
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
