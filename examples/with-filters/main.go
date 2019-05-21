package main

import (
	"fmt"
	"os"
	"path/filepath"

	core "github.com/shoplineapp/captin/core"
	incoming "github.com/shoplineapp/captin/incoming"
)

func main() {
	// Load webhooks configuration
	pwd, _ := os.Getwd()
	configMapper := core.NewConfigurationMapperFromPath(filepath.Join(pwd, "hooks.json"))

	fmt.Println("* Captin loaded with hooks.json")
	captin := core.NewCaptin(*configMapper)

	jsonStr := []byte(`{"event_key":"custom","source":"service_one","payload":{"_id":"A"}}`)
	captin.Execute()
}
