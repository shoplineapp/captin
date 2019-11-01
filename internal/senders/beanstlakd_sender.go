package senders

import (
  "encoding/json"
  models "github.com/shoplineapp/captin/models"
  beanstalk "github.com/beanstalkd/go-beanstalk"
  "fmt"
  "time"
  log "github.com/sirupsen/logrus"
)


// BeanstalkdSender - Send Event to beanstalkd
type BeanstalkdSender struct{}

// SendEvent - #BeanstalkdSender SendEvent
func (c *BeanstalkdSender) SendEvent(e models.IncomingEvent, d models.Destination) error {
  host := fmt.Sprintf("%s:%s", e.Control["beanstalkd_host"], e.Control["beanstalkd_port"])
  conn, err := beanstalk.Dial("tcp", host)

  if err != nil {
    log.Info("Beanstalk create connection failed: ", err)
    panic(err)
  }

  tubeName := fmt.Sprintf("%s.%s", e.Control["tube_namespace"], e.Control["queue_name"])
  conn.Tube = beanstalk.Tube { Conn: conn, Name: tubeName }
  jobBody, err := json.Marshal(e.Payload)
  if err != nil {
    log.Info(fmt.Sprintf("* Beanstalkd job payload format invalid: %s", string(jobBody)))
    panic(err)
  }
  pri, delay, ttr := uint32(1), time.Duration(0), time.Duration(0)
  id, err := conn.Put(jobBody, pri, delay, ttr)
  if err != nil {
    log.Info(fmt.Sprintf("* Beanstalk client put job failed: %s", string(jobBody)))
    panic(err)
  }
  log.Info(fmt.Sprintf("Enqueue job id: %d, body: %s", id, string(jobBody)))
  defer conn.Close()
  return nil
}

