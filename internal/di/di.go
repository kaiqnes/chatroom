package di

import (
	"chatroom/internal/clients"
	"chatroom/internal/config"
	"chatroom/internal/controllers"
	"chatroom/internal/db"
	"chatroom/internal/logger"
	"chatroom/internal/middlewares"
	"chatroom/internal/repositories"
	"chatroom/internal/use_cases"
	"fmt"
	"github.com/gin-gonic/gin"
	socketio "github.com/googollee/go-socket.io"
	"log"
	"net/http"
	"os"
)

type DI struct {
	cfg          *config.Config
	database     *db.Database
	httpServer   *gin.Engine
	socketServer *socketio.Server
	log          *logger.CustomLogger
}

func New(cfg *config.Config, database *db.Database, httpServer *gin.Engine, server *socketio.Server, log *logger.CustomLogger) (*DI, error) {
	// Initialize dependencies
	return &DI{
		cfg:          cfg,
		database:     database,
		httpServer:   httpServer,
		socketServer: server,
		log:          log,
	}, nil
}

func (d *DI) Inject() error {
	// Inject Middlewares
	authMiddleware := middlewares.NewAuthenticationMiddleware(d.cfg, d.log)

	// ================================================================================================================
	// Inject Clients
	stockBotClient := clients.NewStockBot(d.cfg)

	// ================================================================================================================
	// Inject Repositories
	chatRepository := repositories.NewChatRepository(d.database, d.log)
	userRepository := repositories.NewUserRepository(d.database, d.log)

	// Inject Use Cases
	sendMessageUseCase := use_cases.NewSendMessageUseCase(chatRepository, userRepository, d.log)
	authUseCase := use_cases.NewAuthUseCase(d.cfg.JwtKey, userRepository, d.log)

	// Inject Controllers
	health := controllers.NewHealthCheckController(d.httpServer)
	health.SetupEndpoints()

	authController := controllers.NewAuthController(d.httpServer, authUseCase)
	authController.SetupEndpoints()

	chatroom := controllers.NewChatroomController(d.httpServer, d.socketServer, authMiddleware, sendMessageUseCase)
	chatroom.SetupEndpoints()

	d.sockServer()

	// Inject Socket Routes
	d.httpServer.GET("/socket.io/*any", gin.WrapH(d.socketServer))
	d.httpServer.POST("/socket.io/*any", gin.WrapH(d.socketServer))

	dir, _ := os.Getwd()

	// Inject Front-End Files
	d.httpServer.LoadHTMLFiles(
		dir+"/pkg/front-end/auth/auth.html",
		dir+"/pkg/front-end/chatroom/chat_rooms.html")

	// Inject Front-End Routes
	d.httpServer.GET("/sign", func(c *gin.Context) {
		c.HTML(http.StatusOK, "auth.html", gin.H{
			"title": "auth",
		})
	})
	d.httpServer.GET("/chat", authMiddleware.ValidateToken(), func(c *gin.Context) {
		c.HTML(http.StatusOK, "chat_rooms.html", gin.H{
			"title": "Main website",
		})
	})

	return nil
}

func (d *DI) sockServer() {
	d.socketServer.OnConnect("/", func(s socketio.Conn) error {
		fmt.Println("socket connected:", s.ID())
		s.SetContext("")
		return nil
	})

	d.socketServer.OnError("/", func(s socketio.Conn, e error) {
		log.Println("meet error:", e)
	})

	d.socketServer.OnDisconnect("/", func(s socketio.Conn, reason string) {
		log.Println("closed", reason)
	})
}
