package config

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

// AppConfig contiene toda la configuración de la aplicación
type AppConfig struct {
	ServerPort     int    `json:"SERVER_PORT"`
	MQUser         string `json:"RABBIT_USER"`
	MQPassword     string `json:"RABBIT_PASSWORD"`
	MQHost         string `json:"RABBIT_HOST"`
	MQPort         int    `json:"RABBIT_PORT"`
	RabbitVHost    string `json:"RABBIT_VHOST"`
	RABBITMQ_QUEUE string `json:"RABBITMQ_QUEUE"`
	Environment    string `json:"ENVIRONMENT"`
}

// RabbitMQConfig contiene la configuración específica para RabbitMQ
type RabbitMQConfig struct {
	URL       string
	QueueName string
}

// LoadSecrets carga los secretos desde AWS Secrets Manager
func LoadSecrets() (*AppConfig, error) {
	secretName := os.Getenv("AWS_SECRET_NAME")
	if secretName == "" {
		return nil, fmt.Errorf("AWS_SECRET_NAME not set")
	}

	// Cargar configuración de AWS
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(os.Getenv("AWS_REGION")),
	)
	if err != nil {
		return nil, fmt.Errorf("unable to load SDK config: %v", err)
	}

	// Crear cliente de Secrets Manager
	svc := secretsmanager.NewFromConfig(cfg)

	// Obtener el secreto
	result, err := svc.GetSecretValue(context.TODO(), &secretsmanager.GetSecretValueInput{
		SecretId: &secretName,
	})
	if err != nil {
		return nil, fmt.Errorf("error getting secret: %v", err)
	}

	// Decodificar el secreto
	var config AppConfig
	if err := json.Unmarshal([]byte(*result.SecretString), &config); err != nil {
		return nil, fmt.Errorf("error unmarshaling secret: %v", err)
	}

	return &config, nil
}

// GetRabbitMQConfig construye y retorna la configuración de RabbitMQ
func (c *AppConfig) GetRabbitMQConfig() *RabbitMQConfig {
	url := fmt.Sprintf("amqp://%s:%s@%s:%d/%s",
		c.MQUser,
		c.MQPassword,
		c.MQHost,
		c.MQPort,
		c.RabbitVHost)

	return &RabbitMQConfig{
		URL:       url,
		QueueName: c.RABBITMQ_QUEUE,
	}
}
