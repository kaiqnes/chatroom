package db

import "chatroom/internal/config"

type Database struct {
}

func New(cfg *config.Config) (*Database, error) {
	// Initialize database connection
	return &Database{}, nil
}
