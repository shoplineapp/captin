package models_test

import (
	"context"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/trace"

	. "github.com/shoplineapp/captin/v2/models"
)

func mustParseSpanContextConfig(traceIDHex string, spanIDHex string, traceFlags byte, traceStateStr string) trace.SpanContextConfig {
	traceID, err := trace.TraceIDFromHex(traceIDHex)
	if err != nil {
		panic(err)
	}
	spanID, err := trace.SpanIDFromHex(spanIDHex)
	if err != nil {
		panic(err)
	}
	traceState, err := trace.ParseTraceState(traceStateStr)
	if err != nil {
		panic(err)
	}
	return trace.SpanContextConfig{
		TraceID:    traceID,
		SpanID:     spanID,
		TraceFlags: trace.TraceFlags(traceFlags),
		TraceState: traceState,
	}
}

var sampleDistributedTracingInfo = NewDistributedTracingInfoFromContext(trace.ContextWithSpanContext(
	context.Background(),
	trace.NewSpanContext(mustParseSpanContextConfig("11111111111111111111111111111111", "2222222222222222", 0x01, "a=b, c=d")),
))

var sampleIncomingEvent = IncomingEvent{
	Key:        "product.update",
	Source:     "core",
	Payload:    map[string]interface{}{"payload": "data"},
	TargetType: "product",
	TargetId:   "xxxxx",
	Control: map[string]interface{}{
		"extra":        "extra",
		"ts":           99999999999999,
		"host":         "host",
		"ip_addresses": "ip_addresses",
	},
	TargetDocument:         map[string]interface{}{"payload": "data"},
	DistributedTracingInfo: sampleDistributedTracingInfo,
	TraceId:                "11111111-2222-3333-4444-555555555555",
}

func TestNewIncomingEvent(t *testing.T) {
	event_key := "product.update"
	source := "core"
	payload := map[string]interface{}{
		"_id":  "xxxxx",
		"type": "product",
	}
	control := map[string]interface{}{
		"host": "http://example.com",
	}
	target_type := "product"
	target_id := "xxxxx"
	trace_id := uuid.New().String()

	data, _ := json.Marshal(map[string]interface{}{
		"event_key":                event_key,
		"source":                   source,
		"payload":                  payload,
		"control":                  control,
		"target_type":              target_type,
		"target_id":                target_id,
		"trace_id":                 trace_id,
		"distributed_tracing_info": sampleDistributedTracingInfo,
	})
	inc := NewIncomingEvent(data)
	assert.Equal(t, event_key, inc.Key)
	assert.Equal(t, source, inc.Source)
	assert.Equal(t, payload, inc.Payload)
	assert.Equal(t, target_type, inc.TargetType)
	assert.Equal(t, target_id, inc.TargetId)
	assert.Equal(t, control, inc.Control)
	assert.Equal(t, trace_id, inc.TraceId)
	assert.Equal(t, sampleDistributedTracingInfo, inc.DistributedTracingInfo)
}

func TestIsValid(t *testing.T) {
	assert.Equal(t, true, IncomingEvent{Key: "product.update", Source: "core", Payload: map[string]interface{}{}, TargetType: "product", TargetId: "xxxxx"}.IsValid())

	// Test missing required attributes
	assert.Equal(t, false, IncomingEvent{Source: "core", Payload: map[string]interface{}{}, TargetType: "product", TargetId: "xxxxx"}.IsValid())
	assert.Equal(t, false, IncomingEvent{Key: "product.update", Payload: map[string]interface{}{}, TargetType: "product", TargetId: "xxxxx"}.IsValid())
	assert.Equal(t, false, IncomingEvent{Key: "product.update", Source: "core"}.IsValid())

	// Test optional attributes
	assert.Equal(t, true, IncomingEvent{Key: "product.update", Source: "core", TargetType: "product", TargetId: "xxxxx"}.IsValid())
	assert.Equal(t, true, IncomingEvent{Key: "product.update", Source: "core", Payload: map[string]interface{}{"_id": "xxxxx"}}.IsValid())
}

func TestMarshalJSON(t *testing.T) {
	val, err := sampleIncomingEvent.MarshalJSON()
	require.NoError(t, err)
	assert.Equal(t, `{"control":{"host":"host","ip_addresses":"ip_addresses","ts":99999999999999},"id":"xxxxx","key":"product.update","source":"core","trace_id":"11111111-2222-3333-4444-555555555555","type":"product"}`, string(val))
}

func TestString(t *testing.T) {
	val := sampleIncomingEvent.String()
	assert.Equal(t, `{"control":{"host":"host","ip_addresses":"ip_addresses","ts":99999999999999},"id":"xxxxx","key":"product.update","source":"core","trace_id":"11111111-2222-3333-4444-555555555555","type":"product"}`, string(val))
}

func TestToMap(t *testing.T) {
	val := sampleIncomingEvent.ToMap()
	rs := reflect.DeepEqual(val, map[string]interface{}{
		"control":                  map[string]interface{}{"host": "host", "ip_addresses": "ip_addresses", "extra": "extra", "ts": 99999999999999},
		"event_key":                "product.update",
		"payload":                  map[string]interface{}{"payload": "data"},
		"source":                   "core",
		"target_document":          map[string]interface{}{"payload": "data"},
		"target_id":                "xxxxx",
		"target_type":              "product",
		"trace_id":                 "11111111-2222-3333-4444-555555555555",
		"distributed_tracing_info": sampleDistributedTracingInfo,
	})
	assert.Equal(t, true, rs)
}

func TestToJson(t *testing.T) {
	val, err := sampleIncomingEvent.ToJson()
	require.NoError(t, err)
	assert.Equal(t, `{"control":{"extra":"extra","host":"host","ip_addresses":"ip_addresses","ts":99999999999999},"distributed_tracing_info":{"traceparent":"00-11111111111111111111111111111111-2222222222222222-01","tracestate":"a=b,c=d"},"event_key":"product.update","payload":{"payload":"data"},"source":"core","target_document":{"payload":"data"},"target_id":"xxxxx","target_type":"product","trace_id":"11111111-2222-3333-4444-555555555555"}`, string(val))
}
