package models_test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"

	. "captin/internal"
	models "captin/internal/models"
)

func TestExecute(t *testing.T) {
	// When event is not given or is invalid
	var err error

	_, err = Captin{}.Execute(models.IncomingEvent{})
	fmt.Println(reflect.TypeOf(err))

	if assert.Error(t, err, "invalid incoming event") {
		assert.IsType(t, err, &ExecutionError{})
	}
}
