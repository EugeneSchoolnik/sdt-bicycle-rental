package service

import "errors"

var (
	// Common
	ErrInternalError = errors.New("internal server error")

	// Auth
	ErrExpiredToken       = errors.New("token expired")
	ErrInvalidToken       = errors.New("invalid token")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
)
