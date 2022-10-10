package models_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	. "github.com/shoplineapp/captin/models"
)

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

	data, _ := json.Marshal(map[string]interface{}{
		"event_key":   event_key,
		"source":      source,
		"payload":     payload,
		"control":     control,
		"target_type": target_type,
		"target_id":   target_id,
	})
	inc := NewIncomingEvent(data)
	assert.Equal(t, event_key, inc.Key)
	assert.Equal(t, source, inc.Source)
	assert.Equal(t, payload, inc.Payload)
	assert.Equal(t, target_type, inc.TargetType)
	assert.Equal(t, target_id, inc.TargetId)
	assert.Equal(t, control, inc.Control)
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
	e := IncomingEvent{Key: "product.update", Source: "core", Payload: map[string]interface{}{"payload": "data"}, TargetType: "product", TargetId: "xxxxx", Control: map[string]interface{}{"extra": "extra", "ts": 99999999999999, "host": "host", "ip_addresses": "ip_addresses"}, TargetDocument: map[string]interface{}{"payload": "data"}}
	val, _ := e.MarshalJSON()
	assert.Equal(t, `{"control":{"host":"host","ip_addresses":"ip_addresses","ts":99999999999999},"event_key":"product.update","source":"core","target_id":"xxxxx","target_type":"product","trace_id":""}`, string(val))
}

func TestString(t *testing.T) {
	e := IncomingEvent{Key: "product.update", Source: "core", Payload: map[string]interface{}{"payload": "data"}, TargetType: "product", TargetId: "xxxxx", Control: map[string]interface{}{"extra": "extra", "ts": 99999999999999, "host": "host", "ip_addresses": "ip_addresses"}, TargetDocument: map[string]interface{}{"payload": "data"}}
	val := e.String()
	assert.Equal(t, `{"control":{"host":"host","ip_addresses":"ip_addresses","ts":99999999999999},"event_key":"product.update","source":"core","target_id":"xxxxx","target_type":"product","trace_id":""}`, val)
}

func TestToMap(t *testing.T) {
	e := IncomingEvent{Key: "product.update", Source: "core", Payload: map[string]interface{}{"payload": "data"}, TargetType: "product", TargetId: "xxxxx", Control: map[string]interface{}{"extra": "extra", "ts": 99999999999999, "host": "host", "ip_addresses": "ip_addresses"}, TargetDocument: map[string]interface{}{"payload": "data"}}
	val := e.ToMap()
	rs := reflect.DeepEqual(val, map[string]interface{}{
		"control":         map[string]interface{}{"host": "host", "ip_addresses": "ip_addresses", "extra": "extra", "ts": 99999999999999},
		"event_key":       "product.update",
		"payload":         map[string]interface{}{"payload": "data"},
		"source":          "core",
		"target_document": map[string]interface{}{"payload": "data"},
		"target_id":       "xxxxx",
		"target_type":     "product",
		"trace_id":        "",
	})
	assert.Equal(t, true, rs)
}

func TestToJson(t *testing.T) {
	e := IncomingEvent{Key: "product.update", Source: "core", Payload: map[string]interface{}{"payload": "data"}, TargetType: "product", TargetId: "xxxxx", Control: map[string]interface{}{"extra": "extra", "ts": 99999999999999, "host": "host", "ip_addresses": "ip_addresses"}, TargetDocument: map[string]interface{}{"payload": "data"}}
	val, _ := e.ToJson()
	assert.Equal(t, `{"control":{"extra":"extra","host":"host","ip_addresses":"ip_addresses","ts":99999999999999},"event_key":"product.update","payload":{"payload":"data"},"source":"core","target_document":{"payload":"data"},"target_id":"xxxxx","target_type":"product","trace_id":""}`, string(val))
}
