package server

import (
	"chatroom/internal/config"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	socketio "github.com/googollee/go-socket.io"
)

func New(cfg *config.Config) (*gin.Engine, *socketio.Server, error) {
	// Initialize servers
	httpServer := setupHTTPServer()
	socketServer := setupSocketServer()

	return httpServer, socketServer, nil
}

func setupSocketServer() *socketio.Server {
	socketServer := socketio.NewServer(nil)

	return socketServer
}

func setupHTTPServer() *gin.Engine {
	httpServer := gin.Default()
	httpServer.Use(cors.Default())
	return httpServer
}
