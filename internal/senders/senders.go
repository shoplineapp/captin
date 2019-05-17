package senders

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	models "github.com/shoplineapp/captin/internal/models"
)

// EventSenderInterface - Event Sender Interface
type EventSenderInterface interface {
	SendEvent(e models.IncomingEvent, config models.Configuration) error
}

// ConsoleEventSender - Present Event in console
type ConsoleEventSender struct{}

// SendEvent - #ConsoleEventSender SendEvent
func (c *ConsoleEventSender) SendEvent(e models.IncomingEvent, config models.Configuration) error {
	fmt.Println("Configuration: \t\t", config.Name)
	fmt.Println("Process Event ID: \t", e.TargetId)
	fmt.Println("Process Event Type: \t", e.TargetType)
	return nil
}

// HTTPResponse - HTTP Response
type HTTPResponse struct {
	url      string
	response *http.Response
	err      error
}

// HTTPEventSender - Send Event through HTTP
type HTTPEventSender struct{}

// SendEvent - #HttpEventSender SendEvent
func (c *HTTPEventSender) SendEvent(e models.IncomingEvent, config models.Configuration) error {
	url := config.CallbackURL
	payload, err := json.Marshal(e)
	if err != nil {
		return err
	}

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
