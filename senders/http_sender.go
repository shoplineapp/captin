package senders

import (
	"bytes"
	"context"
	"crypto/tls"
	"io/ioutil"
	"net/http"

	interfaces "github.com/shoplineapp/captin/v2/interfaces"
	"github.com/shoplineapp/captin/v2/internal/helpers"
	models "github.com/shoplineapp/captin/v2/models"
	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/codes"
)

var hLogger = log.WithFields(log.Fields{"class": "HttpEventSender"})

var _ interfaces.EventSenderInterface = &HTTPEventSender{}

// HTTPEventSender - Send Event through HTTP
type HTTPEventSender struct{}

func (c *HTTPEventSender) SendEvent(ctx context.Context, ev interfaces.IncomingEventInterface, dv interfaces.DestinationInterface) (err error) {
	ctx, span := helpers.Tracer().Start(ctx, "captin.HTTPEventSender.SendEvent")
	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
		span.End()
	}()

	e := ev.(models.IncomingEvent)
	d := dv.(models.Destination)

	url := d.GetCallbackURL()
	e.DistributedTracingInfo.InjectContext(ctx)
	payload, err := e.ToJson()

	if err != nil {
		return err
	}

	req, reqErr := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(payload))
	if reqErr != nil {
		return reqErr
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		// for tracing
		Transport: otelhttp.NewTransport(&http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}),
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
