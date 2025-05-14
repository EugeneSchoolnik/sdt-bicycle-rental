package models

import "time"

type Station struct {
	ID             uint64     `gorm:"primaryKey;autoIncrement;type:BIGINT"`
	LocationStreet string     `gorm:"type:varchar(255);not null" validate:"required,min=8,max=100"`
	BikesAvailable int        `gorm:"type:int;not null;check: bikes_available >= 0;default:0"`
	BikesTotal     int        `gorm:"type:int;not null;check: bikes_total >= 0;default:0"`
	CreatedAt      *time.Time `gorm:"type:timestamp;default:now()"`
}
