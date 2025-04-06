package assembly

import (
	"github.com/graduate-work-mirea/data-collector-service/config"
	"github.com/graduate-work-mirea/data-collector-service/internal/rabbitmq"
	"github.com/graduate-work-mirea/data-collector-service/reader"
	"github.com/graduate-work-mirea/data-collector-service/repository"
	"github.com/graduate-work-mirea/data-collector-service/service"
	"go.uber.org/zap"
)

type ServiceLocator struct {
	Config        *config.Config
	RabbitClient  *rabbitmq.Client
	DataReader    *reader.DataReader
	RabbitRepo    *repository.RabbitRepository
	DataProcessor *service.DataProcessor
	Logger        *zap.SugaredLogger
}

func NewServiceLocator(cfg *config.Config, logger *zap.SugaredLogger) (*ServiceLocator, error) {
	rabbitClient, err := rabbitmq.NewClient(cfg.RabbitMQURL, logger)
	if err != nil {
		return nil, err
	}

	_, err = rabbitClient.DeclareQueue(cfg.DataQueueName)
	if err != nil {
		rabbitClient.Close()
		return nil, err
	}

	// Create data reader
	dataReader := reader.NewDataReader(cfg.DatasetPath, logger)

	// Create rabbit repository
	rabbitRepo := repository.NewRabbitRepository(rabbitClient, cfg.DataQueueName, logger)

	// Create data processor service
	dataProcessor := service.NewDataProcessor(dataReader, rabbitRepo, logger)

	return &ServiceLocator{
		Config:        cfg,
		RabbitClient:  rabbitClient,
		DataReader:    dataReader,
		RabbitRepo:    rabbitRepo,
		DataProcessor: dataProcessor,
		Logger:        logger,
	}, nil
}

func (l *ServiceLocator) Close() {
	if l.RabbitClient != nil {
		l.RabbitClient.Close()
	}
}
