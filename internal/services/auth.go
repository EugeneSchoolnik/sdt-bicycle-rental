package services

import (
	"errors"
	"fmt"
	"log/slog"
	"sdt-bicycle-rental/internal/models"
	"sdt-bicycle-rental/lib/logger/sl"
	. "sdt-bicycle-rental/lib/ptr"
	"sdt-bicycle-rental/lib/validation"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

//go:generate mockery --name=UserRepository
type UserRepository interface {
	Create(user *models.User) error
	GetByID(id uint64) (*models.User, error)
	GetByIDWithRelations(id uint64) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	Update(user *models.User) error
	AnonymizeAndMarkDeleted(id uint64) error
}

type AuthService struct {
	repo      UserRepository
	log       *slog.Logger
	jwtSecret string
}

func NewAuthService(repo UserRepository, log *slog.Logger, jwtSecret string) *AuthService {
	return &AuthService{repo: repo, log: log, jwtSecret: jwtSecret}
}

func (s *AuthService) Register(user *models.User) (*models.User, string, error) {
	const op = "services.AuthService.Register"

	// Validate user data
	err := validate.Struct(user)
	if err != nil {
		s.log.Info(op, "validation error", sl.Err(err))
		return nil, "", validation.PrettyError(err.(validator.ValidationErrors))
	}

	// Hash password
	hashedPassword, err := s.hashPassword(*user.Password)
	if err != nil {
		s.log.Error(op, "failed to hash password", sl.Err(err))
		return nil, "", ErrInternalError
	}
	// Set hashed password
	user.Password = &hashedPassword

	// Set user status
	user.Status = Ptr(models.UserStatusActive)

	// Create new user
	err = s.repo.Create(user)
	if err != nil {
		// Сheck if user already exists
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			s.log.Info(op, "user already exists", slog.String("email", *user.Email))
			return nil, "", ErrUserAlreadyExists
		}
		s.log.Error(op, "failed to create user", sl.Err(err))
		return nil, "", ErrInternalError
	}

	// Generate JWT token
	token, err := s.GenerateToken(user)
	if err != nil {
		s.log.Error(op, "failed to generate token", sl.Err(err))
		return nil, "", ErrInternalError
	}

	return user, token, nil
}

func (s *AuthService) Login(email, password string) (*models.User, string, error) {
	const op = "services.AuthService.Login"

	// Validate email and password
	if validate.Var(email, "required,email") != nil || validate.Var(password, "required,min=8,max=255") != nil {
		s.log.Info(op, "validation error", slog.String("error", "invalid email or password"))
		return nil, "", ErrInvalidCredentials
	}

	// Get user by email
	user, err := s.repo.GetByEmail(email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.log.Info(op, "user not found", slog.String("email", email))
			fmt.Println("User not found:", err)
			return nil, "", ErrInvalidCredentials
		}
		// Handle other errors
		s.log.Error(op, "failed to get user", sl.Err(err))
		return nil, "", ErrInternalError
	}

	// Check password
	if !s.checkPassword(*user.Password, password) {
		return nil, "", ErrInvalidCredentials
	}

	// Generate JWT token
	token, err := s.GenerateToken(user)
	if err != nil {
		s.log.Error(op, "failed to generate token", sl.Err(err))
		return nil, "", ErrInternalError
	}

	return user, token, nil
}

func (s *AuthService) GenerateToken(user *models.User) (string, error) {
	// Define expiration time for the token
	expirationTime := time.Now().Add(24 * time.Hour)

	// Create claims (payload) for the token
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"exp":     expirationTime.Unix(),
	}

	// Create a new token with the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key
	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		s.log.Error("failed to sign token", sl.Err(err))
		return "", err
	}

	return tokenString, nil
}

func (s *AuthService) ValidateToken(tokenString string) (*jwt.Token, error) {
	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Check if the signing method is valid
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			s.log.Error("unexpected signing method", slog.String("method", token.Header["alg"].(string)))
			return nil, ErrInvalidToken
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		s.log.Error("failed to parse token", sl.Err(err))
		return nil, ErrInvalidToken
	}

	if !token.Valid {
		s.log.Info("invalid token", slog.String("token", tokenString))
		return nil, ErrInvalidToken
	}

	return token, nil
}

func (s *AuthService) hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func (s *AuthService) checkPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
