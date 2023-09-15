package di

import (
	"chatroom/internal/clients/message_broker"
	"chatroom/internal/clients/rest"
	"chatroom/internal/config"
	"chatroom/internal/controllers"
	"chatroom/internal/db"
	"chatroom/internal/domain"
	"chatroom/internal/logger"
	"chatroom/internal/middlewares"
	"chatroom/internal/repositories"
	"chatroom/internal/servers"
	"chatroom/internal/use_cases"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	socketio "github.com/googollee/go-socket.io"
)

type DI struct {
	cfg          *config.Config
	database     *db.Database
	httpServer   *gin.Engine
	socketServer *socketio.Server
	messageQueue domain.MessageQueue
	customLogger logger.CustomLogger
}

func New() (*DI, error) {
	// Initialize dependencies
	// Read config file
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Error loading config: %s", err.Error())
	}

	// Initialize message broker
	mq, err := message_broker.NewMQ(cfg)
	if err != nil {
		log.Fatalf("Error initializing message_broker: %s", err.Error())
	}

	// Initialize logger
	customLogger, err := logger.New(cfg.Env)
	if err != nil {
		log.Fatalf("Error initializing logger: %s", err.Error())
	}

	// Initialize database
	database, err := db.New(cfg, &customLogger)
	if err != nil {
		log.Fatalf("Error initializing database: %s", err.Error())
	}

	// Initialize server
	httpServer, socketServer, err := servers.New()
	if err != nil {
		log.Fatalf("Error initializing server: %s", err.Error())
	}

	return &DI{
		cfg:          cfg,
		database:     database,
		httpServer:   httpServer,
		socketServer: socketServer,
		messageQueue: mq,
		customLogger: customLogger,
	}, nil
}

func (d *DI) Inject() {
	// Inject Middlewares
	authMiddleware := middlewares.NewAuthenticationMiddleware(d.cfg, d.customLogger)

	// ================================================================================================================
	// Inject Clients
	stockBotClient := rest.NewStockBot(d.cfg)

	// ================================================================================================================
	// Inject Repositories
	chatRepository := repositories.NewChatRepository(d.database, d.customLogger)
	userRepository := repositories.NewUserRepository(d.database, d.customLogger)

	// ================================================================================================================
	// Inject Use Cases
	sendMessageUseCase := use_cases.NewSendMessageUseCase(chatRepository, userRepository, stockBotClient, d.customLogger)
	authUseCase := use_cases.NewAuthUseCase(d.cfg.JwtKey, userRepository, d.customLogger)

	// ================================================================================================================
	// Inject Controllers
	health := controllers.NewHealthCheckController(d.httpServer)
	health.SetupEndpoints()

	authController := controllers.NewAuthController(d.httpServer, authUseCase)
	authController.SetupEndpoints()

	chatroom := controllers.NewChatroomController(d.httpServer, d.socketServer, authMiddleware, sendMessageUseCase)
	chatroom.SetupEndpoints()

	// ================================================================================================================
	// Inject Socket Routes
	d.httpServer.GET("/socket.io/*any", gin.WrapH(d.socketServer))
	d.httpServer.POST("/socket.io/*any", gin.WrapH(d.socketServer))

	dir, _ := os.Getwd()

	// ================================================================================================================
	// Inject Front-End Files
	d.httpServer.LoadHTMLFiles(
		dir+"/pkg/front-end/auth/auth.html",
		dir+"/pkg/front-end/chatroom/chat_rooms.html")

	// ================================================================================================================
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
}

func (d *DI) RunServers() {
	// Start message broker
	go func() {
		err := d.messageQueue.Listen()
		if err != nil {
			log.Fatalf("Error listening on message_broker: %s", err.Error())
		}
	}()

	// Start Socket Server
	go func() {
		if err := d.socketServer.Serve(); err != nil {
			log.Fatalf("socketio listen error: %s\n", err)
		}
		defer d.socketServer.Close()
	}()

	// Start HTTP Server
	err := d.httpServer.Run(":" + d.cfg.Port)
	if err != nil {
		log.Fatalf("Error running server: %s", err.Error())
	}
}
