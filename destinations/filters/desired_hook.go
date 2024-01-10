package destination_filters

import (
	"context"

	models "github.com/shoplineapp/captin/v2/models"
)

func isPresent(str string, list []string) bool {
	for _, v := range list {
		if v == str {
			return true
		}
	}
	return false
}

func stringList(list []interface{}) []string {
	sList := make([]string, len(list))
	for i, v := range list {
		sList[i] = v.(string)
	}
	return sList
}

// DesiredHookFilter - Filter destination if given event has desired destination
type DesiredHookFilter struct {
	DestinationFilterInterface
}

// Run - Get desired hooks in control and filter out exclusion
func (f DesiredHookFilter) Run(ctx context.Context, e models.IncomingEvent, d models.Destination) (bool, error) {
	hook := d.Config.GetName()
	list := e.Control["desired_hooks"]
	switch list.(type) {
	case []interface{}:
		list = stringList(list.([]interface{}))
		return isPresent(hook, list.([]string)), nil
	case []string:
		return isPresent(hook, list.([]string)), nil
	default:
		return false, nil
	}
}

// Applicable - Check if desired hooks is present
func (f DesiredHookFilter) Applicable(ctx context.Context, e models.IncomingEvent, d models.Destination) bool {
	return e.Control["desired_hooks"] != nil
}
