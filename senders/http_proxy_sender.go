package senders

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"net/http"

	interfaces "github.com/shoplineapp/captin/v2/interfaces"
	models "github.com/shoplineapp/captin/v2/models"
	log "github.com/sirupsen/logrus"
)

var hpLogger = log.WithFields(log.Fields{"class": "HTTPProxyEventSender"})

// HTTPProxyResponse - HTTP Response
type HTTPProxyResponse struct {
	url      string
	response *http.Response
	err      error
}

// HTTPProxyEventSender - Send Event through HTTP with payload only
// Different from HTTPEventSender, which parses the whole event body
// in order to pass event meta data to destinations,
// HTTPProxyEventSender only parses payload for general usage of 
// third party API calls.
func (c *HTTPProxyEventSender) SendEvent(ctx context.Context, ev interfaces.IncomingEventInterface, dv interfaces.DestinationInterface) (err error) {
	e := ev.(models.IncomingEvent)
	d := dv.(models.Destination)

	url := d.GetCallbackURL()
	payload, err := json.Marshal(e.Payload)

	if err != nil {
		return err
	}

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
	hpLogger.WithFields(log.Fields{"result": string(body)}).Debug("Send http event with result")

	return nil
}
