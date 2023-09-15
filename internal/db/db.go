package db

import (
	"fmt"
	"time"

	"chatroom/internal/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database struct {
	DBInstance *gorm.DB
}

func New(cfg *config.Config) (*Database, error) {
	// Initialize database connection with GORM
	fmt.Println("[db] Initializing database connection")
	dsn := getDSN(cfg)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		NowFunc: time.Now().UTC,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}

	return &Database{
		DBInstance: db,
	}, nil
}

func getDSN(cfg *config.Config) string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=UTC",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort, cfg.DBSSLMode)
}
