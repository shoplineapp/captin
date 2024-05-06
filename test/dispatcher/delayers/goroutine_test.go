package dispatcher_delayers_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	. "github.com/shoplineapp/captin/v2/dispatcher/delayers"
	models "github.com/shoplineapp/captin/v2/models"
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
	delayer.Execute(context.Background(), evt, dest, callback)
	assert.Equal(t, true, isCallbackCalled)
}
