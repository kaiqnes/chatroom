package server

import (
	"chatroom/internal/config"
	"github.com/gin-gonic/gin"
)

func New(cfg *config.Config) (*gin.Engine, error) {
	// Initialize server
	engine := gin.Default()

	return engine, nil
}
