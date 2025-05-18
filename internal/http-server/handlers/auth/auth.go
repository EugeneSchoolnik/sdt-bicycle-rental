package auth

import (
	"log/slog"
	"sdt-bicycle-rental/internal/http-server/handlers/auth/login"
	"sdt-bicycle-rental/internal/http-server/handlers/auth/register"
	auth_service "sdt-bicycle-rental/internal/service/auth"

	"github.com/go-chi/chi/v5"
)

func AuthRoute(log *slog.Logger, userRepo auth_service.UserRepository, jwtSecret string) func(chi.Router) {
	return func(r chi.Router) {
		authService := auth_service.New(userRepo, log, jwtSecret)

		r.Post("/register", register.New(authService, log))
		r.Post("/login", login.New(authService, log))
	}
}
