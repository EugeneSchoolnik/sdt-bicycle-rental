package login

import (
	"errors"
	"log/slog"
	"net/http"
	"sdt-bicycle-rental/internal/models"
	"sdt-bicycle-rental/internal/services"
	"sdt-bicycle-rental/lib/logger/sl"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type Request struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Response struct {
	User  *models.User `json:"user,omitempty"`
	Token string       `json:"token,omitempty"`
	Error string       `json:"error,omitempty"`
}

//go:generate mockery --name=UserLoginer
type UserLoginer interface {
	Login(email, password string) (*models.User, string, error)
}

func New(s UserLoginer, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.auth.login.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		if err := render.DecodeJSON(r.Body, &req); err != nil {
			log.Error("failed to decode request body", sl.Err(err))

			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, Response{Error: "invalid input"})
			return
		}

		user, token, err := s.Login(req.Email, req.Password)
		if err != nil {
			if errors.Is(err, services.ErrInternalError) {
				w.WriteHeader(http.StatusInternalServerError)
			} else {
				w.WriteHeader(http.StatusBadRequest)
			}
			render.JSON(w, r, Response{Error: err.Error()})
			return
		}

		log.Info("user authorized", slog.Uint64("id", user.ID))

		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, Response{User: user, Token: token})
	}
}
