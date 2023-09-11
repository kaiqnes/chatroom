package di

import (
	"chatroom/internal/config"
	"chatroom/internal/controllers"
	"chatroom/internal/db"
	"chatroom/internal/logger"
	"chatroom/internal/middlewares"
	"github.com/gin-gonic/gin"
)

type di struct {
	database *db.Database
	engine   *gin.Engine
	log      logger.CustomLogger
}

func New(cfg *config.Config, database *db.Database, engine *gin.Engine, log logger.CustomLogger) (*di, error) {
	// Initialize dependencies
	return &di{
		database: database,
		engine:   engine,
		log:      log,
	}, nil
}

func (d *di) Inject() error {
	// Inject Middlewares
	authMiddleware := middlewares.NewAuthenticationMiddleware(d.log)

	// Inject Controllers
	health := controllers.NewHealthCheckController(d.engine)
	health.SetupEndpoints()

	auth := controllers.NewAuthController(d.engine, authMiddleware)
	auth.SetupEndpoints()

	chatroom := controllers.NewChatroomController(d.engine, authMiddleware)
	chatroom.SetupEndpoints()

	return nil
}
