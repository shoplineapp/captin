package models_test

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"

	. "github.com/shoplineapp/captin/internal/models"
)

func TesttNewEven(t *testing.T) {
	event_key := "product.update"
	source := "core"
	payload := map[string]interface{}{
		"_id":  "xxxxx",
		"type": "product",
	}
	target_type := "product"
	target_id := "xxxxx"

	data, _ := json.Marshal(map[string]interface{}{
		"event_key":   event_key,
		"source":      source,
		"payload":     payload,
		"target_type": target_type,
		"target_id":   target_id,
	})
	event := NewEvent(data)
	assert.Equal(t, event_key, event.Key)
	assert.Equal(t, source, event.Source)
	assert.Equal(t, payload, event.Payload)
	assert.Equal(t, target_type, event.TargetType)
	assert.Equal(t, target_id, event.TargetId)
}

func TestIsValid(t *testing.T) {
	assert.Equal(t, true, Event{Key: "product.update", Source: "core", Payload: map[string]interface{}{}, TargetType: "product", TargetId: "xxxxx"}.IsValid())

	// Test missing required attributes
	assert.Equal(t, false, Event{Source: "core", Payload: map[string]interface{}{}, TargetType: "product", TargetId: "xxxxx"}.IsValid())
	assert.Equal(t, false, Event{Key: "product.update", Payload: map[string]interface{}{}, TargetType: "product", TargetId: "xxxxx"}.IsValid())
	assert.Equal(t, false, Event{Key: "product.update", Source: "core"}.IsValid())

	// Test optional attributes
	assert.Equal(t, true, Event{Key: "product.update", Source: "core", TargetType: "product", TargetId: "xxxxx"}.IsValid())
	assert.Equal(t, true, Event{Key: "product.update", Source: "core", Payload: map[string]interface{}{"_id": "xxxxx"}}.IsValid())
}
