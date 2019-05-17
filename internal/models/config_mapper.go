package models

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// ConfigurationMapper - Action to configuration mapper
type ConfigurationMapper struct {
	ActionMap map[string][]Configuration
}

// NewConfigurationMapper - Create ConfigurationMapper with array of Configurations
func NewConfigurationMapper(configs []Configuration) *ConfigurationMapper {
	result := ConfigurationMapper{
		ActionMap: make(map[string][]Configuration),
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

	configs := []Configuration{}
	json.Unmarshal(data, &configs)

	return NewConfigurationMapper(configs)
}
