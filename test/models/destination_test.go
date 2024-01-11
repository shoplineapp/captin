package models_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	. "github.com/shoplineapp/captin/v2/models"
)

func TestDestination_GetCallbackURL(t *testing.T) {
	var config Configuration
	var subject Destination

	config = Configuration{Name: "site_a", CallbackURL: "http://site-a.com/callback"}
	subject = Destination{Config: config}
	assert.Equal(t, subject.GetCallbackURL(), config.CallbackURL)

	overriden := "http://google.com"
	os.Setenv("HOOK_SITE_B_CALLBACK_URL", overriden)
	config = Configuration{Name: "site_b", CallbackURL: "http://site-b.com/callback"}
	subject = Destination{Config: config}
	assert.Equal(t, subject.GetCallbackURL(), overriden)
}

func TestDestination_GetDocumentStore(t *testing.T) {
	var config Configuration
	var subject Destination

	config = Configuration{Name: "callback_a"}
	subject = Destination{Config: config}
	assert.Equal(t, "default", subject.GetDocumentStore())

	config = Configuration{Name: "callback_b", DocumentStore: "store_b"}
	subject = Destination{Config: config}
	assert.Equal(t, subject.GetDocumentStore(), "store_b")

	overriden := "store_c"
	os.Setenv("HOOK_CALLBACK_C_DOCUMENT_STORE", overriden)
	config = Configuration{Name: "callback_c"}
	subject = Destination{Config: config}
	assert.Equal(t, overriden, subject.GetDocumentStore())
}

func TestDestination_RequireDelay(t *testing.T) {
	var config Configuration
	var subject Destination
	var event IncomingEvent

	// Hook without delay and event without control
	config = Configuration{Name: "callback_without_delay"}
	subject = Destination{Config: config}
	event = IncomingEvent{}
	assert.Equal(t, false, subject.RequireDelay(event))

	// Hook with delay and event without control
	config = Configuration{Name: "callback_without_delay", Delay: "10s"}
	subject = Destination{Config: config}
	event = IncomingEvent{}
	assert.Equal(t, true, subject.RequireDelay(event))

	// Hook with delay and event with control
	config = Configuration{Name: "callback_without_delay", Delay: "10s"}
	subject = Destination{Config: config}
	event = IncomingEvent{Control: map[string]interface{}{ "outstanding_delay_seconds": "10" }}
	assert.Equal(t, true, subject.RequireDelay(event))

	// Hook with delay and event with 0 outstanding delay in control
	config = Configuration{Name: "callback_without_delay", Delay: "10s"}
	subject = Destination{Config: config}
	event = IncomingEvent{Control: map[string]interface{}{ "outstanding_delay_seconds": "0" }}
	assert.Equal(t, false, subject.RequireDelay(event))
}

func TestDestination_GetRetryBackoffSeconds(t *testing.T) {
	var config Configuration
	var subject Destination
	var event IncomingEvent

	// Hook without retry backoff AND event without any retry record
	config = Configuration{Name: "no_retry_backoff_no_retry"}
	subject = Destination{Config: config}
	event = IncomingEvent{}
	assert.Equal(t, int64(DEFAULT_RETRY_BACKOFF_SECONDS), subject.GetRetryBackoffSeconds(event))

	// Hook without retry backoff AND event with retry record
	config = Configuration{Name: "no_retry_backoff_has_retry"}
	subject = Destination{Config: config}
	event = IncomingEvent{Control: map[string]interface{}{"retry_count": float64(2)}}
	assert.Equal(t, int64(DEFAULT_RETRY_BACKOFF_SECONDS), subject.GetRetryBackoffSeconds(event))

	// Hook with retry backoff AND event without any retry record
	config = Configuration{Name: "retry_backoff_no_retry", RetryBackoff: "5,30,200,600"}
	subject = Destination{Config: config}
	event = IncomingEvent{}
	assert.Equal(t, int64(5), subject.GetRetryBackoffSeconds(event))

	// Hook with retry backoff AND event with retry record
	config = Configuration{Name: "retry_backoff_no_retry", RetryBackoff: "5,30,200,600"}
	subject = Destination{Config: config}
	event = IncomingEvent{Control: map[string]interface{}{"retry_count": float64(4)}}
	assert.Equal(t, int64(600), subject.GetRetryBackoffSeconds(event))

	// Hook with retry backoff AND event with retry record
	config = Configuration{Name: "retry_backoff_has_retry", RetryBackoff: "5,30,200,600"}
	subject = Destination{Config: config}
	event = IncomingEvent{Control: map[string]interface{}{"retry_count": float64(2)}}
	assert.Equal(t, int64(200), subject.GetRetryBackoffSeconds(event))

	// Hook with retry backoff AND event with retries beyond config
	config = Configuration{Name: "retry_backoff_too_many_retry", RetryBackoff: "5,30,200,600"}
	subject = Destination{Config: config}
	event = IncomingEvent{Control: map[string]interface{}{"retry_count": float64(100)}}
	assert.Equal(t, int64(600), subject.GetRetryBackoffSeconds(event))
}
