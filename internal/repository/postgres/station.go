package postgres

import (
	"sdt-bicycle-rental/internal/models"

	"gorm.io/gorm"
)

type StationRepository struct {
	db *gorm.DB
}

func NewStationRepository(db *gorm.DB) *StationRepository {
	return &StationRepository{db: db}
}

func (r *StationRepository) Create(station *models.Station) error {
	return r.db.Create(station).Error
}

func (r *StationRepository) GetByID(id uint64) (*models.Station, error) {
	var station models.Station
	if err := r.db.First(&station, id).Error; err != nil {
		return nil, err
	}
	return &station, nil
}

func (r *StationRepository) Update(station *models.Station) error {
	tx := r.db.Updates(station)
	if err := tx.Error; err != nil {
		return err
	}
	if tx.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *StationRepository) Delete(id uint64) error {
	return r.db.Delete(&models.Station{}, id).Error
}

func (r *StationRepository) UpdateBikesAvailable(id uint64, delta int) error {
	tx := r.db.Model(&models.Station{}).Where("id = ?", id).
		Update("bikes_available", gorm.Expr("bikes_available + ?", delta))

	if err := tx.Error; err != nil {
		return err
	}
	if tx.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *StationRepository) UpdateBikesTotal(id uint64, delta int) error {
	tx := r.db.Model(&models.Station{}).Where("id = ?", id).
		Update("bikes_total", gorm.Expr("bikes_total + ?", delta))

	if err := tx.Error; err != nil {
		return err
	}
	if tx.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
