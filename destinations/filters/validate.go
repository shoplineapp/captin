package destination_filters

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/robertkrimen/otto"
	models "github.com/shoplineapp/captin/v2/models"
	log "github.com/sirupsen/logrus"
)

var vLogger = log.WithFields(log.Fields{"class": "ValidateFilter"})

var _ DestinationFilterInterface = ValidateFilter{}

type ValidateFilter struct {
}

func (f ValidateFilter) Run(ctx context.Context, e models.IncomingEvent, d models.Destination) (bool, error) {
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

func (f ValidateFilter) Applicable(ctx context.Context, e models.IncomingEvent, d models.Destination) bool {
	return (d.Config.GetValidate()) != ""
}
