package main

import (
	"log/slog"
	"sdt-bicycle-rental/internal/config"
	"sdt-bicycle-rental/lib/logger"
)

func main() { // CONFIG_PATH=config/local.yaml go run ./cmd/bicycle-rental/main.go
	// Load the configuration
	cfg := config.MustLoad()
	// Initialize the logger
	log := logger.InitLogger(cfg.Env)
	log.Info("Logger initialized", slog.String("env", cfg.Env))

	// Initialize the database

	// Initialize the HTTP server

	// Start the server
}
