package models

import (
	"encoding/json"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"

	interfaces "github.com/shoplineapp/captin/interfaces"
)

type IncomingEvent struct {
	interfaces.IncomingEventInterface

	TraceId           string                   `json:"trace_id"`
	Key               string                   `json:"event_key"`                    // Required, The identifier of an event, usually form as PREFIX.MODEL.ACTION
	Source            string                   `json:"source"`                       // Required, Event source from
	Payload           map[string]interface{}   `json:"payload"`                      // Optional, custom payload / document from caller
	ThrottledPayloads []map[string]interface{} `json:"throttled_payloads,omitempty"` // for response only
	Control           map[string]interface{}   `json:"control"`                      // Optional, custom control values from caller

	// Optional with payload, Captin will try to fetch the document from the default database
	TargetType         string                   `json:"target_type"`
	TargetId           string                   `json:"target_id"`
	TargetDocument     map[string]interface{}   `json:"target_document,omitempty"`
	ThrottledDocuments []map[string]interface{} `json:"throttled_documents,omitempty"` // for response only
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
		"key":      e.Key,
		"source":   e.Source,
		"type":     e.TargetType,
		"id":       e.TargetId,
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

func (e IncomingEvent) String() string {
	val, err := e.MarshalJSON()
	if err != nil {
		panic(err)
	}
	return string(val)
}

// default to exclude certain fields to reduce size of the resulted json, to get full payload, use ToJson
func (e IncomingEvent) MarshalJSON() ([]byte, error) {
	val := e.GetTraceInfo()

	if e.Control != nil {
		val["control"] = map[string]interface{}{
			"ts":           e.Control["ts"],
			"host":         e.Control["host"],
			"ip_addresses": e.Control["ip_addresses"],
		}
	}

	return json.Marshal(val)
}

func (e IncomingEvent) ToJson() ([]byte, error) {
	return json.Marshal(e.ToMap())
}

func (e IncomingEvent) ToMap() map[string]interface{} {
	out := make(map[string]interface{})

	v := reflect.ValueOf(e)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		fi := t.Field(i)
		if tagValue := fi.Tag.Get("json"); tagValue != "" {
			parts := strings.Split(tagValue, ",")
			if parts[0] == "-" || parts[0] == "" {
				continue
			}
			if len(parts) > 1 && parts[1] == "omitempty" {
				if v.Field(i).IsZero() {
					continue
				}
			}
			out[parts[0]] = v.Field(i).Interface()
		}
	}

	return out
}
