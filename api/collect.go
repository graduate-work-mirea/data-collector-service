package api

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// handleCollect handles the collect endpoint
func (s *Server) handleCollect(c *gin.Context) {
	// Check if config is set
	if s.config.SourceURL == "" || s.config.APIKey == "" {
		s.logger.Println("Config not set")
		c.JSON(http.StatusBadRequest, Response{
			Status:    "error",
			ErrorCode: 1002,
			Error:     "Data source not configured. Call /collect-config first",
		})
		return
	}

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Create request
	req, err := http.NewRequest("GET", s.config.SourceURL, nil)
	if err != nil {
		s.logger.Printf("Error creating request: %v", err)
		c.JSON(http.StatusInternalServerError, Response{
			Status:    "error",
			ErrorCode: 1003,
			Error:     "Failed to create request",
		})
		return
	}

	// Add API key to request header
	req.Header.Add("Authorization", s.config.APIKey)

	// Send request
	s.logger.Printf("Sending request to %s", s.config.SourceURL)
	resp, err := client.Do(req)
	if err != nil {
		s.logger.Printf("Error sending request: %v", err)
		c.JSON(http.StatusInternalServerError, Response{
			Status:    "error",
			ErrorCode: 1004,
			Error:     "Failed to connect to data source",
		})
		return
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		s.logger.Printf("Error response from data source: %s", resp.Status)
		c.JSON(http.StatusInternalServerError, Response{
			Status:    "error",
			ErrorCode: 1005,
			Error:     "Data source returned an error: " + resp.Status,
		})
		return
	}

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		s.logger.Printf("Error reading response body: %v", err)
		c.JSON(http.StatusInternalServerError, Response{
			Status:    "error",
			ErrorCode: 1006,
			Error:     "Failed to read data from source",
		})
		return
	}

	// Parse JSON
	var rawData json.RawMessage
	if err := json.Unmarshal(body, &rawData); err != nil {
		s.logger.Printf("Error parsing JSON: %v", err)
		c.JSON(http.StatusInternalServerError, Response{
			Status:    "error",
			ErrorCode: 1007,
			Error:     "Failed to parse data from source",
		})
		return
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Save to database
	id, err := s.db.SaveRawData(ctx, rawData)
	if err != nil {
		s.logger.Printf("Error saving to database: %v", err)
		c.JSON(http.StatusInternalServerError, Response{
			Status:    "error",
			ErrorCode: 1008,
			Error:     "Failed to save data to database",
		})
		return
	}
	s.logger.Printf("Data saved to database with ID: %d", id)

	// Publish to RabbitMQ
	if err := s.queue.PublishRawData(ctx, rawData); err != nil {
		s.logger.Printf("Error publishing to RabbitMQ: %v", err)
		c.JSON(http.StatusInternalServerError, Response{
			Status:    "error",
			ErrorCode: 1009,
			Error:     "Failed to publish data to queue",
		})
		return
	}
	s.logger.Println("Data published to RabbitMQ")

	// Return success
	c.JSON(http.StatusOK, Response{
		Status:    "success",
		ErrorCode: 0,
	})
}
