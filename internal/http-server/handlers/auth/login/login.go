package login

import (
	"errors"
	"log/slog"
	"net/http"
	"sdt-bicycle-rental/internal/models"
	"sdt-bicycle-rental/internal/service"
	"sdt-bicycle-rental/lib/logger/sl"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type Request struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type SuccessResponse struct {
	User  *models.User `json:"user"`
	Token string       `json:"token"`
}
type ErrorResponse struct {
	Error string `json:"error"`
}

//go:generate mockery --name=UserLoginer
type UserLoginer interface {
	Login(email, password string) (*models.User, string, error)
}

// New returns login handler
//
//	@Summary      Login
//	@Description  login a user
//	@Tags         auth
//	@Accept       json
//	@Produce      json
//	@Param        request body 		Request true "User login data"
//	@Success      201  {object}   	SuccessResponse
//	@Failure      400  {object}		ErrorResponse
//	@Failure      409  {object}		ErrorResponse
//	@Failure      500  {object}		ErrorResponse
//	@Router       /auth/login [post]
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
			render.JSON(w, r, ErrorResponse{Error: "invalid input"})
			return
		}

		user, token, err := s.Login(req.Email, req.Password)
		if err != nil {
			if errors.Is(err, service.ErrInternalError) {
				w.WriteHeader(http.StatusInternalServerError)
			} else {
				w.WriteHeader(http.StatusBadRequest)
			}
			render.JSON(w, r, ErrorResponse{Error: err.Error()})
			return
		}

		log.Info("user authorized", slog.Uint64("id", user.ID))

		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, SuccessResponse{User: user, Token: token})
	}
}
