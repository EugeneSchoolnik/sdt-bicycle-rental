package services

import (
	"errors"
	"log/slog"
	"sdt-bicycle-rental/internal/models"
	"sdt-bicycle-rental/internal/repository/dto"
	"sdt-bicycle-rental/lib/logger/sl"
	"sdt-bicycle-rental/lib/validation"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type UserService struct {
	repo UserRepository
	log  *slog.Logger
}

func NewUserService(repo UserRepository, log *slog.Logger) *UserService {
	return &UserService{repo: repo, log: log}
}

func (s *UserService) ProfileByID(id uint64) (*models.User, error) {
	const op = "services.UserService.ProfileByID"

	// Get user by ID
	user, err := s.repo.GetByIDWithRelations(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.log.Info(op, "user not found", slog.Uint64("id", id))
			return nil, ErrInvalidCredentials
		}
		// Handle other errors
		s.log.Error(op, "failed to get user", sl.Err(err))
		return nil, ErrInternalError
	}

	return user, nil
}

// TODO: return updated user
func (s *UserService) Update(id uint64, user *dto.UpdateUserDTO) error {
	const op = "services.UserService.Update"

	// Validate user
	err := validate.Struct(user)
	if err != nil {
		s.log.Error(op, "validation failed", slog.String("error", err.Error()))
		return validation.PrettyError(err.(validator.ValidationErrors))
	}

	updateUser := models.User{
		ID:       id,
		Name:     user.Name,
		LastName: user.LastName,
		Email:    user.Email,
		Phone:    user.Phone,
	}

	// Update user
	err = s.repo.Update(&updateUser)
	if err != nil {
		s.log.Error(op, "failed to update user", slog.String("error", err.Error()))
		return ErrInternalError
	}

	return nil
}

func (s *UserService) Delete(id uint64) error {
	const op = "services.UserService.Delete"

	// Delete user
	err := s.repo.AnonymizeAndMarkDeleted(id)
	if err != nil {
		s.log.Error(op, "failed to delete user", slog.String("error", err.Error()))
		return ErrInternalError
	}

	return nil
}
