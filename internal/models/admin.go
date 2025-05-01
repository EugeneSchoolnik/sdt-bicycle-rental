package models

type Admin struct {
	UserID uint64 `gorm:"primaryKey"`
	User   *User  `gorm:"foreignKey:UserID;references:ID"`
}
