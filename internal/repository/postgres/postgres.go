package postgres

import (
	"fmt"
	"sdt-bicycle-rental/internal/config"
	"sdt-bicycle-rental/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func New(cfg config.Postgres) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		cfg.Host,
		cfg.User,
		cfg.Password,
		cfg.DBName,
		cfg.Port,
		cfg.SSLMode,
		cfg.TimeZone,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to db: %w", err)
	}

	if err := migrate(db); err != nil {
		return nil, fmt.Errorf("failed to migrate db: %w", err)
	}

	return db, nil
}

func migrate(db *gorm.DB) error {
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
