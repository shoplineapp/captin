package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"

	incoming "redis_example/incoming"
	stores "redis_example/stores"

	core "github.com/shoplineapp/captin/v2/core"
	"github.com/shoplineapp/captin/v2/models"
)

func main() {
	fmt.Println("Starting in port:", getEnv("CAPTIN_PORT", "3000"))
	port := fmt.Sprintf(":%s", getEnv("CAPTIN_PORT", "3000"))

	// Load webhooks configuration
	pwd, _ := os.Getwd()
	path := os.Args[1:][0]
	absPath := filepath.Join(pwd, path)
	configMapper := models.NewConfigurationMapperFromPath(absPath)
	captin := core.NewCaptin(*configMapper)

	// Set up redis store
	redisHost := getEnv("CAPTIN_REDIS_HOST", "localhost")
	redisPort := getEnv("CAPTIN_REDIS_PORT", "6379")

	store := stores.NewRedisStore(fmt.Sprintf("%s:%s", redisHost, redisPort))
	captin.SetStore(*store)
	fmt.Printf("[Main] %+v\n", store)

	// Set up api server
	router := gin.Default()
	handler := incoming.HttpEventHandler{}
	handler.Setup(*captin)
	handler.SetRoutes(router)

	fmt.Printf("* Binding captin on 0.0.0.0%s\n", port)
	router.Run(port)
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return defaultValue
	}
	return value
}
