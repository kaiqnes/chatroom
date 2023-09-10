package main

import (
	"log"

	"chatroom/internal/config"
	"chatroom/internal/db"
	"chatroom/internal/di"
	"chatroom/internal/logger"
	"chatroom/internal/server"
)

func main() {
	// Read config file
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Error loading config: %s", err.Error())
	}

	// Initialize logger
	customLogger, err := logger.New(cfg.Env)
	if err != nil {
		log.Fatalf("Error initializing logger: %s", err.Error())
	}

	// Initialize database
	database, err := db.New(cfg)
	if err != nil {
		log.Fatalf("Error initializing database: %s", err.Error())
	}

	// Initialize server
	srv, err := server.New(cfg)
	if err != nil {
		log.Fatalf("Error initializing server: %s", err.Error())
	}

	// Inject dependencies
	dependencies, err := di.New(cfg, database, customLogger)
	if err != nil {
		log.Fatalf("Error initializing dependencies: %s", err.Error())
	}

	err = dependencies.Inject()
	if err != nil {
		log.Fatalf("Error injecting dependencies: %s", err.Error())
	}

	// Start server
	err = srv.Run(":" + cfg.Port)
	if err != nil {
		log.Fatalf("Error running server: %s", err.Error())
	}
}
