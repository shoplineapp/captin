package models

import (
	"encoding/json"
)

type IncomingEvent struct {
	Key     string                 `json:"event_key"` // Required, The identifier of an event, usually form as PREFIX.MODEL.ACTION
	Source  string                 `json:"source"`    // Required, Event source from
	Payload map[string]interface{} `json:"payload"`   // Optional, custom payload / document from caller
	Control map[string]interface{} `json:"control"`   // Optional, custom control values from caller

	// Optional with payload, Captin will try to fetch the document from the default database
	TargetType     string                  `json:"target_type"`
	TargetId       string                  `json:"target_id"`
	TargetDocument map[string]interface{}  `json:"target_document,omitempty"`
}

func NewIncomingEvent(data []byte) IncomingEvent {
	e := IncomingEvent{}
	json.Unmarshal(data, &e)
	return e
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
