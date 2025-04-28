package logger

import (
	"log/slog"
	"os"
)

const (
	envLocal       = "local"
	envDevelopment = "dev"
	envProduction  = "prod"
)

func InitLogger(env string) *slog.Logger {
	// Initialize the logger based on the environment
	var logger *slog.Logger

	switch env {
	case envLocal:
		// Initialize local logger
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDevelopment:
		// Initialize development logger
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProduction:
		// Initialize production logger
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	default:
		panic("failed to initialize logger: unknown environment")
	}

	return logger
}
