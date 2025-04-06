package service

import (
	"context"
	"fmt"
	"time"

	"github.com/graduate-work-mirea/data-collector-service/repository"
	"go.uber.org/zap"
)

// DataProcessor handles processing of the dataset data
type DataProcessor struct {
	reader     *repository.DataReader
	repository *repository.RabbitRepository
	logger     *zap.SugaredLogger
}

// NewDataProcessor creates a new instance of DataProcessor
func NewDataProcessor(
	reader *repository.DataReader,
	repository *repository.RabbitRepository,
	logger *zap.SugaredLogger,
) *DataProcessor {
	return &DataProcessor{
		reader:     reader,
		repository: repository,
		logger:     logger,
	}
}

// ProcessData reads data from the dataset, processes it and sends to RabbitMQ
func (p *DataProcessor) ProcessData(ctx context.Context) error {
	p.logger.Info("Starting data processing")

	// Read data from the dataset
	products, err := p.reader.ReadData()
	if err != nil {
		return fmt.Errorf("failed to read data: %w", err)
	}

	p.logger.Infof("Processing %d products", len(products))

	// Process and send data to RabbitMQ
	err = p.repository.SendBatch(ctx, products)
	if err != nil {
		return fmt.Errorf("failed to send data to RabbitMQ: %w", err)
	}

	p.logger.Info("Data processing completed successfully")
	return nil
}

// StartScheduler starts a scheduler that processes data at specified intervals
func (p *DataProcessor) StartScheduler(ctx context.Context, interval time.Duration) {
	p.logger.Infof("Starting scheduler with interval: %v", interval)

	// Process data immediately on start
	if err := p.ProcessData(ctx); err != nil {
		p.logger.Errorf("Initial data processing failed: %v", err)
	}

	// Set up ticker for periodic processing
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			p.logger.Info("Scheduler triggered data processing")
			if err := p.ProcessData(ctx); err != nil {
				p.logger.Errorf("Scheduled data processing failed: %v", err)
			}
		case <-ctx.Done():
			p.logger.Info("Scheduler stopped")
			return
		}
	}
}
