package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	filters "example-with-filters/filters"
	middlewares "example-with-filters/middlewares"
	core "github.com/shoplineapp/captin/core"
	interfaces "github.com/shoplineapp/captin/interfaces"
	models "github.com/shoplineapp/captin/models"
)

func main() {
	// Load webhooks configuration
	pwd, _ := os.Getwd()
	configMapper := models.NewConfigurationMapperFromPath(filepath.Join(pwd, "hooks.json"))

	fmt.Println("* Captin loaded with hooks.json")
	captin := core.NewCaptin(*configMapper)
	captin.SetDestinationFilters([]interfaces.DestinationFilter{filters.TargetIdFilter{}})
	captin.SetDestinationMiddlewares([]interfaces.DestinationMiddleware{middlewares.LoggerMiddleware{}})

	jsonBytes := []byte(`{"event_key":"custom","source":"service_one","payload":{"_id":"A"},"target_id":"A"}`)
	event := models.IncomingEvent{}
	fmt.Println("= Sending event with Target ID A")
	json.Unmarshal(jsonBytes, &event)
	captin.Execute(event)

	jsonBytes = []byte(`{"event_key":"custom","source":"service_one","payload":{"_id":"A"},"target_id":"B"}`)
	event = models.IncomingEvent{}
	fmt.Println("= Sending event with Target ID B, expected to be filtered")
	json.Unmarshal(jsonBytes, &event)
	captin.Execute(event)
}
