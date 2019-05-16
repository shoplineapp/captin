package models_test

import (
	"github.com/stretchr/testify/assert"
	"testing"

	. "github.com/shoplineapp/captin/internal/models"
)

func tuples(args ...interface{}) []interface{} {
	return args
}

func TestEvaluatorRunValidate(t *testing.T) {
	payload := map[string]interface{}{
		"_id":  "xxxxxx",
		"type": "line",
	}
	assert.Equal(t, true, tuples(Evaluator{Payload: payload, Config: Configuration{Validate: "document.type == 'line'"}}.Run())[0])
	assert.Equal(t, false, tuples(Evaluator{Payload: payload, Config: Configuration{Validate: "document.type != 'line'"}}.Run())[0])

	// Error validate check
	assert.Equal(t, false, tuples(Evaluator{Payload: payload, Config: Configuration{Validate: "document.noneExist"}}.Run())[0])
	assert.Equal(t, false, tuples(Evaluator{Payload: payload, Config: Configuration{Validate: "invalidCall()"}}.Run())[0])
	assert.Equal(t, false, tuples(Evaluator{Payload: payload, Config: Configuration{Validate: "asd"}}.Run())[0])

	// Calling with nil payload should not affect test result
	assert.Equal(t, true, tuples(Evaluator{Payload: nil, Config: Configuration{Validate: "true"}}.Run())[0])
	assert.Equal(t, false, tuples(Evaluator{Payload: nil, Config: Configuration{Validate: "false"}}.Run())[0])

	// Calling without validate method should always return true
	assert.Equal(t, true, tuples(Evaluator{Payload: nil, Config: Configuration{}}.Run())[0])
}
