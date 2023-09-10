package di

import (
	"chatroom/internal/config"
	"chatroom/internal/db"
	"chatroom/internal/logger"
)

type di struct {
}

func New(cfg *config.Config, db *db.Database, logger logger.CustomLogger) (*di, error) {
	// Initialize dependencies
	return &di{}, nil
}

func (d *di) Inject() error {
	// Inject dependencies
	return nil
}
