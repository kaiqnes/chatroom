package controllers

import (
	"fmt"
	"net/http"
	"time"

	"chatroom/internal/domain"
	"chatroom/internal/middlewares"
	"github.com/gin-gonic/gin"
	socketio "github.com/googollee/go-socket.io"
)

type chatroomController struct {
	httpServer         *gin.Engine
	socketServer       *socketio.Server
	authMiddleware     middlewares.AuthenticationMiddleware
	sendMessageUseCase domain.SendMessageUseCase
}

func NewChatroomController(httpServer *gin.Engine, socketServer *socketio.Server,
	authMiddleware middlewares.AuthenticationMiddleware, sendMessageUseCase domain.SendMessageUseCase) Controller {
	return &chatroomController{httpServer: httpServer, socketServer: socketServer,
		authMiddleware: authMiddleware, sendMessageUseCase: sendMessageUseCase}
}

func (c *chatroomController) SetupEndpoints() {
	c.socketServer.OnEvent("/", "chat message", c.OnMessage)
}

func (c *chatroomController) OnMessage(socket socketio.Conn, req messageInputDto) {
	//socket.SetContext(req)

	// Call use case
	err := c.sendMessageUseCase.SendMessage(req.Username, req.RoomID, req.Message)
	if err != nil {
		// Handle error
		fmt.Printf("error to persist message: %+v\n", err)
		return
	}

	resp := MessageOutputDto{
		RoomID:    req.RoomID,
		Message:   req.Message,
		Username:  req.Username,
		Timestamp: time.Now(),
	}

	_ = c.socketServer.BroadcastToNamespace("/", "chat message", resp)
}

func (c *chatroomController) ListMessages(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"list": "messages!",
	})
}
