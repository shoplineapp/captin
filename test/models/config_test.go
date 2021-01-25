package models_test

import (
	"encoding/json"
	"io/ioutil"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	. "github.com/shoplineapp/captin/models"
)

func TestDecodeConfigurationJson(t *testing.T) {
	fixture, err := ioutil.ReadFile("fixtures/config.json")
	if err != nil {
		panic(err)
	}
	subject := Configuration{}
	json.Unmarshal(fixture, &subject)

	assert.Equal(t, subject.CallbackURL, "http://callback_url/sync")
	assert.Equal(t, subject.ConfigID, "1234")
	assert.Equal(t, subject.Validate, "function(obj) { return !!obj.wapos_id }")
	assert.Equal(t, subject.Throttle, "500ms")
	assert.Equal(t, len(subject.Actions), 6)
	assert.Equal(t, subject.Source, "core-api")
	assert.Equal(t, subject.Name, "sync_service")
	assert.Equal(t, subject.Sender, "mock")
	assert.Equal(t, subject.IncludeDocument, false)
}

func TestGetThrottle(t *testing.T) {
	subject := Configuration{}
	subject.Throttle = "50ms"
	assert.Equal(t, subject.GetThrottleValue(), time.Duration(50)*time.Millisecond)

	subject.Throttle = "50s"
	assert.Equal(t, subject.GetThrottleValue(), time.Duration(50)*time.Second)
}

func TestGetDelay(t *testing.T) {
	subject := Configuration{}
	subject.Delay = "500ms"
	assert.Equal(t, subject.GetDelayValue(), time.Duration(500)*time.Millisecond)

	subject.Delay = "1s"
	assert.Equal(t, subject.GetDelayValue(), time.Duration(1)*time.Second)
}

func TestConfiguration_GetDocumentStore(t *testing.T) {
	subject := Configuration{}
	subject.DocumentStore = "another"
	assert.Equal(t, subject.DocumentStore, "another")
}
