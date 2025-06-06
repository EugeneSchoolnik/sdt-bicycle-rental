package auth_service_test

import (
	"log/slog"
	"reflect"
	"sdt-bicycle-rental/internal/models"
	"sdt-bicycle-rental/internal/repository/dto"
	"sdt-bicycle-rental/internal/service"
	auth_service "sdt-bicycle-rental/internal/service/auth"

	mocks "sdt-bicycle-rental/internal/service/auth/mocks"
	"sdt-bicycle-rental/lib/logger/handlers/slogdiscard"
	"sdt-bicycle-rental/lib/util"
	"testing"

	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

const (
	validEmail   = "valid@email.com"
	invalidEmail = "invalid-email"
)

func TestAuthService_Register(t *testing.T) {
	type fields struct {
		repo      auth_service.UserRepository
		log       *slog.Logger
		jwtSecret string
	}

	defaultFields := fields{
		repo:      mocks.NewUserRepository(t),
		log:       slogdiscard.NewDiscardLogger(),
		jwtSecret: "secret",
	}

	tests := []struct {
		name    string
		fields  fields
		argUser *dto.CreateUser
		want    *models.User
		wantErr bool
	}{
		{
			name:   "success",
			fields: defaultFields,
			argUser: &dto.CreateUser{
				Name:     "John",
				Lastname: "Doe",
				Email:    validEmail,
				Phone:    "1234567890",
				Password: "password",
			},
			want: &models.User{
				Name:     util.Ptr("John"),
				Lastname: util.Ptr("Doe"),
				Email:    util.Ptr(validEmail),
				Phone:    util.Ptr("1234567890"),
				Status:   util.Ptr(models.UserStatusActive),
				Password: nil,
			},
			wantErr: false,
		},
		{
			name:   "validation error: name",
			fields: defaultFields,
			argUser: &dto.CreateUser{
				Name:     "",
				Lastname: "Doe",
				Email:    validEmail,
				Phone:    "1234567890",
				Password: "password",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:   "validation error: email",
			fields: defaultFields,
			argUser: &dto.CreateUser{
				Name:     "John",
				Lastname: "Doe",
				Email:    "invalid-emal",
				Phone:    "1234567890",
				Password: "password",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:   "validation error: password",
			fields: defaultFields,
			argUser: &dto.CreateUser{
				Name:     "John",
				Lastname: "Doe",
				Email:    validEmail,
				Phone:    "1234567890",
				Password: "",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:   "create error",
			fields: defaultFields,
			argUser: &dto.CreateUser{
				Name:     "John",
				Lastname: "Doe",
				Email:    validEmail,
				Phone:    "1234567890",
				Password: "password",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := auth_service.New(tt.fields.repo, tt.fields.log, tt.fields.jwtSecret)

			switch tt.name {
			case "success":
				tt.fields.repo.(*mocks.UserRepository).
					On("Create", mock.MatchedBy(func(u *models.User) bool { return true })).
					Return(nil).Once()
			case "create error":
				tt.fields.repo.(*mocks.UserRepository).
					On("Create", mock.MatchedBy(func(u *models.User) bool { return true })).
					Return(service.ErrInternalError).Once()
			}

			got, got1, err := s.Register(tt.argUser)
			isErr := err != nil

			if isErr != tt.wantErr {
				t.Errorf("UserService.Register() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !isErr {
				if bcrypt.CompareHashAndPassword([]byte(*got.Password), []byte("password")) != nil {
					t.Errorf("UserService.Register() password hash mismatch")
				}

				got.Password = nil // Remove password because we already checked it
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("UserService.Register() got = %v, want %v", got, tt.want)
				}

				token, err := s.ValidateToken(got1)
				if err != nil || token == nil {
					t.Errorf("UserService.Register() token validation error = %v", err)
				}
			}
		})
	}
}

func TestAuthService_Login(t *testing.T) {
	type fields struct {
		repo      auth_service.UserRepository
		log       *slog.Logger
		jwtSecret string
	}

	defaultFields := fields{
		repo:      mocks.NewUserRepository(t),
		log:       slogdiscard.NewDiscardLogger(),
		jwtSecret: "secret",
	}
	type args struct {
		email    string
		password string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *models.User
		wantErr bool
	}{
		{
			name:   "success",
			fields: defaultFields,
			args: args{
				email:    validEmail,
				password: "password",
			},
			want: &models.User{
				Name:     util.Ptr("John"),
				Lastname: util.Ptr("Doe"),
				Email:    util.Ptr(validEmail),
				Phone:    util.Ptr("1234567890"),
				Status:   util.Ptr(models.UserStatusActive),
				Password: util.Ptr("$2a$10$Qz4ERCPWmdyNe7DR5H19RubOlA7drtlD9VCVYl8N9QjcqhueonsM6"),
			},
			wantErr: false,
		},
		{
			name:   "invalid email",
			fields: defaultFields,
			args: args{
				email:    invalidEmail,
				password: "password",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:   "too small password",
			fields: defaultFields,
			args: args{
				email:    validEmail,
				password: "pass",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := auth_service.New(tt.fields.repo, tt.fields.log, tt.fields.jwtSecret)

			switch tt.name {
			case "success":
				tt.fields.repo.(*mocks.UserRepository).On("GetByEmail", tt.args.email).Return(tt.want, nil).Once()
			}

			got, got1, err := s.Login(tt.args.email, tt.args.password)
			t.Logf("Error Message: %v", err)
			isErr := err != nil
			if isErr != tt.wantErr {
				t.Errorf("AuthService.Login() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !isErr {
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("AuthService.Login() got = %v, want %v", got, tt.want)
				}

				token, err := s.ValidateToken(got1)
				if err != nil || token == nil {
					t.Errorf("AuthService.Login() token validation error = %v", err)
				}
			}
		})
	}
}
