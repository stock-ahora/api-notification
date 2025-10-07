package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
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

type RabbitMQSecret struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	Vhost    string `json:"vhost"`
}

func Load() (*Config, error) {
	// Cargar .env si existe (local)
	godotenv.Load()

	port := getEnv("PORT", "8084")
	queueName := getEnv("RABBITMQ_QUEUE", "notifications")
	awsRegion := getEnv("AWS_REGION", "us-east-1")
	secretName := getEnv("AWS_SECRET_NAME", "prod/rabbitmq/credentials")

	// Obtener credenciales de AWS Secret Manager
	rabbitMQURL, err := getRabbitMQCredentials(awsRegion, secretName)
	if err != nil {
		return nil, fmt.Errorf("error getting RabbitMQ credentials: %w", err)
	}

	return &Config{
		Port: port,
		RabbitMQ: RabbitMQConfig{
			URL:       rabbitMQURL,
			QueueName: queueName,
		},
	}, nil
}

func getRabbitMQCredentials(region, secretName string) (string, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		return "", err
	}

	svc := secretsmanager.New(sess)
	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretName),
	}

	result, err := svc.GetSecretValue(input)
	if err != nil {
		return "", err
	}

	var secret RabbitMQSecret
	if err := json.Unmarshal([]byte(*result.SecretString), &secret); err != nil {
		return "", err
	}

	// Construir URL de RabbitMQ
	vhost := secret.Vhost
	if vhost == "" {
		vhost = "/"
	}

	url := fmt.Sprintf("amqp://%s:%s@%s:%s%s",
		secret.Username,
		secret.Password,
		secret.Host,
		secret.Port,
		vhost,
	)

	return url, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
