package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"

	core "github.com/shoplineapp/captin/core"
	incoming "github.com/shoplineapp/captin/incoming"
	models "github.com/shoplineapp/captin/models"
)

func main() {
	fmt.Println("Starting in port:", os.Getenv("CAPTIN_PORT"))
	port := fmt.Sprintf(":%s", os.Getenv("CAPTIN_PORT"))

	// Load webhooks configuration
	pwd, _ := os.Getwd()
	path := os.Args[1:][0]
	absPath := filepath.Join(pwd, path)
	configMapper := models.NewConfigurationMapperFromPath(absPath)

	captin := core.NewCaptin(*configMapper)

	// Set up api server
	router := gin.Default()
	handler := incoming.HttpEventHandler{}
	handler.Setup(*captin)
	handler.SetRoutes(router)

	fmt.Printf("* Binding captin on 0.0.0.0%s\n", port)
	router.Run(port)
}
