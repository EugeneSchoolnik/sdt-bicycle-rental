package models

import "time"

const (
	BicycleStatusAvailable = "available"
	BicycleStatusRented    = "rented"
	BicycleStatusInService = "in_service"
)

type Bicycle struct {
	ID          uint64     `gorm:"primaryKey;autoIncrement;type:BIGINT"`
	StationID   uint64     `gorm:"type:BIGINT;not null"`
	Status      string     `gorm:"type:varchar(64);not null;"`
	LastService *time.Time `gorm:"type:timestamp"`

	Station *Station `gorm:"foreignKey:StationID;references:ID"`
}
