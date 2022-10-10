package senders

import (
	"bytes"
	"crypto/tls"
	"io/ioutil"
	"net/http"

	interfaces "github.com/shoplineapp/captin/interfaces"
	models "github.com/shoplineapp/captin/models"
	log "github.com/sirupsen/logrus"
)

var hLogger = log.WithFields(log.Fields{"class": "HttpEventSender"})

// HTTPResponse - HTTP Response
type HTTPResponse struct {
	url      string
	response *http.Response
	err      error
}

// HTTPEventSender - Send Event through HTTP
type HTTPEventSender struct {
	interfaces.EventSenderInterface
}

// SendEvent - #HttpEventSender SendEvent
func (c *HTTPEventSender) SendEvent(ev interfaces.IncomingEventInterface, dv interfaces.DestinationInterface) error {
	e := ev.(models.IncomingEvent)
	d := dv.(models.Destination)

	url := d.GetCallbackURL()
	payload, err := e.ToJson()

	if err != nil {
		return err
	}

	// TODO: Read from config
	req, reqErr := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if reqErr != nil {
		return reqErr
	}
	req.Header.Set("Content-Type", "application/json")

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{
		Transport: tr,
	}

	res, resErr := client.Do(req)
	if resErr != nil {
		return resErr
	}
	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)
	hLogger.WithFields(log.Fields{"result": string(body)}).Debug("Send http event with result")

	return nil
}
