package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	RabbitMQURL       string
	DataQueueName     string
	DatasetPath       string
	SchedulerInterval time.Duration
}

func New() (*Config, error) {
	rabbitMQURL := os.Getenv("RABBITMQ_URL")
	if rabbitMQURL == "" {
		rabbitMQURL = "amqp://guest:guest@localhost:5672/"
	}

	dataQueueName := os.Getenv("DATA_QUEUE_NAME")
	if dataQueueName == "" {
		dataQueueName = "marketplace_data"
	}

	datasetPath := os.Getenv("DATASET_PATH")
	if datasetPath == "" {
		datasetPath = "./dataset/ozon.csv"
	}

	schedulerIntervalStr := os.Getenv("SCHEDULER_INTERVAL_HOURS")
	schedulerInterval := 24 * time.Hour // Default: once per day
	if schedulerIntervalStr != "" {
		interval, err := strconv.Atoi(schedulerIntervalStr)
		if err == nil && interval > 0 {
			schedulerInterval = time.Duration(interval) * time.Hour
		}
	}

	return &Config{
		RabbitMQURL:       rabbitMQURL,
		DataQueueName:     dataQueueName,
		DatasetPath:       datasetPath,
		SchedulerInterval: schedulerInterval,
	}, nil
}
