package register

import (
	"errors"
	"log/slog"
	"net/http"
	"sdt-bicycle-rental/internal/models"
	"sdt-bicycle-rental/internal/repository/dto"
	"sdt-bicycle-rental/internal/service"
	"sdt-bicycle-rental/lib/logger/sl"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type Request struct {
	User dto.CreateUser `json:"user"`
}
type SuccessResponse struct {
	User  *models.User `json:"user"`
	Token string       `json:"token"`
}
type ErrorResponse struct {
	Error string `json:"error"`
}

//go:generate mockery --name=UserRegisterer
type UserRegisterer interface {
	Register(user *dto.CreateUser) (*models.User, string, error)
}

// New returns register handler
//
//	@Summary      Register
//	@Description  register a user
//	@Tags         auth
//	@Accept       json
//	@Produce      json
//	@Param        request body 		Request true "User registration data"
//	@Success      201  {object}   	SuccessResponse
//	@Failure      400  {object}		ErrorResponse
//	@Failure      409  {object}		ErrorResponse
//	@Failure      500  {object}		ErrorResponse
//	@Router       /auth/register [post]
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
			render.JSON(w, r, ErrorResponse{Error: "invalid input"})

			return
		}

		user, token, err := s.Register(&req.User)
		if err != nil {
			if errors.Is(err, service.ErrInternalError) {
				// internal error
				w.WriteHeader(http.StatusInternalServerError)
			} else if errors.Is(err, service.ErrUserAlreadyExists) {
				// user exists error
				w.WriteHeader(http.StatusConflict)
			} else {
				// other errors
				w.WriteHeader(http.StatusBadRequest)
			}
			render.JSON(w, r, ErrorResponse{Error: err.Error()})
			return
		}

		log.Info("user registered", slog.Uint64("id", user.ID))

		w.WriteHeader(http.StatusCreated)
		render.JSON(w, r, SuccessResponse{User: user, Token: token})
	}
}
