package models

import (
	"encoding/json"
	"fmt"
	"github.com/robertkrimen/otto"
)

type Evaluator struct {
	Payload map[string]interface{}
	Config  Configuration
}

func (evaluator Evaluator) Run() (bool, error) {
	if (evaluator.Config.Validate) == "" {
		return true, nil
	}
	payloadJson, _ := json.Marshal(evaluator.Payload)
	configJson, _ := json.Marshal(evaluator.Config)
	template := fmt.Sprintf(
		`(function() {
			var document = %s || {};
			var config = %s || {};
			return !!(eval(config.validate));
		})()`,
		string(payloadJson),
		string(configJson))

	runtime := otto.New()
	result, err := runtime.Run(template)

	valid, errToB := result.ToBoolean()
	if errToB != nil {
		err = errToB
	}
	if err != nil {
		fmt.Printf("[Evaluator] Unable to parse result %s", err)
	}
	return valid, err
}
