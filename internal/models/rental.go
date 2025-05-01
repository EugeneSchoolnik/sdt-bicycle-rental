package models

import "time"

type Rental struct {
	ID             uint64     `gorm:"primaryKey;autoIncrement;type:BIGINT"`
	UserID         uint64     `gorm:"type:BIGINT;not null"`
	BicycleID      uint64     `gorm:"type:BIGINT;not null"`
	StationStartID uint64     `gorm:"type:BIGINT;not null"`
	StationEndID   uint64     `gorm:"type:BIGINT;not null"`
	StartTime      *time.Time `gorm:"type:TIMESTAMP;not null"`
	EndTime        *time.Time `gorm:"type:TIMESTAMP;not null"`
	TotalCost      float64    `gorm:"type:DECIMAL(10,2);not null"`
	User           *User      `gorm:"foreignKey:UserID;references:ID"`
	Bicycle        *Bicycle   `gorm:"foreignKey:BicycleID;references:ID"`
	StationStart   *Station   `gorm:"foreignKey:StationStartID;references:ID"`
	StationEnd     *Station   `gorm:"foreignKey:StationEndID;references:ID"`
}
