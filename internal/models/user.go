package models

import "time"

type User struct {
	ID        uint64     `gorm:"primaryKey;autoIncrement;type:BIGINT"`
	Name      *string    `gorm:"type:varchar(64)"`
	LastName  *string    `gorm:"type:varchar(64)"`
	Email     *string    `gorm:"type:varchar(255);uniqueIndex"`
	Phone     *string    `gorm:"type:varchar(64);uniqueIndex"`
	Status    *string    `gorm:"type:varchar(64)"`
	Password  *string    `gorm:"type:varchar(255)"`
	CreatedAt *time.Time `gorm:"type:timestamp;default:now()"`
}
