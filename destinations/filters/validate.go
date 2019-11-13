package destination_filters

import (
	"encoding/json"
	"fmt"
	"github.com/robertkrimen/otto"
	interfaces "github.com/shoplineapp/captin/interfaces"
	models "github.com/shoplineapp/captin/models"
	log "github.com/sirupsen/logrus"
)

var vLogger = log.WithFields(log.Fields{"class": "ValidateFilter"})

type ValidateFilter struct {
	interfaces.DestinationFilter
}

func (f ValidateFilter) Run(e models.IncomingEvent, d models.Destination) (bool, error) {
	payloadJson, _ := json.Marshal(e.Payload)
	configJson, _ := json.Marshal(d.Config)
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
		vLogger.WithFields(log.Fields{"error": err}).Error("Unable to parse result")
	}
	return valid, err
}

func (f ValidateFilter) Applicable(e models.IncomingEvent, d models.Destination) bool {
	return (d.Config.Validate) != ""
}
