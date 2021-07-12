package models

import (
	"time"
	"strconv"
	"encoding/json"
	"github.com/google/uuid"

	interfaces "github.com/shoplineapp/captin/interfaces"
)

type IncomingEvent struct {
	interfaces.IncomingEventInterface

	TraceId string
	Key     string                 `json:"event_key"` // Required, The identifier of an event, usually form as PREFIX.MODEL.ACTION
	Source  string                 `json:"source"`    // Required, Event source from
	Payload map[string]interface{} `json:"payload"`   // Optional, custom payload / document from caller
	ThrottledPayloads []map[string]interface{} `json:"throttled_payloads,omitempty"`   // for response only
	Control map[string]interface{} `json:"control"`   // Optional, custom control values from caller

	// Optional with payload, Captin will try to fetch the document from the default database
	TargetType     string                 `json:"target_type"`
	TargetId       string                 `json:"target_id"`
	TargetDocument map[string]interface{} `json:"target_document,omitempty"`
	ThrottledDocuments []map[string]interface{} `json:"throttled_documents,omitempty"`   // for response only
}

func NewIncomingEvent(data []byte) IncomingEvent {
	e := IncomingEvent{}
	json.Unmarshal(data, &e)

	// Reuse trace ID if it's given for tracing retry with the same ID
	if e.TraceId == "" {
		e.TraceId = uuid.New().String()
	}
	return e
}

func (e IncomingEvent) GetTraceInfo() map[string]interface{} {
	return map[string]interface{}{
		"trace_id": e.TraceId,
		"key":     e.Key,
		"source":  e.Source,
		"type":    e.TargetType,
		"id":      e.TargetId,
	}
}

func (e IncomingEvent) GetControl() map[string]interface{} {
	return e.Control
}

func (e IncomingEvent) GetOutstandingDelaySeconds() time.Duration {
	if e.Control["outstanding_delay_seconds"] == nil {
		return time.Duration(-1) * time.Second
	}

	s, err := strconv.Atoi(e.Control["outstanding_delay_seconds"].(string))
	if err != nil || s <= 0 {
		s = 0
	}

	return time.Duration(s) * time.Second
}

func (e IncomingEvent) IsValid() bool {
	if e.Key == "" || e.Source == "" {
		return false
	}
	if e.TargetType == "" && e.TargetId == "" && len(e.Payload) == 0 {
		return false
	}
	return true
}
