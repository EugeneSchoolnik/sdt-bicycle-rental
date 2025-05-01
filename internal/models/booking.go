package models

import "time"

type Booking struct {
	ID        uint64     `gorm:"primaryKey;autoIncrement;type:BIGINT"`
	UserID    uint64     `gorm:"type:BIGINT;not null"`
	BicycleID uint64     `gorm:"type:BIGINT;not null"`
	StationID uint64     `gorm:"type:BIGINT;not null"`
	PaymentID uint64     `gorm:"type:BIGINT;not null"`
	CreatedAt *time.Time `gorm:"type:timestamp;default:now()"`
	ExpiresAt *time.Time `gorm:"type:timestamp"`
	User      *User      `gorm:"foreignKey:UserID;references:ID"`
	Bicycle   *Bicycle   `gorm:"foreignKey:BicycleID;references:ID"`
	Station   *Station   `gorm:"foreignKey:StationID;references:ID"`
	Payment   *Payment   `gorm:"foreignKey:PaymentID;references:ID"`
}
