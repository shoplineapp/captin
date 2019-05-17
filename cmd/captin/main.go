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

	// TODO: Use Mapped Configuration for webhook logics
	fmt.Printf("%+v\n", mappedConfigs.ActionMap)
}
