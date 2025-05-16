package models

import (
	"time"
)

const (
	UserStatusActive  = "active"
	UserStatusDeleted = "deleted"
	UserStatusBanned  = "banned"
)

type User struct {
	ID        uint64     `gorm:"primaryKey;autoIncrement;type:BIGINT"`
	Name      *string    `gorm:"type:varchar(64)" validate:"required,min=1,max=64"`
	Lastname  *string    `gorm:"type:varchar(64)" validate:"required,min=1,max=64"`
	Email     *string    `gorm:"type:varchar(255);uniqueIndex" validate:"required,email"`
	Phone     *string    `gorm:"type:varchar(64);uniqueIndex" validate:"required,max=64"`
	Status    *string    `gorm:"type:varchar(64)"`
	Password  *string    `gorm:"type:varchar(255)" validate:"required,min=8,max=255"`
	CreatedAt *time.Time `gorm:"type:timestamp;default:now()"`

	Bookings []Booking `gorm:"foreignKey:UserID;references:ID"`
	Payments []Payment `gorm:"foreignKey:UserID;references:ID"`
	Rentals  []Rental  `gorm:"foreignKey:UserID;references:ID"`
}
