package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type healthCheckController struct {
	routes *gin.Engine
}

func NewHealthCheckController(routes *gin.Engine) Controller {
	return &healthCheckController{routes: routes}
}

func (c *healthCheckController) SetupEndpoints() {
	c.routes.GET("/health", c.healthCheck)
}

func (c *healthCheckController) healthCheck(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"response": "https://www.youtube.com/watch?v=xos2MnVxe-c",
	})
}
