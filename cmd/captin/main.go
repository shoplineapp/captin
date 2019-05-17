package main

import (
	"fmt"
	"os"
	"path/filepath"

	internal "github.com/shoplineapp/captin/internal"
	models "github.com/shoplineapp/captin/internal/models"
)

func main() {
	fmt.Println("* Starting captin")

	pwd, _ := os.Getwd()
	path := os.Args[1:][0]
	absPath := filepath.Join(pwd, path)

	mappedConfigs := models.NewConfigurationMapperFromPath(absPath)

	captin := internal.Captin{ConfigMap: *mappedConfigs}
	captin.Execute(models.IncomingEvent{Key: "product.update", Source: "core", Payload: map[string]interface{}{"field1": 1}})
}
