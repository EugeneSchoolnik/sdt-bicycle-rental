package repository

import (
	"fmt"
	"sdt-bicycle-rental/internal/models"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	var modelsToMigrate = []any{
		&models.User{},
		&models.Admin{},
		&models.Station{},
		&models.Bicycle{},
		&models.Payment{},
		&models.Booking{},
		&models.Rental{},
	}

	for _, model := range modelsToMigrate {
		if err := db.AutoMigrate(model); err != nil {
			return fmt.Errorf("failed to migrate model %T: %w", model, err)
		}
	}

	return nil
}
