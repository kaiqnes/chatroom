package servers

import (
	"chatroom/internal/config"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	socketio "github.com/googollee/go-socket.io"
	"log"
)

func New(cfg *config.Config) (*gin.Engine, *socketio.Server, error) {
	// Initialize servers
	httpServer := setupHTTPServer()
	socketServer := setupSocketServer()

	return httpServer, socketServer, nil
}

func setupSocketServer() *socketio.Server {
	socketServer := socketio.NewServer(nil)

	socketServer.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		return nil
	})

	socketServer.OnError("/", func(s socketio.Conn, e error) {
		log.Println("meet error:", e)
	})

	socketServer.OnDisconnect("/", func(s socketio.Conn, reason string) {
		log.Println("closed", reason)
	})

	return socketServer
}

func setupHTTPServer() *gin.Engine {
	httpServer := gin.Default()
	httpServer.Use(cors.Default())
	return httpServer
}
