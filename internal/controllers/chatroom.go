package controllers

import (
	"net/http"

	"chatroom/internal/middlewares"
	"github.com/gin-gonic/gin"
)

type chatroomController struct {
	routes         *gin.Engine
	authMiddleware middlewares.AuthenticationMiddleware
}

func NewChatroomController(routes *gin.Engine, authMiddleware middlewares.AuthenticationMiddleware) Controller {
	return &chatroomController{routes: routes, authMiddleware: authMiddleware}
}

func (c *chatroomController) SetupEndpoints() {
	c.routes.GET("/home", c.authMiddleware.ValidateToken(), c.chatroom)
}

func (c *chatroomController) chatroom(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"chat": "room!",
	})
}
