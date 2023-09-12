package controllers

import (
	"chatroom/internal/domain"
	"net/http"

	"chatroom/internal/middlewares"
	"github.com/gin-gonic/gin"
)

type chatroomController struct {
	routes             *gin.Engine
	authMiddleware     middlewares.AuthenticationMiddleware
	sendMessageUseCase domain.SendMessageUseCase
}

func NewChatroomController(routes *gin.Engine, authMiddleware middlewares.AuthenticationMiddleware,
	sendMessageUseCase domain.SendMessageUseCase) Controller {
	return &chatroomController{routes: routes, authMiddleware: authMiddleware, sendMessageUseCase: sendMessageUseCase}
}

func (c *chatroomController) SetupEndpoints() {
	c.routes.GET("/list-messages", c.authMiddleware.ValidateToken(), c.ListMessages)
	c.routes.POST("/send-message", c.authMiddleware.ValidateToken(), c.SendMessage)
}

func (c *chatroomController) ListMessages(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"list": "messages!",
	})
}

func (c *chatroomController) SendMessage(context *gin.Context) {
	// Get user ID from context
	mockUserID := "45952e2e-0913-40d8-8e86-15d9d6c2eccc"

	// Get room ID from context
	mockRoomID := "746bd2ff-886b-422e-8d3b-21a42fa7b213"

	// Get body from request
	var body messageInputDto
	err := context.ShouldBindJSON(&body)
	if err != nil {
		// Handle error
		context.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Call use case
	err = c.sendMessageUseCase.SendMessage(mockUserID, mockRoomID, body.Message)
	if err != nil {
		// Handle error
		context.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}

	// Return response
	context.JSON(http.StatusNoContent, gin.H{})
}
