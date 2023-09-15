package config

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
	"strings"
	"time"
)

type Config struct {
	Env                     string        `mapstructure:"env"`
	AppName                 string        `mapstructure:"app_name"`
	Port                    string        `mapstructure:"app_port"`
	LogLevel                string        `mapstructure:"log_level"`
	JwtKey                  string        `mapstructure:"jwt_key"`
	DBHost                  string        `mapstructure:"db_host"`
	DBPort                  string        `mapstructure:"db_port"`
	DBUser                  string        `mapstructure:"db_user"`
	DBPassword              string        `mapstructure:"db_password"`
	DBName                  string        `mapstructure:"db_name"`
	DBSSLMode               string        `mapstructure:"db_ssl_mode"`
	DBMaxIdleConnections    int           `mapstructure:"db_max_idle_conns"`
	DBMaxOpenConnections    int           `mapstructure:"db_max_open_conns"`
	DBMaxConnectionLifetime time.Duration `mapstructure:"db_max_life_time"`
	DBMaxIdleTime           time.Duration `mapstructure:"db_max_idle_time"`
	StockBotTemplateURL     string        `mapstructure:"stock_bot_template_url"`
	RabbitMQHost            string        `mapstructure:"rabbitmq_host"`
	RabbitMQHostTemplate    string        `mapstructure:"rabbitmq_host_template"`
	RabbitMQPort            string        `mapstructure:"rabbitmq_port"`
	RabbitMQUser            string        `mapstructure:"rabbitmq_user"`
	RabbitMQPassword        string        `mapstructure:"rabbitmq_password"`
	RabbitMQQueue           string        `mapstructure:"rabbitmq_queue"`
}

const defaultConfigPath = "./internal/config/"

func Load() (*Config, error) {
	// Read environment variables
	env := getEnv()

	// Read config file
	viper.SetConfigName("config_" + env)
	viper.AddConfigPath(defaultConfigPath)
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	for _, k := range viper.AllKeys() {
		value := viper.GetString(k)
		if strings.HasPrefix(value, "${") && strings.HasSuffix(value, "}") {
			viper.Set(k, getEnvOrPanic(strings.TrimSuffix(strings.TrimPrefix(value, "${"), "}")))
		}
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		fmt.Println(err)
		return nil, err
	}

	cfg.Env = env

	return &cfg, nil
}

func getEnv() string {
	env := os.Getenv("ENV")
	if env == "" {
		env = "dev"
	}
	return env
}

func getEnvOrPanic(env string) string {
	res := os.Getenv(env)
	if len(res) == 0 {
		panic("Mandatory env variable not found:" + env)
	}
	return res
}
