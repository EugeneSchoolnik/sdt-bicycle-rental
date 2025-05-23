package models

import "time"

type Payment struct {
	ID            uint64     `gorm:"primaryKey;autoIncrement;type:BIGINT"`
	UserID        uint64     `gorm:"type:BIGINT;not null"`
	Method        string     `gorm:"type:varchar(64);not null"`
	Amount        float64    `gorm:"type:decimal(10,2);not null"`
	TransactionID string     `gorm:"type:varchar(255)"`
	Status        string     `gorm:"type:varchar(64);not null"`
	CreatedAt     *time.Time `gorm:"type:timestamp;default:now()"`
	User          *User      `gorm:"foreignKey:UserID;references:ID"`
}
