package main

import (
	"fmt"
	// . "github.com/shoplineapp/captin/internal/helpers"
	. "github.com/shoplineapp/captin/internal/models"
	"os"
	"path/filepath"
	// "fmt"
)

func main() {
	fmt.Println("* Starting captin")

	pwd, _ := os.Getwd()
	path := "configs/hooks.json"
	absPath := filepath.Join(pwd, path)
	mappedConfigs := NewConfigurationMapperFromPath(absPath)

	// TODO: Use Mapped Configuration for webhook logics
	fmt.Printf("%+v\n", mappedConfigs.ActionMap["product.update"])

	captin := Captin{}
	captin.Wake(*mappedConfigs)
	captin.Execute(IncomingEvent{Key: "product.update", Source: "core", Payload: map[string]interface{}{"wtf": 1}})
}
