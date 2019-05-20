package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"

	internal "captin/internal"
	models "captin/internal/models"
)

func main() {
	fmt.Println("* Starting captin (Press ENTER to quit)")

	pwd, _ := os.Getwd()
	path := os.Args[1:][0]
	absPath := filepath.Join(pwd, path)

	mappedConfigs := models.NewConfigurationMapperFromPath(absPath)

	captin := internal.Captin{ConfigMap: *mappedConfigs}
	captin.Execute(models.IncomingEvent{
		Key:        "product.update",
		Source:     "core",
		Payload:    map[string]interface{}{"field1": 1},
		TargetType: "Product",
		TargetId:   "product_id",
	})

	bufio.NewReader(os.Stdin).ReadBytes('\n')
}
