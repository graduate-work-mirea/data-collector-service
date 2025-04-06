package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/graduate-work-mirea/data-collector-service/assembly"
	"github.com/graduate-work-mirea/data-collector-service/config"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

// @title Data Collector Service
// @version 1.0
// @description Collect data from marketplaces
func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()

	if err := godotenv.Load(); err != nil {
		sugar.Warnf("Error loading .env file: %v", err)
	}

	cfg, err := config.New()
	if err != nil {
		sugar.Fatalf("Failed to load config: %v", err)
	}

	locator, err := assembly.NewServiceLocator(cfg, sugar)
	if err != nil {
		sugar.Fatalf("Failed to initialize service locator: %v", err)
	}
	defer locator.Close()
	
	// Create context with cancellation for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	
	// Set up signal handling for graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	
	// Start the scheduler in a goroutine
	sugar.Infof("Starting data collector scheduler with interval: %v", cfg.SchedulerInterval)
	go locator.DataProcessor.StartScheduler(ctx, cfg.SchedulerInterval)
	
	// Wait for termination signal
	sig := <-sigCh
	sugar.Infof("Received signal: %v, shutting down...", sig)
	
	// Cancel context to stop the scheduler
	cancel()
}

func startServerWithGracefulShutdown(srv *http.Server, logger *zap.SugaredLogger) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		logger.Infof("Starting server on port %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("Failed to start server: %v", err)
		}
	}()

	<-quit
	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatalf("Server forced to shutdown: %v", err)
	}

	logger.Info("Server exited properly")
}
