package dto

type UpdateUserDTO struct {
	Name     *string `validate:"omitempty,min=1,max=64"`
	LastName *string `validate:"omitempty,min=1,max=64"`
	Email    *string `validate:"omitempty,email"`
	Phone    *string `validate:"omitempty,max=64"`
}
