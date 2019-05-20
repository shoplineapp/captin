package outgoing_filters

import (
	models "captin/internal/models"
	"encoding/json"
	"fmt"
	"github.com/robertkrimen/otto"
)

type ValidateFilter struct {
	Filter
	Event models.IncomingEvent
}

func (f ValidateFilter) Run(c models.Configuration) (bool, error) {
	payloadJson, _ := json.Marshal(f.Event.Payload)
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

func (f ValidateFilter) Applicable(c models.Configuration) bool {
	return (c.Validate) != ""
}
