package dispatcher_delayers_test

import (
	"github.com/stretchr/testify/assert"
	"testing"

	. "github.com/shoplineapp/captin/dispatcher/delayers"
	models "github.com/shoplineapp/captin/models"
)

func TestGoroutineDelayer_Execute(t *testing.T) {
	isCallbackCalled := false
	callback := func() {
		isCallbackCalled = true
	}

	evt := models.IncomingEvent{}
	dest := models.Destination{
		Config: models.Configuration{
			Delay: "1ms",
		},
	}
	delayer := GoroutineDelayer{}
	delayer.Execute(evt, dest, callback)
	assert.Equal(t, true, isCallbackCalled)
}
