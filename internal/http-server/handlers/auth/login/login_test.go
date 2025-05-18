package login_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sdt-bicycle-rental/internal/http-server/handlers/auth/login"
	"sdt-bicycle-rental/internal/http-server/handlers/auth/login/mocks"
	"sdt-bicycle-rental/internal/models"
	"sdt-bicycle-rental/internal/services"
	"sdt-bicycle-rental/lib/logger/handlers/slogdiscard"
	. "sdt-bicycle-rental/lib/ptr"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoginHandler(t *testing.T) {
	type resp struct {
		Code  int
		Error string
	}

	cases := []struct {
		name      string
		email     string
		password  string
		resp      resp
		mockError error
		mockUser  *models.User
	}{
		{
			name:     "success",
			email:    "valid@email.com",
			password: "password",
			resp:     resp{Code: http.StatusOK},
			mockUser: &models.User{
				ID:       1,
				Email:    Ptr("valid@email.com"),
				Password: Ptr("hashed_password"),
			},
		},
		{
			name:      "invalid password",
			email:     "valid@email.com",
			password:  "1234",
			resp:      resp{Code: http.StatusBadRequest, Error: "field   is not valid"},
			mockError: errors.New("field   is not valid"),
		},
		{
			name:      "incorrect password",
			email:     "valid@email.com",
			password:  "12345678",
			resp:      resp{Code: http.StatusBadRequest, Error: services.ErrInvalidCredentials.Error()},
			mockError: services.ErrInvalidCredentials,
		},
		{
			name:      "internal error",
			email:     "valid@email.com",
			password:  "12345678",
			resp:      resp{Code: http.StatusInternalServerError, Error: services.ErrInternalError.Error()},
			mockError: services.ErrInternalError,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			userLoginerMock := mocks.NewUserLoginer(t)

			if tc.resp.Error == "" || tc.mockError != nil {
				mockCall := userLoginerMock.On("Login", tc.email, tc.password)
				mockCall.Return(tc.mockUser, "token", tc.mockError).Once()
			}

			handler := login.New(userLoginerMock, slogdiscard.NewDiscardLogger())

			input := fmt.Sprintf(`{"email": "%s", "password": "%s"}`, tc.email, tc.password)

			req, err := http.NewRequest(http.MethodPost, "/login", bytes.NewReader([]byte(input)))
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			require.Equal(t, rr.Code, tc.resp.Code)
			body := rr.Body.String()

			if rr.Code == http.StatusCreated {
				var resp login.SuccessResponse
				require.NoError(t, json.Unmarshal([]byte(body), &resp))
				assert.NotEqual(t, resp.User, nil)
				assert.NotEmpty(t, resp.Token)
				return
			}

			var resp login.ErrorResponse
			require.NoError(t, json.Unmarshal([]byte(body), &resp))
			require.Equal(t, tc.resp.Error, resp.Error)
		})
	}
}
