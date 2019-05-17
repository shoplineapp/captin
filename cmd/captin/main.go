package main

import (
	"fmt"
	"os"
	"path/filepath"

	models "github.com/shoplineapp/captin/internal/models"
)

func main() {
	fmt.Println("* Starting captin")

	pwd, _ := os.Getwd()
	path := os.Args[1:][0]
	absPath := filepath.Join(pwd, path)

	mappedConfigs := models.NewConfigurationMapperFromPath(absPath)

	captin := Captin{ConfigMap: mappedConfigs}
	captin.Execute(IncomingEvent{Key: "product.update", Source: "core", Payload: map[string]interface{}{"field1": 1}})
}
