package main

import (
	"fmt"
	"os"
	"path/filepath"

	filters "example-with-filters/filters"
	core "github.com/shoplineapp/captin/core"
	models "github.com/shoplineapp/captin/models"
)

func main() {
	// Load webhooks configuration
	pwd, _ := os.Getwd()
	configMapper := models.NewConfigurationMapperFromPath(filepath.Join(pwd, "hooks.json"))

	fmt.Println("* Captin loaded with hooks.json")
	captin := core.NewCaptin(*configMapper)
	captin.SetCustomFilters([]interfaces.CustomFilter{filters.CallbackUrlFilter})

	jsonBytes := []byte(`{"event_key":"custom","source":"service_one","payload":{"_id":"A"},"callback_url": "http://some-service/webhook"}`)
	event := models.IncomingEvent{}
	json.Unmarshal(jsonBytes, &event)
	success, _ = captin.Execute(event)
}
