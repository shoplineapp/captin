package senders

import (
  "encoding/json"
  "fmt"
  "time"

  interfaces "github.com/shoplineapp/captin/interfaces"
  models "github.com/shoplineapp/captin/models"
  beanstalk "github.com/beanstalkd/go-beanstalk"
  log "github.com/sirupsen/logrus"
)

bLogger = log.WithFields(log.Fields{"class": "BeanstalkdSender"})

// BeanstalkdSender - Send Event to beanstalkd
type BeanstalkdSender struct {
  interfaces.EventSenderInterface
}

// SendEvent - #BeanstalkdSender SendEvent
func (c *BeanstalkdSender) SendEvent(e models.IncomingEvent, d models.Destination) error {
  conn, err := beanstalk.Dial("tcp", e.Control["beanstalkd_host"])
  if err != nil {
    bLogger.WithFields(log.Fields{
      "error": err
    }).Error("Beanstalk create connection failed.")
    return err
  }

  tubeName := fmt.Sprintf("%s.%s", e.Control["tube_namespace"], e.Control["queue_name"])
  conn.Tube = beanstalk.Tube { Conn: conn, Name: tubeName }

  jobBody, err := json.Marshal(e.Payload)
  if err != nil {
    bLogger.WithFields(log.Fields{
      "error": err
    }).Error("Beanstalkd job payload format invalid.")
    return err
  }

  pri, delay, ttr := uint32(1), time.Duration(0), time.Duration(0)
  id, err := conn.Put(jobBody, pri, delay, ttr)
  if err != nil {
    bLogger.WithFields(log.Fields{
      "error": err
    }).Error("Beanstalk client put job failed.")
    return err
  }

  bLogger.WithFields(log.Fields{
    "id": id
    "jobBody": string(jobBody)
  }).Info("Enqueue job.")

  defer conn.Close()
  return nil
}

