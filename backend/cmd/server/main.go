package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"real-time-forum/internal/config"
	"real-time-forum/packages/logger"
)

func main() {
	appLogger := logger.NewLogger(os.Stdout, logger.InfoLevel)

	config, err := config.LoadConfig()
	if err != nil {
		appLogger.Fatal("Failed to load config", "error", err)
	}

	if err := config.Validate(); err != nil {
		appLogger.Fatal("Invalid configuration", "error", err)
	}

	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", config.Server.Port),
		ReadTimeout:  config.Server.ReadTimeout,
		WriteTimeout: config.Server.WriteTimeout,
		IdleTimeout:  config.Server.IdleTimeout,
	}

	go func() {
		appLogger.Info("Starting server", "port", config.Server.Port, "environment", config.Environment)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			appLogger.Fatal("Failed to start server", "error", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	appLogger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		appLogger.Error("Failed to shutdown server gracefully", "error", err)
	}

	appLogger.Info("Server stopped successfully")
}
