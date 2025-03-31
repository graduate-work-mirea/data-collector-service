package main

import (
	"log"
	"os"

	"github.com/dasha/data-collector-service/api"
	"github.com/dasha/data-collector-service/config"
	"github.com/dasha/data-collector-service/db"
	"github.com/dasha/data-collector-service/queue"
)

func main() {
	// Initialize logger
	logger := log.New(os.Stdout, "[DATA-COLLECTOR] ", log.LstdFlags)
	logger.Println("Starting Data Collector Service...")

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize database connection
	database, err := db.NewPostgresDB(cfg.PostgresDSN)
	if err != nil {
		logger.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	// Initialize RabbitMQ connection
	rabbitMQ, err := queue.NewRabbitMQ(cfg.RabbitMQURL)
	if err != nil {
		logger.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer rabbitMQ.Close()

	// Initialize and start API server
	server := api.NewServer(logger, database, rabbitMQ)
	if err := server.Start(":8080"); err != nil {
		logger.Fatalf("Failed to start server: %v", err)
	}
}
