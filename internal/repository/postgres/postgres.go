package postgres

import (
	"fmt"
	"sdt-bicycle-rental/internal/config"
	"sdt-bicycle-rental/internal/repository"

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

	if err := repository.Migrate(db); err != nil {
		return nil, fmt.Errorf("failed to migrate db: %w", err)
	}

	return db, nil
}
