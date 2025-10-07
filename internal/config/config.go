package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port     string
	RabbitMQ RabbitMQConfig
}

type RabbitMQConfig struct {
	URL       string
	QueueName string
}

func Load() (*Config, error) {
	// Cargar .env si existe (local)
	godotenv.Load()

	// Variables de entorno
	port := getEnv("PORT", "8084")
	queueName := getEnv("RABBITMQ_QUEUE", "notifications")

	// RabbitMQ desde variables de entorno
	mqHost := getEnv("MQ_HOST", "")
	mqPort := getEnv("MQ_PORT", "5672")
	mqUser := getEnv("MQ_USER", "")
	mqPassword := getEnv("MQ_PASSWORD", "")
	mqVhost := getEnv("RMQ_VHOST", "/")

	if mqHost == "" || mqUser == "" || mqPassword == "" {
		return nil, fmt.Errorf("missing RabbitMQ credentials: MQ_HOST, MQ_USER, or MQ_PASSWORD")
	}

	// Construir URL de RabbitMQ
	rabbitMQURL := fmt.Sprintf("amqp://%s:%s@%s:%s%s",
		mqUser,
		mqPassword,
		mqHost,
		mqPort,
		mqVhost,
	)

	return &Config{
		Port: port,
		RabbitMQ: RabbitMQConfig{
			URL:       rabbitMQURL,
			QueueName: queueName,
		},
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
