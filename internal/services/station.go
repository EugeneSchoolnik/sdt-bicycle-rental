package services

import (
	"errors"
	"log/slog"
	"sdt-bicycle-rental/internal/models"
	"sdt-bicycle-rental/lib/logger/sl"
	"sdt-bicycle-rental/lib/validation"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

//go:generate mockery --name=StationRepositoty
type StationRepositoty interface {
	Create(station *models.Station) error
	GetByID(id uint64) (*models.Station, error)
	UpdateBikesAvailable(id uint64, delta int) error
	UpdateBikesTotal(id uint64, delta int) error
	Update(station *models.Station) error
	Delete(id uint64) error
}

type StationService struct {
	repo StationRepositoty
	log  *slog.Logger
}

func NewStationService(repo StationRepositoty, log *slog.Logger) *StationService {
	return &StationService{repo, log}
}

func (s *StationService) Create(station *models.Station) (*models.Station, error) {
	const op = "services.StationService.Create"

	// Validate station data
	err := validate.Struct(station)
	if err != nil {
		s.log.Info(op, "validation error", sl.Err(err))
		return nil, validation.PrettyError(err.(validator.ValidationErrors))
	}

	err = s.repo.Create(station)
	if err != nil {
		s.log.Error(op, "failed to create station", sl.Err(err))
		return nil, ErrInternalError
	}

	return station, nil
}

func (s *StationService) ByID(id uint64) (*models.Station, error) {
	const op = "services.StationService.ByID"

	station, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.log.Info(op, "station not found", slog.Uint64("id", id))
			return nil, ErrInvalidCredentials
		}
		s.log.Error(op, "failed to get station", sl.Err(err))
		return nil, ErrInternalError
	}

	return station, nil
}

func (s *StationService) UpdateLocation(id uint64, location string) error {
	const op = "services.StationService.UpdateLocation"

	station := &models.Station{
		ID:             id,
		LocationStreet: location,
	}

	// Validate station data
	err := validate.Struct(station)
	if err != nil {
		s.log.Info(op, "validation error", sl.Err(err))
		return validation.PrettyError(err.(validator.ValidationErrors))
	}

	err = s.repo.Update(station)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.log.Info(op, "station not found", slog.Uint64("id", id), slog.String("location", location))
			return ErrInvalidCredentials
		}
		s.log.Error(op, "failed to udpate station", sl.Err(err))
		return ErrInternalError
	}

	return nil
}

func (s *StationService) Delete(id uint64) error {
	const op = "services.StationService.Delete"

	err := s.repo.Delete(id)
	if err != nil {
		s.log.Error(op, "failed to delete station", slog.Uint64("id", id), sl.Err(err))
		return ErrInternalError
	}

	return nil
}
