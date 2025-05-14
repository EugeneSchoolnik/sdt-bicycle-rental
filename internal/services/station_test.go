package services

import (
	"log/slog"
	"reflect"
	"sdt-bicycle-rental/internal/models"
	mocks "sdt-bicycle-rental/internal/services/mocks"
	"sdt-bicycle-rental/lib/logger/handlers/slogdiscard"
	"testing"
)

func TestStationService_Create(t *testing.T) {
	type fields struct {
		repo StationRepositoty
		log  *slog.Logger
	}

	defaultFields := fields{
		repo: mocks.NewStationRepositoty(t),
		log:  slogdiscard.NewDiscardLogger(),
	}

	tests := []struct {
		name       string
		fields     fields
		argStation *models.Station
		want       *models.Station
		wantErr    bool
	}{
		{
			name:   "success",
			fields: defaultFields,
			argStation: &models.Station{
				LocationStreet: "some street 8, house 4",
			},
			want: &models.Station{
				LocationStreet: "some street 8, house 4",
			},
			wantErr: false,
		},
		{
			name:   "too small location",
			fields: defaultFields,
			argStation: &models.Station{
				LocationStreet: "some",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:   "repository error",
			fields: defaultFields,
			argStation: &models.Station{
				LocationStreet: "some street 8, house 4",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &StationService{
				repo: tt.fields.repo,
				log:  tt.fields.log,
			}

			switch tt.name {
			case "success":
				tt.fields.repo.(*mocks.StationRepositoty).On("Create", tt.argStation).Return(nil).Once()
			case "repository error":
				tt.fields.repo.(*mocks.StationRepositoty).On("Create", tt.argStation).Return(ErrInternalError).Once()
			}

			got, err := s.Create(tt.argStation)
			if (err != nil) != tt.wantErr {
				t.Errorf("StationService.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StationService.Create() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStationService_UpdateLocation(t *testing.T) {
	type fields struct {
		repo StationRepositoty
		log  *slog.Logger
	}
	defaultFields := fields{
		repo: mocks.NewStationRepositoty(t),
		log:  slogdiscard.NewDiscardLogger(),
	}
	type args struct {
		id       uint64
		location string
	}
	type mockData struct {
		notNeeded bool
		arg       *models.Station
		resp      error
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		mock    mockData
		wantErr bool
	}{
		{
			name:   "success",
			fields: defaultFields,
			args: args{
				id:       1,
				location: "some location",
			},
			mock: mockData{
				arg: &models.Station{
					ID:             1,
					LocationStreet: "some location",
				},
				resp: nil,
			},
			wantErr: false,
		},
		{
			name:   "too small location error",
			fields: defaultFields,
			args: args{
				id:       1,
				location: "some",
			},
			mock: mockData{
				notNeeded: true,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &StationService{
				repo: tt.fields.repo,
				log:  tt.fields.log,
			}

			if !tt.mock.notNeeded {
				s.repo.(*mocks.StationRepositoty).On("Update", tt.mock.arg).Return(tt.mock.resp).Once()
			}

			if err := s.UpdateLocation(tt.args.id, tt.args.location); (err != nil) != tt.wantErr {
				t.Errorf("StationService.UpdateLocation() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
