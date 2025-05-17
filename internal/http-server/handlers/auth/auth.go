package auth

import (
	"log/slog"
	"sdt-bicycle-rental/internal/http-server/handlers/auth/login"
	"sdt-bicycle-rental/internal/http-server/handlers/auth/register"
	"sdt-bicycle-rental/internal/services"

	"github.com/go-chi/chi/v5"
)

func AuthRoute(log *slog.Logger, userRepo services.UserRepository, jwtSecret string) func(chi.Router) {
	return func(r chi.Router) {
		authService := services.NewAuthService(userRepo, log, jwtSecret)

		r.Post("/register", register.New(authService, log))
		r.Post("/login", login.New(authService, log))
	}
}
