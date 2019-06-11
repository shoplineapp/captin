package models_test

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"

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
