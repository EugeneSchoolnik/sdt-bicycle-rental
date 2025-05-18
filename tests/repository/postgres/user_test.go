package repository_postgres_test

import (
	"sdt-bicycle-rental/internal/models"
	"sdt-bicycle-rental/internal/repository/postgres"
	. "sdt-bicycle-rental/lib/util"
	test_postgres "sdt-bicycle-rental/tests/util/db/postgres"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestUserRepository(t *testing.T) {
	db, cleanup := test_postgres.SetupTestDB(t)
	defer cleanup()

	test_postgres.ClearTable(t, db, "users")

	repo := postgres.NewUserRepository(db)

	user := &models.User{
		Name:     Ptr("Test"),
		Lastname: Ptr("User"),
		Email:    Ptr("test@example.com"),
		Phone:    Ptr("123456"),
		Status:   Ptr("active"),
		Password: Ptr("password123"),
	}

	t.Run("create", func(t *testing.T) {
		err := repo.Create(user)
		require.NoError(t, err)
		require.NotZero(t, user.ID)

		// Check in database
		var saved models.User
		err = db.First(&saved, user.ID).Error
		assert.NoError(t, err)
		assert.Equal(t, *user.Email, *saved.Email)
	})

	t.Run("get by id", func(t *testing.T) {
		saved, err := repo.GetByID(user.ID)
		require.NoError(t, err)
		assert.Equal(t, *user.Email, *saved.Email)

		// user not found
		_, err = repo.GetByID(404)
		assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
	})

	t.Run("create duplicate", func(t *testing.T) {
		err := repo.Create(user)
		require.ErrorIs(t, err, gorm.ErrDuplicatedKey)
	})

	t.Run("update", func(t *testing.T) {
		user.Name = Ptr("Updated")
		err := repo.Update(user)
		require.NoError(t, err)

		updated, err := repo.GetByID(user.ID)
		require.NoError(t, err)
		require.Equal(t, *user.Name, *updated.Name)
	})

	t.Run("get by email", func(t *testing.T) {
		saved, err := repo.GetByEmail(*user.Email)
		require.NoError(t, err)
		assert.Equal(t, *user.Email, *saved.Email)
	})

	t.Run("anonymize and mark deleted", func(t *testing.T) {
		err := repo.AnonymizeAndMarkDeleted(user.ID)
		require.NoError(t, err)
		deletedUser, err := repo.GetByID(user.ID)
		require.NoError(t, err)
		assert.NotZero(t, deletedUser.ID)
		assert.Empty(t, deletedUser.Name)
		assert.Empty(t, deletedUser.Lastname)
		assert.Empty(t, deletedUser.Email)
		assert.Empty(t, deletedUser.Phone)
		assert.Empty(t, deletedUser.Password)
		assert.Empty(t, deletedUser.CreatedAt)

		require.NotEmpty(t, deletedUser.Status)
		assert.Equal(t, *deletedUser.Status, models.UserStatusDeleted)
	})
}
