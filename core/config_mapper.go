package core

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	models "github.com/shoplineapp/captin/internal/models"
)

// ConfigurationMapper - Action to configuration mapper
type ConfigurationMapper struct {
	ActionMap map[string][]models.Configuration
}

// NewConfigurationMapper - Create ConfigurationMapper with array of Configurations
func NewConfigurationMapper(configs []models.Configuration) *ConfigurationMapper {
	result := ConfigurationMapper{
		ActionMap: make(map[string][]models.Configuration),
	}
	for _, config := range configs {
		for _, action := range config.Actions {
			list := result.ActionMap[action]
			list = append(list, config)
			result.ActionMap[action] = list
		}
	}

	return &result
}

// NewConfigurationMapperFromPath - Read Configuration from path
func NewConfigurationMapperFromPath(path string) *ConfigurationMapper {
	data, err := ioutil.ReadFile(path)

	if err != nil {
		fmt.Println("[Configuration] Failed to load file")
		panic(err)
	}

	configs := []models.Configuration{}
	json.Unmarshal(data, &configs)

	return NewConfigurationMapper(configs)
}

func (cm ConfigurationMapper) ConfigsForKey(eventKey string) []models.Configuration {
	return cm.ActionMap[eventKey]
}
