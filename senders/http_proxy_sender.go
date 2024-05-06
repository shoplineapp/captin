package senders

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"net/http"

	interfaces "github.com/shoplineapp/captin/v2/interfaces"
	"github.com/shoplineapp/captin/v2/internal/helpers"
	models "github.com/shoplineapp/captin/v2/models"
	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
)

var hpLogger = log.WithFields(log.Fields{"class": "HTTPProxyEventSender"})

var _ interfaces.EventSenderInterface = &HTTPProxyEventSender{}

// HTTPProxyEventSender - Send Event through HTTP with payload only
// Different from HTTPEventSender, which parses the whole event body
// in order to pass event meta data to destinations,
// HTTPProxyEventSender only parses payload for general usage of
// third party API calls.
type HTTPProxyEventSender struct{}

func (c *HTTPProxyEventSender) SendEvent(ctx context.Context, ev interfaces.IncomingEventInterface, dv interfaces.DestinationInterface) (err error) {
	ctx, span := helpers.Tracer().Start(ctx, "captin.HTTPProxyEventSender.SendEvent")
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
	// clear the distributed tracing context before sending the request to external receivers
	e.DistributedTracingInfo.ClearContext()
	payload, err := json.Marshal(e.Payload)

	if err != nil {
		return err
	}

	req, reqErr := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(payload))
	if reqErr != nil {
		return reqErr
	}
	req.Header.Set("Content-Type", "application/json")

	// For tracing the event of sending the request
	// For external receivers, we don't want to propagate the trace context, while we still want to trace the event of sending the request
	// Therefore we use a CompositeTextMapPropagator with no base propagators as a no-op propagator
	noopPropagator := propagation.NewCompositeTextMapPropagator()
	client := &http.Client{
		Transport: otelhttp.NewTransport(&http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}, otelhttp.WithPropagators(noopPropagator)),
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
