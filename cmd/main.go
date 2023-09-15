package main

import (
	"chatroom/internal/di"
	"log"
)

func main() {
	// Inject dependencies
	dependencies, err := di.New()
	if err != nil {
		log.Fatalf("Error initializing dependencies: %s", err.Error())
	}
	dependencies.Inject()

	// Run servers
	dependencies.RunServers()
}
