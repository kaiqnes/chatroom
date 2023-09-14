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
	database, err := db.New(cfg, &customLogger)
	if err != nil {
		log.Fatalf("Error initializing database: %s", err.Error())
	}

	// Initialize server
	httpServer, socketServer, err := server.New(cfg)
	if err != nil {
		log.Fatalf("Error initializing server: %s", err.Error())
	}

	// Inject dependencies
	dependencies, err := di.New(cfg, database, httpServer, socketServer, &customLogger)
	if err != nil {
		log.Fatalf("Error initializing dependencies: %s", err.Error())
	}
	dependencies.Inject()

	// Start Socket Server
	go func() {
		if err := socketServer.Serve(); err != nil {
			log.Fatalf("socketio listen error: %s\n", err)
		}
	}()
	defer socketServer.Close()

	// Start HTTP Server
	err = httpServer.Run(":" + cfg.Port)
	if err != nil {
		log.Fatalf("Error running server: %s", err.Error())
	}
}
