package services

import (
	"errors"

	"github.com/go-playground/validator/v10"
)

var (
	// Common
	ErrInternalError = errors.New("internal server error")

	// Auth
	ErrExpiredToken       = errors.New("token expired")
	ErrInvalidToken       = errors.New("invalid token")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

var validate = validator.New(validator.WithRequiredStructEnabled())
