package senders

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	models "github.com/shoplineapp/captin/internal/models"
)

// HTTPResponse - HTTP Response
type HTTPResponse struct {
	url      string
	response *http.Response
	err      error
}

// HTTPEventSender - Send Event through HTTP
type HTTPEventSender struct{}

// SendEvent - #HttpEventSender SendEvent
func (c *HTTPEventSender) SendEvent(e models.IncomingEvent, d models.Destination) error {
	url := d.Config.CallbackURL
	payload, err := json.Marshal(e)
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

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)
	fmt.Println("response Body:", string(body))

	return nil
}
