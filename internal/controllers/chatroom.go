package controllers

import (
	"fmt"
	"net/http"
	"time"

	"chatroom/internal/domain"

	"github.com/gin-gonic/gin"
	socketio "github.com/googollee/go-socket.io"
)

type chatroomController struct {
	httpServer         *gin.Engine
	socketServer       *socketio.Server
	authMiddleware     domain.AuthenticationMiddleware
	sendMessageUseCase domain.SendMessageUseCase
}

func NewChatroomController(httpServer *gin.Engine, socketServer *socketio.Server,
	authMiddleware domain.AuthenticationMiddleware, sendMessageUseCase domain.SendMessageUseCase) domain.ChatController {
	return &chatroomController{httpServer: httpServer, socketServer: socketServer,
		authMiddleware: authMiddleware, sendMessageUseCase: sendMessageUseCase}
}

func (c *chatroomController) SetupEndpoints() {
	c.socketServer.OnEvent("/", "chat message", c.OnMessage)
}

func (c *chatroomController) SendMessageFromBot(req domain.MessageRequestDto) {
	resp := domain.MessageResponseDto{
		RoomID:    req.RoomID,
		Message:   req.Message,
		Username:  req.Username,
		Timestamp: time.Now(),
	}
	_ = c.socketServer.BroadcastToNamespace("/", "chat message", resp)
}

func (c *chatroomController) OnMessage(socket socketio.Conn, req domain.MessageRequestDto) {
	socket.SetContext(req)

	// Call use case
	resp, err := c.sendMessageUseCase.SendMessage(req.Username, req.RoomID, req.Message)
	if err != nil {
		// Handle error
		fmt.Printf("error to persist message: %+v\n", err)
		return
	}

	// Broadcast to all clients
	for _, r := range resp {
		_ = c.socketServer.BroadcastToNamespace("/", "chat message", r)
	}
}

func (c *chatroomController) ListMessages(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"list": "messages!",
	})
}
