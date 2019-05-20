package incoming

import (
	interfaces "captin/interfaces"
	models "captin/internal/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

type HttpEventHandler struct {
	Handler
	captin interfaces.CaptinInterface
}

func (h *HttpEventHandler) Setup(c interfaces.CaptinInterface) {
	h.captin = c
}

func (h HttpEventHandler) SetRoutes(router *gin.Engine) {
	router.GET("/", func(c *gin.Context) {
		c.String(200, "Captin aboard")
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

	_, err := h.captin.Execute(event)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"code": "created"})
}
