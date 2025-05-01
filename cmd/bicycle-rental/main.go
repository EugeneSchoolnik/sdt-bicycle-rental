package main

import (
	"log/slog"
	"sdt-bicycle-rental/internal/config"
	"sdt-bicycle-rental/internal/repository/postgres"
	"sdt-bicycle-rental/lib/logger"
)

func main() { // go run ./cmd/bicycle-rental/main.go
	// Load the configuration
	cfg := config.MustLoad()
	// Initialize the logger
	log := logger.InitLogger(cfg.Env)
	log.Info("Logger initialized", slog.String("env", cfg.Env))

	// Initialize the database
	_, err := postgres.New(cfg.Postgres)
	if err != nil {
		log.Error("Failed to initialize database", slog.String("error", err.Error()))
		return
	}

	// Initialize the HTTP server

	// Start the server
}
