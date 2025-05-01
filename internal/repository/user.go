package repository

import "sdt-bicycle-rental/internal/models"

type UserRepository interface {
	Create(user *models.User) error
	GetByID(id uint64) (*models.User, error)
	Update(user *models.User) error
	Delete(id uint64) error
}
