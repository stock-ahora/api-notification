package config

import (
	"time"

	"github.com/stock-ahora/api-notification/internal/messaging"
	"github.com/stock-ahora/api-notification/pkg/utils"
)

type Config struct {
	Server   ServerConfig
	RabbitMQ messaging.RabbitConfig
}

type ServerConfig struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

func Load() (*Config, error) {
	return &Config{
		Server: ServerConfig{
			Port:         utils.GetEnv("PORT", "8084"),
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
		},
		RabbitMQ: messaging.RabbitConfig{
			URL: utils.GetEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/"),
		},
	}, nil
}
