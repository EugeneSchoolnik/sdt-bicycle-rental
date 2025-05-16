package main

import (
	"log/slog"
	"net/http"
	"sdt-bicycle-rental/internal/config"
	"sdt-bicycle-rental/internal/http-server/handlers/auth/register"
	"sdt-bicycle-rental/internal/repository/postgres"
	"sdt-bicycle-rental/internal/services"
	"sdt-bicycle-rental/lib/logger"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() { // go run ./cmd/bicycle-rental/main.go
	// Load the configuration
	cfg := config.MustLoad()
	// Initialize the logger
	log := logger.InitLogger(cfg.Env)
	log.Info("Logger initialized", slog.String("env", cfg.Env))

	// Initialize the database
	db, err := postgres.New(cfg.Postgres)
	if err != nil {
		log.Error("Failed to initialize database", slog.String("error", err.Error()))
		return
	}
	log.Info("Database initialized", slog.String("db_name", cfg.Postgres.DBName))

	// Initialize the HTTP server
	router := chi.NewRouter()

	// middleware
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	// routes
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello world"))
	})
	router.Route("/auth", func(r chi.Router) {
		authService := services.NewAuthService(postgres.NewUserRepository(db), log, cfg.JwtSecret)
		r.Post("/register", register.New(authService, log))
	})

	// Start the server
	httpAddr := ":" + strconv.Itoa(cfg.HTTPServer.Port)
	log.Info("starting server", slog.String("address", httpAddr))

	srv := &http.Server{
		Addr:         httpAddr,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server")
	}

	log.Error("server stopped")
}
