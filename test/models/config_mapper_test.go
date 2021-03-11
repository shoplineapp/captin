package models_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	. "github.com/shoplineapp/captin/models"
	interfaces "github.com/shoplineapp/captin/interfaces"
)

func setup() []interfaces.ConfigurationInterface {
	result := []interfaces.ConfigurationInterface{}
	for i := 0; i < 3; i++ {
		switch i {
		case 0:
			val := Configuration{
				Name:    "0",
				Actions: []string{"action:0", "action:2"},
			}
			result = append(result, val)
		case 1:
			val := Configuration{
				Name:    "1",
				Actions: []string{"action:1", "action:0", "action:3"},
			}
			result = append(result, val)
		case 2:
			val := Configuration{
				Name:    "2",
				Actions: []string{"action:0", "action:2", "action:3"},
			}
			result = append(result, val)
		}
	}

	return result
}

func getNames(configs []interfaces.ConfigurationInterface) []string {
	result := []string{}
	for _, config := range configs {
		result = append(result, config.GetName())
	}
	return result
}

func TestMapActionToConfig(t *testing.T) {
	subject := NewConfigurationMapper(setup())

	action := subject.ActionMap["action:0"]
	assert.Equal(t, 3, len(action))
	names := getNames(action)
	assert.Contains(t, names, "0")
	assert.Contains(t, names, "1")
	assert.Contains(t, names, "2")

	action = subject.ActionMap["action:1"]
	assert.Equal(t, 1, len(action))
	names = getNames(action)
	assert.Contains(t, names, "1")

	action = subject.ActionMap["action:2"]
	assert.Equal(t, 2, len(action))
	names = getNames(action)
	assert.Contains(t, names, "0")
	assert.Contains(t, names, "2")

	action = subject.ActionMap["action:3"]
	assert.Equal(t, 2, len(action))
	names = getNames(action)
	assert.Contains(t, names, "1")
	assert.Contains(t, names, "2")
}

func TestReadLocalFile(t *testing.T) {
	pwd, _ := os.Getwd()
	absPath := filepath.Join(pwd, "fixtures/config_list.json")

	subject := NewConfigurationMapperFromPath(absPath)
	action := subject.ActionMap["product.update"]
	assert.Equal(t, 2, len(action))
	names := getNames(action)
	assert.Contains(t, names, "sync_service")
	assert.Contains(t, names, "sync_service2")
}

func TestConfigsForKey(t *testing.T) {
	subject := NewConfigurationMapper(setup())
	action := subject.ConfigsForKey("action:0")
	assert.Equal(t, 3, len(action))
	names := getNames(action)
	assert.Contains(t, names, "0")
	assert.Contains(t, names, "1")
}
