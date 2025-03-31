package api

import (
	"log"
	"net/http"

	"github.com/dasha/data-collector-service/db"
	"github.com/dasha/data-collector-service/queue"
	"github.com/gin-gonic/gin"
)

// Server represents the API server
type Server struct {
	router   *gin.Engine
	logger   *log.Logger
	db       *db.PostgresDB
	queue    *queue.RabbitMQ
	config   *SourceConfig
}

// SourceConfig represents the configuration for the data source
type SourceConfig struct {
	SourceURL string `json:"source_url"`
	APIKey    string `json:"api_key"`
}

// Response represents the API response
type Response struct {
	Status    string `json:"status"`
	ErrorCode int    `json:"errorCode"`
	Error     string `json:"error,omitempty"`
}

// NewServer creates a new API server
func NewServer(logger *log.Logger, db *db.PostgresDB, queue *queue.RabbitMQ) *Server {
	router := gin.Default()
	server := &Server{
		router: router,
		logger: logger,
		db:     db,
		queue:  queue,
		config: &SourceConfig{},
	}

	// Register routes
	router.POST("/collect-config", server.handleCollectConfig)
	router.POST("/collect", server.handleCollect)

	return server
}

// Start starts the API server
func (s *Server) Start(addr string) error {
	return s.router.Run(addr)
}

// handleCollectConfig handles the collect-config endpoint
func (s *Server) handleCollectConfig(c *gin.Context) {
	var config SourceConfig
	if err := c.ShouldBindJSON(&config); err != nil {
		s.logger.Printf("Error binding JSON: %v", err)
		c.JSON(http.StatusBadRequest, Response{
			Status:    "error",
			ErrorCode: 1000,
			Error:     "Invalid request body",
		})
		return
	}

	// Validate config
	if config.SourceURL == "" || config.APIKey == "" {
		s.logger.Println("Invalid config: missing source_url or api_key")
		c.JSON(http.StatusBadRequest, Response{
			Status:    "error",
			ErrorCode: 1001,
			Error:     "Missing source_url or api_key",
		})
		return
	}

	// Store config
	s.config = &config
	s.logger.Printf("Config updated: source_url=%s", config.SourceURL)

	c.JSON(http.StatusOK, Response{
		Status:    "success",
		ErrorCode: 0,
	})
}
