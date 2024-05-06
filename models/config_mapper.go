package models

import (
	"encoding/json"
	"io/ioutil"

	log "github.com/sirupsen/logrus"
	interfaces "github.com/shoplineapp/captin/v2/interfaces"
)

var cmLogger = log.WithFields(log.Fields{"class": "ConfigurationMapper"})

// ConfigurationMapper - Action to configuration mapper
type ConfigurationMapper struct {
	ActionMap map[string][]interfaces.ConfigurationInterface
}

// NewConfigurationMapper - Create ConfigurationMapper with array of Configurations
func NewConfigurationMapper(configs []interfaces.ConfigurationInterface) *ConfigurationMapper {
	result := ConfigurationMapper{
		ActionMap: make(map[string][]interfaces.ConfigurationInterface),
	}
	for _, config := range configs {
		for _, action := range config.GetActions() {
			list := result.ActionMap[action]
			list = append(list, config)
			result.ActionMap[action] = list
		}
	}
	return &result
}

// NewConfigurationMapperFromPath - Read Configuration from path
func NewConfigurationMapperFromPath(path string) *ConfigurationMapper {
	pathLogger := cmLogger.WithFields(log.Fields{"path": path})
	data, err := ioutil.ReadFile(path)

	if err != nil {
		pathLogger.Error("Failed to load file")
		panic(err)
	}

	raw := []Configuration{}
	jsonErr := json.Unmarshal(data, &raw)
	if jsonErr != nil {
		pathLogger.Error("Invalid configuration file format")
		panic(jsonErr)
	}

	configs := []interfaces.ConfigurationInterface{}
	for _, c := range raw {
		configs = append(configs, c)
	}
	return NewConfigurationMapper(configs)
}

func (cm ConfigurationMapper) ConfigsForKey(eventKey string) []interfaces.ConfigurationInterface {
	return cm.ActionMap[eventKey]
}
