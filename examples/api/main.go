package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"os"
	"path/filepath"

	core "github.com/shoplineapp/captin/core"
	incoming "github.com/shoplineapp/captin/incoming"
	models "github.com/shoplineapp/captin/models"
)

func main() {
	port := ":8080"

	pwd, _ := os.Getwd()
	configMapper := models.NewConfigurationMapperFromPath(filepath.Join(pwd, "hooks.json"))

	captin := core.NewCaptin(*configMapper)

	// Set up api server
	router := gin.Default()
	handler := incoming.HttpEventHandler{}
	handler.Setup(*captin)
	handler.SetRoutes(router)

	router.POST("/callback", func(c *gin.Context) {
		fmt.Println("Webhook callback received")
		c.String(200, "Received")
	})

	fmt.Printf("* Binding captin on 0.0.0.0%s\n", port)
	router.Run(port)
}
