package di

import (
	"chatroom/internal/config"
	"chatroom/internal/controllers"
	"chatroom/internal/db"
	"chatroom/internal/logger"
	"chatroom/internal/middlewares"
	"chatroom/internal/repositories"
	"chatroom/internal/use_cases"
	"github.com/gin-gonic/gin"
)

type DI struct {
	database *db.Database
	engine   *gin.Engine
	log      *logger.CustomLogger
}

func New(cfg *config.Config, database *db.Database, engine *gin.Engine, log *logger.CustomLogger) (*DI, error) {
	// Initialize dependencies
	return &DI{
		database: database,
		engine:   engine,
		log:      log,
	}, nil
}

func (d *DI) Inject() error {
	// Inject Middlewares
	authMiddleware := middlewares.NewAuthenticationMiddleware(d.log)

	// Inject Repositories
	chatRepository := repositories.NewChatRepository(d.database, d.log)
	userRepository := repositories.NewUserRepository(d.database, d.log)

	// Inject Use Cases
	sendMessageUseCase := use_cases.NewSendMessageUseCase(chatRepository, userRepository, d.log)

	// Inject Controllers
	health := controllers.NewHealthCheckController(d.engine)
	health.SetupEndpoints()

	auth := controllers.NewAuthController(d.engine, authMiddleware)
	auth.SetupEndpoints()

	chatroom := controllers.NewChatroomController(d.engine, authMiddleware, sendMessageUseCase)
	chatroom.SetupEndpoints()

	return nil
}
