package models

import (
	"encoding/json"
)

type Event struct {
	Key     string                 `json:"event_key"` // Required, The identifier of an event, usually form as PREFIX.MODEL.ACTION
	Source  string                 `json:"source"`    // Required, Event source from
	Payload map[string]interface{} `json:"payload"`   // Optional, custom payload / document from caller

	// Optional with payload, Captin will try to fetch the document from the default database
	TargetType string `json:"target_type"`
	TargetId   string `json:"target_id"`
}

func NewEvent(data []byte) Event {
	event := Event{}
	json.Unmarshal(data, &event)
	return event
}

func (e Event) IsValid() bool {
	if e.Key == "" || e.Source == "" {
		return false
	}
	if e.TargetType == "" && e.TargetId == "" && len(e.Payload) == 0 {
		return false
	}
	return true
}
