package register

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
	User *models.User `json:"user"`
}

type Response struct {
	User  *models.User `json:"user,omitempty"`
	Token string       `json:"token,omitempty"`
	Error string       `json:"error,omitempty"`
}

//go:generate mockery --name=UserRegisterer
type UserRegisterer interface {
	Register(user *models.User) (*models.User, string, error)
}

func New(s UserRegisterer, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.auth.register.New"

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

		user, token, err := s.Register(req.User)
		if err != nil {
			if errors.Is(err, services.ErrInternalError) {
				// internal error
				w.WriteHeader(http.StatusInternalServerError)
			} else if errors.Is(err, services.ErrUserAlreadyExists) {
				// user exists error
				w.WriteHeader(http.StatusConflict)
			} else {
				// other errors
				w.WriteHeader(http.StatusBadRequest)
			}
			render.JSON(w, r, Response{Error: err.Error()})
			return
		}

		log.Info("user registered", slog.Uint64("id", user.ID))

		w.WriteHeader(http.StatusCreated)
		render.JSON(w, r, Response{User: user, Token: token})
	}
}
