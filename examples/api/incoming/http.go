package incoming

import (
	"github.com/gin-gonic/gin"
	interfaces "github.com/shoplineapp/captin/v2/interfaces"
	models "github.com/shoplineapp/captin/v2/models"
	"net/http"
)

type HttpEventHandler struct {
	interfaces.IncomingHandler
	captin interfaces.CaptinInterface
}

func (h *HttpEventHandler) Setup(c interfaces.CaptinInterface) {
	h.captin = c
}

func (h HttpEventHandler) SetRoutes(router *gin.Engine) {
	router.GET("/", func(c *gin.Context) {
		c.String(200, "github.com/shoplineapp/captin/v2 aboard")
	})
	router.POST("/api/events", func(c *gin.Context) {
		h.HandleEventCreation(c)
	})
	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"code": "not_found"})
	})
}

func (h HttpEventHandler) HandleEventCreation(c *gin.Context) {
	var event models.IncomingEvent
	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, errors := h.captin.Execute(c, event)
	if len(errors) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error occurred when handling event"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"code": "created"})
}
