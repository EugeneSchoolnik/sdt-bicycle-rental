package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"sdt-bicycle-rental/internal/http-server/handlers/auth"
	"sdt-bicycle-rental/internal/http-server/handlers/auth/register"
	"sdt-bicycle-rental/internal/models"
	"sdt-bicycle-rental/internal/repository/postgres"
	"sdt-bicycle-rental/internal/service"
	"sdt-bicycle-rental/lib/logger/handlers/slogdiscard"
	test_postgres "sdt-bicycle-rental/tests/util/db/postgres"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthHandler(t *testing.T) {
	db, cleanup := test_postgres.SetupTestDB(t)
	defer cleanup()

	test_postgres.ClearTable(t, db, "users")

	userRepo := postgres.NewUserRepository(db)

	r := chi.NewRouter()
	r.Route("/auth", auth.AuthRoute(slogdiscard.NewDiscardLogger(), userRepo, "secret"))

	t.Run("register", func(t *testing.T) {
		type resp struct {
			Code  int
			Error string
		}

		tests := []struct {
			name     string
			body     string
			wantResp resp
		}{
			{
				name:     "success",
				body:     `{"user":{"name":"John","lastname":"Doe","email":"john@example.com","phone":"123456","password":"12345678"}}`,
				wantResp: resp{Code: http.StatusCreated},
			},
			{
				name:     "invalid name",
				body:     `{"user":{"name":"","lastname":"Doe","email":"john@example.com","phone":"123456","password":"12345678"}}`,
				wantResp: resp{Code: http.StatusBadRequest, Error: "field Name is a required field"},
			},
			{
				name:     "invalid email",
				body:     `{"user":{"name":"John","lastname":"Doe","email":"example.com","phone":"123456","password":"12345678"}}`,
				wantResp: resp{Code: http.StatusBadRequest, Error: "field Email is not a valid email"},
			},
			{
				name:     "invalid password",
				body:     `{"user":{"name":"John","lastname":"Doe","email":"john@example.com","phone":"123456","password":"1234"}}`,
				wantResp: resp{Code: http.StatusBadRequest, Error: "field Password is not valid"},
			},
			{
				name:     "invalid name and email",
				body:     `{"user":{"name":"","lastname":"Doe","email":"example.com","phone":"123456","password":"12345678"}}`,
				wantResp: resp{Code: http.StatusBadRequest, Error: "field Name is a required field, field Email is not a valid email"},
			},
			{
				name:     "invalid body",
				body:     `not user`,
				wantResp: resp{Code: http.StatusBadRequest, Error: "invalid input"},
			},
			{
				name:     "user exists error",
				body:     `{"user":{"name":"John","lastname":"Doe","email":"john@example.com","phone":"123456","password":"12345678"}}`,
				wantResp: resp{Code: http.StatusConflict, Error: service.ErrUserAlreadyExists.Error()},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				req := httptest.NewRequest(http.MethodPost, "/auth/register", strings.NewReader(tt.body))
				req.Header.Set("Content-Type", "application/json")

				resp := httptest.NewRecorder()
				r.ServeHTTP(resp, req)

				assert.Equal(t, tt.wantResp.Code, resp.Code)

				if tt.wantResp.Error == "" {
					var response register.SuccessResponse
					require.NoError(t, render.DecodeJSON(resp.Body, &response))
					// database
					var user models.User
					err := db.First(&user, "email = ?", *response.User.Email).Error
					require.NoError(t, err)
					assert.Equal(t, *response.User.Name, *user.Name)

					// response check
					assert.NotEmpty(t, response.Token)
					return
				}

				var response register.ErrorResponse
				require.NoError(t, render.DecodeJSON(resp.Body, &response))
				assert.Equal(t, tt.wantResp.Error, response.Error)

			})
		}
	})
}
