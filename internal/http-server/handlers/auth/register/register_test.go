package register_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sdt-bicycle-rental/internal/http-server/handlers/auth/register"
	"sdt-bicycle-rental/internal/http-server/handlers/auth/register/mocks"
	"sdt-bicycle-rental/internal/repository/dto"
	"sdt-bicycle-rental/internal/service"
	"sdt-bicycle-rental/lib/logger/handlers/slogdiscard"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegisterHandler(t *testing.T) {
	type user struct {
		Name     string
		Lastname string
		Email    string
		Phone    string
		Password string
	}
	type resp struct {
		Code  int
		Error string
	}

	cases := []struct {
		name      string
		user      user
		resp      resp
		mockError error
	}{
		{
			name: "success",
			user: user{Name: "John", Lastname: "Doe", Email: "valid@email.com", Phone: "1234567890", Password: "superpass"},
			resp: resp{Code: http.StatusCreated},
		},
		{
			name:      "invalid email",
			user:      user{Name: "John", Lastname: "Doe", Email: "invalidemail.com", Phone: "1234567890", Password: "superpass"},
			resp:      resp{Code: http.StatusBadRequest, Error: "field Email is not a valid email"},
			mockError: errors.New("field Email is not a valid email"),
		},
		{
			name:      "user exists",
			user:      user{Name: "John", Lastname: "Doe", Email: "valid@email.com", Phone: "1234567890", Password: "superpass"},
			resp:      resp{Code: http.StatusConflict, Error: service.ErrUserAlreadyExists.Error()},
			mockError: service.ErrUserAlreadyExists,
		},
		{
			name:      "internal error",
			user:      user{Name: "John", Lastname: "Doe", Email: "valid@email.com", Phone: "1234567890", Password: "superpass"},
			resp:      resp{Code: http.StatusInternalServerError, Error: service.ErrInternalError.Error()},
			mockError: service.ErrInternalError,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			userRegistererMock := mocks.NewUserRegisterer(t)

			var userModel dto.CreateUser
			inputUser, _ := json.Marshal(tc.user)
			json.Unmarshal(inputUser, &userModel)

			if tc.resp.Error == "" || tc.mockError != nil {
				mockCall := userRegistererMock.On("Register", &userModel)
				mockCall.Return(userModel.Model(), "token", tc.mockError).Once()
			}

			handler := register.New(userRegistererMock, slogdiscard.NewDiscardLogger())

			input := fmt.Sprintf(`{"user": %s}`, inputUser)

			req, err := http.NewRequest(http.MethodPost, "/register", bytes.NewReader([]byte(input)))
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			require.Equal(t, rr.Code, tc.resp.Code)
			body := rr.Body.String()

			if rr.Code == http.StatusCreated {
				var resp register.SuccessResponse
				require.NoError(t, json.Unmarshal([]byte(body), &resp))
				assert.NotEqual(t, resp.User, nil)
				assert.NotEmpty(t, resp.Token)
				return
			}

			var resp register.ErrorResponse
			require.NoError(t, json.Unmarshal([]byte(body), &resp))
			require.Equal(t, tc.resp.Error, resp.Error)
		})
	}
}
