package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Controller struct for API endpoints
type Controller struct{}

// NewController creates a new Controller instance
func NewController() *Controller {
	return &Controller{}
}

// RegisterRoutes registers all API endpoints to the router
func (c *Controller) RegisterRoutes(router *gin.Engine) {
	api := router.Group("/api")
	{
		api.GET("/ping", c.Ping)
		// Přidejte další endpointy zde
	}
}

// Ping je ukázkový endpoint
func (c *Controller) Ping(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "pong"})
}
