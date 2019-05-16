package models_test

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"

	. "github.com/shoplineapp/captin/internal/models"
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
	assert.Equal(t, subject.IncludeDocument, false)
}
