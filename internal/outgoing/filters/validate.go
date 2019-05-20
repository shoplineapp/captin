package outgoing_filters

import (
	models "captin/internal/models"
	"encoding/json"
	"fmt"
	"github.com/robertkrimen/otto"
)

type ValidateFilter struct {
	Filter
}

func (f ValidateFilter) Run(e models.IncomingEvent, c models.Configuration) (bool, error) {
	payloadJson, _ := json.Marshal(e.Payload)
	configJson, _ := json.Marshal(c)
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
		fmt.Printf("[ValidateFilter] Unable to parse result %s", err)
	}
	return valid, err
}

func (f ValidateFilter) Applicable(e models.IncomingEvent, c models.Configuration) bool {
	return (c.Validate) != ""
}
