package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/graduate-work-mirea/data-collector-service/entity"
	"github.com/graduate-work-mirea/data-collector-service/internal/rabbitmq"
	"go.uber.org/zap"
)

// RabbitRepository handles sending data to RabbitMQ
type RabbitRepository struct {
	client    *rabbitmq.Client
	queueName string
	logger    *zap.SugaredLogger
}

// NewRabbitRepository creates a new instance of RabbitRepository
func NewRabbitRepository(client *rabbitmq.Client, queueName string, logger *zap.SugaredLogger) *RabbitRepository {
	return &RabbitRepository{
		client:    client,
		queueName: queueName,
		logger:    logger,
	}
}

// SendData marshals the product data to JSON and sends it to the RabbitMQ queue
func (r *RabbitRepository) SendData(ctx context.Context, data entity.ProductData) error {
	// Marshal data to JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal product data: %w", err)
	}

	// Send data to RabbitMQ
	err = r.client.PublishMessage(ctx, r.queueName, jsonData)
	if err != nil {
		return fmt.Errorf("failed to publish message to RabbitMQ: %w", err)
	}

	return nil
}

// SendBatch sends a batch of product data to RabbitMQ
func (r *RabbitRepository) SendBatch(ctx context.Context, dataList []entity.ProductData) error {
	r.logger.Infof("Sending batch of %d products to RabbitMQ", len(dataList))
	
	for i, data := range dataList {
		if err := r.SendData(ctx, data); err != nil {
			r.logger.Errorf("Failed to send product %d: %v", i, err)
			continue
		}
	}
	
	r.logger.Info("Batch processing completed")
	return nil
}
