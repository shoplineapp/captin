package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	models "github.com/shoplineapp/captin/internal/models"
)

func main() {
	fmt.Println("* Starting captin")
	fmt.Println(os.Args[1:])

	pwd, _ := os.Getwd()
	path := os.Args[1:][0]
	absPath := filepath.Join(pwd, path)

	data, err := ioutil.ReadFile(absPath)

	if err != nil {
		fmt.Println("[Configuration] Failed to load file")
		panic(err)
	}

	configs := []models.Configuration{}
	json.Unmarshal(data, &configs)

	mappedConfigs := models.NewConfigurationMapper(configs)

	// TODO: Use Mapped Configuration for webhook logics
	fmt.Printf("%+v\n", mappedConfigs.ActionMap)
}
