package postgres

import (
	"errors"
	"sdt-bicycle-rental/internal/models"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *models.User) error {
	if err := r.db.Create(user).Error; err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return gorm.ErrDuplicatedKey // 23505 = unique_violation
			}
		}
		return err
	}
	return nil
}

func (r *UserRepository) GetByID(id uint64) (*models.User, error) {
	var user models.User
	if err := r.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetByIDWithRelations(id uint64) (*models.User, error) {
	var user models.User
	err := r.db.
		Preload("Bookings", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at DESC").Limit(10)
		}).
		Preload("Payments", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at DESC").Limit(10)
		}).
		Preload("Rentals", func(db *gorm.DB) *gorm.DB {
			return db.Order("start_time DESC").Limit(10)
		}).
		First(&user, id).Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) Update(user *models.User) error {
	// Updates func ignore nil fields
	tx := r.db.Updates(user)
	if err := tx.Error; err != nil {
		return err
	}
	if tx.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *UserRepository) AnonymizeAndMarkDeleted(id uint64) error {
	return r.db.Model(&models.User{}).Where("id = ?", id).
		Updates(map[string]interface{}{
			"name":       nil,
			"lastname":   nil,
			"email":      nil,
			"phone":      nil,
			"password":   nil,
			"created_at": nil,
			"status":     models.UserStatusDeleted,
		}).Error
}
