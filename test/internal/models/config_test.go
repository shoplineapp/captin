package models_test

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	. "github.com/shoplineapp/captin/internal/models"
)

func TestAbs(t *testing.T) {
	fixture, err := ioutil.ReadFile("fixtures/config.json")
	if err != nil {
		panic(err)
	}
	subject := Configuration{}
	json.Unmarshal(fixture, &subject)

	if subject.CallbackURL != "http://callback_url/sync" {
		t.Errorf("subject.CallbackURL = %s; want http://callback_url/sync", subject.CallbackURL)
	}

	if subject.ConfigID != "1234" {
		t.Errorf("subject.ConfigID = %s; want 1234", subject.ConfigID)
	}

	if subject.Validate != "function(obj) { return !!obj.wapos_id }" {
		t.Errorf("subject.Validate = %s; want function(obj) { return !!obj.wapos_id }", subject.Validate)
	}

	if subject.Throttle != "500ms" {
		t.Errorf("subject.Throttle = %s; want 500ms", subject.Throttle)
	}

	if len(subject.Actions) != 6 {
		t.Errorf("len(subject.Actions) = %d; want 6", 6)
	}

	if subject.Source != "core-api" {
		t.Errorf("subject.Source = %s; want core-api", subject.Source)
	}

	if subject.Name != "sync_service" {
		t.Errorf("subject.Name = %s; want sync_service", subject.Name)
	}

	if subject.IncludeDocument != false {
		t.Errorf("subject.IncludeDocument = true; want false")
	}
}
