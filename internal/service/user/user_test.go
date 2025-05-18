package user_service_test

import (
	"log/slog"
	"reflect"
	"sdt-bicycle-rental/internal/models"
	"sdt-bicycle-rental/internal/repository/dto"
	user_service "sdt-bicycle-rental/internal/service/user"
	mocks "sdt-bicycle-rental/internal/service/user/mocks"
	"sdt-bicycle-rental/lib/logger/handlers/slogdiscard"
	"sdt-bicycle-rental/lib/util"
	"testing"

	"gorm.io/gorm"
)

func TestUserService_ProfileByID(t *testing.T) {
	type fields struct {
		repo user_service.UserRepository
		log  *slog.Logger
	}

	defaultFields := fields{
		repo: mocks.NewUserRepository(t),
		log:  slogdiscard.NewDiscardLogger(),
	}

	tests := []struct {
		name    string
		fields  fields
		argID   uint64
		want    *models.User
		wantErr bool
	}{
		{
			name:   "success",
			fields: defaultFields,
			argID:  43,
			want: &models.User{
				ID:       43,
				Name:     util.Ptr("John"),
				Lastname: util.Ptr("Doe"),
				Email:    util.Ptr("valid@email.com"),
			},
			wantErr: false,
		},
		{
			name:    "not found",
			fields:  defaultFields,
			argID:   404,
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := user_service.New(tt.fields.repo, tt.fields.log)

			switch tt.name {
			case "success":
				tt.fields.repo.(*mocks.UserRepository).On("GetByIDWithRelations", tt.argID).Return(tt.want, nil).Once()
			case "not found":
				tt.fields.repo.(*mocks.UserRepository).On("GetByIDWithRelations", tt.argID).Return(tt.want, gorm.ErrRecordNotFound).Once()
			}

			got, err := s.ProfileByID(tt.argID)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserService.ProfileByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserService.ProfileByID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserService_Update(t *testing.T) {
	type fields struct {
		repo user_service.UserRepository
		log  *slog.Logger
	}

	defaultFields := fields{
		repo: mocks.NewUserRepository(t),
		log:  slogdiscard.NewDiscardLogger(),
	}

	type args struct {
		id   uint64
		user *dto.UpdateUser
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "successfully",
			fields: defaultFields,
			args: args{
				id: 1,
				user: &dto.UpdateUser{
					Name:     util.Ptr("John"),
					Lastname: util.Ptr("Doe"),
					Email:    util.Ptr("valid@email.com"),
				},
			},
			wantErr: false,
		},
		{
			name:   "empty name",
			fields: defaultFields,
			args: args{
				id: 1,
				user: &dto.UpdateUser{
					Name:  util.Ptr(""),
					Email: util.Ptr("valid@email.com"),
				},
			},
			wantErr: true,
		},
		{
			name:   "invalid email",
			fields: defaultFields,
			args: args{
				id: 1,
				user: &dto.UpdateUser{
					Name:  util.Ptr("Tomas"),
					Email: util.Ptr("invalid@email"),
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := user_service.New(tt.fields.repo, tt.fields.log)

			updateModel := models.User{
				ID:       tt.args.id,
				Name:     tt.args.user.Name,
				Lastname: tt.args.user.Lastname,
				Email:    tt.args.user.Email,
				Phone:    tt.args.user.Phone,
			}

			switch tt.name {
			case "successfully":
				tt.fields.repo.(*mocks.UserRepository).On("Update", &updateModel).Return(nil).Once()
			}

			if err := s.Update(tt.args.id, tt.args.user); (err != nil) != tt.wantErr {
				t.Errorf("UserService.Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUserService_Delete(t *testing.T) {
	type fields struct {
		repo user_service.UserRepository
		log  *slog.Logger
	}

	defaultFields := fields{
		repo: mocks.NewUserRepository(t),
		log:  slogdiscard.NewDiscardLogger(),
	}

	tests := []struct {
		name    string
		fields  fields
		argID   uint64
		wantErr bool
	}{
		{
			name:    "success",
			fields:  defaultFields,
			argID:   1,
			wantErr: false,
		},
		{
			name:    "not found",
			fields:  defaultFields,
			argID:   1,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := user_service.New(tt.fields.repo, tt.fields.log)

			switch tt.name {
			case "success":
				tt.fields.repo.(*mocks.UserRepository).On("AnonymizeAndMarkDeleted", tt.argID).Return(nil).Once()
			case "not found":
				tt.fields.repo.(*mocks.UserRepository).On("AnonymizeAndMarkDeleted", tt.argID).Return(gorm.ErrRecordNotFound).Once()
			}

			if err := s.Delete(tt.argID); (err != nil) != tt.wantErr {
				t.Errorf("UserService.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
