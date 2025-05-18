package dto

import "sdt-bicycle-rental/internal/models"

type CreateUser struct {
	Name     string `validate:"required,min=1,max=64"`
	Lastname string `validate:"required,min=1,max=64"`
	Email    string `validate:"required,email"`
	Phone    string `validate:"required,max=64"`
	Password string `validate:"required,min=8,max=255"`
}

func (dto *CreateUser) Model() *models.User {
	return &models.User{
		Name:     &dto.Name,
		Lastname: &dto.Lastname,
		Email:    &dto.Email,
		Phone:    &dto.Phone,
		Password: &dto.Password,
	}
}

type UpdateUser struct {
	Name     *string `validate:"omitempty,min=1,max=64"`
	Lastname *string `validate:"omitempty,min=1,max=64"`
	Email    *string `validate:"omitempty,email"`
	Phone    *string `validate:"omitempty,max=64"`
}
