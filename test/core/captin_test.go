package models_test

import (
	"github.com/stretchr/testify/assert"
	"testing"

	. "github.com/shoplineapp/captin/core"
	models "github.com/shoplineapp/captin/models"
)

func TestExecute(t *testing.T) {
	// When event is not given or is invalid
	var err error

	_, err = Captin{}.Execute(models.IncomingEvent{})

	if assert.Error(t, err, "invalid incoming event") {
		assert.IsType(t, err, &ExecutionError{})
	}
}
