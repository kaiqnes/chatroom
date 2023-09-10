package config

import "time"

type Config struct {
	Env      string
	AppName  string
	LogLevel string
	Port     string
	Database Database
}

type Database struct {
	Host                  string
	Port                  string
	User                  string
	Password              string
	Name                  string
	MaxIdleConnections    int
	MaxOpenConnections    int
	MaxConnectionLifetime time.Duration
	MaxIdleTime           time.Duration
}

const defaultConfigPath = "./internal/config/config.yml"

func Load() (*Config, error) {
	// Read environment variables
	// Read config file
	return &Config{}, nil
}
