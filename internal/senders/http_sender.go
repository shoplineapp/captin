package senders

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	models "captin/internal/models"
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

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)
	fmt.Println("response Body:", string(body))

	return nil
}
