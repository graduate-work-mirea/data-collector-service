package config

import (
	"os"
)

// Config holds the application configuration
type Config struct {
	RabbitMQURL string
	PostgresDSN string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	// Default values
	rabbitMQURL := "amqp://guest:guest@rabbitmq:5672/"
	postgresDSN := "postgres://user:password@postgres:5432/db"

	// Override with environment variables if provided
	if envRabbitMQ := os.Getenv("RABBITMQ_URL"); envRabbitMQ != "" {
		rabbitMQURL = envRabbitMQ
	}

	if envPostgres := os.Getenv("POSTGRES_DSN"); envPostgres != "" {
		postgresDSN = envPostgres
	}

	return &Config{
		RabbitMQURL: rabbitMQURL,
		PostgresDSN: postgresDSN,
	}, nil
}
