package repository_postgres_test

import (
	"sdt-bicycle-rental/internal/repository"
	"testing"

	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) (*gorm.DB, func()) {
	dsn := "host=localhost user=postgres password=postgres dbname=bicycle-rental-test port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	require.NoError(t, err)

	err = repository.Migrate(db)
	require.NoError(t, err)

	// Cleanup fucntion
	cleanup := func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}

	return db, cleanup
}

func clearTable(t *testing.T, db *gorm.DB, tableName string) {
	err := db.Exec("TRUNCATE TABLE " + tableName + " RESTART IDENTITY CASCADE").Error
	if err != nil {
		t.Fatalf("failed to truncate table %s: %v", tableName, err)
	}
}
