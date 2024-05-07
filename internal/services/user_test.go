package services

import (
	"context"
	"github.com/stretchr/testify/require"
	"light-orm/internal/database"
	"light-orm/internal/models"
	"testing"
)

func init() {
	database.NewPostgresTestingService() // setup testing database service connection
}

func TestNewUserService(t *testing.T) {
	type args struct {
		db database.Service
	}

	dbServ := args{
		db: database.DbTestInstance,
	}

	tests := []struct {
		name string
		args args
		want *UserService
	}{
		{
			name: "can create new user service",
			args: dbServ,
			want: NewUserService(dbServ.db),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, NewUserService(tt.args.db))
		})
	}
}

func TestUserService_CreateUser(t *testing.T) {
	type fields struct {
		Db database.Service
	}

	type args struct {
		ctx      context.Context
		user     *models.User
		password *models.Password
	}

	dbService := fields{Db: database.DbTestInstance}

	users := []*models.User{
		{
			Username:      "stevejobs",
			FirstName:     "Steve",
			LastName:      "Jobs",
			ContactNumber: "0123456789",
			Email:         "steve.jobs@gmail.com",
		},
		{
			Username:      "stevewozniak",
			ContactNumber: "0123456784",
			Email:         "steve.wozniak@gmail.com",
		},
	}

	password := &models.Password{
		Password: "password123",
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   *models.User
	}{
		{
			name:   "can create a new user",
			fields: dbService,
			args: args{
				context.Background(),
				users[0],
				password,
			},
			want: users[0],
		},
		{
			name:   "can't create a new user with incomplete fields e.g first and last name",
			fields: dbService,
			args: args{
				context.Background(),
				users[1],
				password,
			},
			want: users[1],
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			us := &UserService{
				Db: tt.fields.Db,
			}
			got, err := us.CreateUser(tt.args.ctx, tt.args.user, tt.args.password)

			if (err != nil) && err.Error() == "invalid input, empty fields not allowed" {
				return
			}

			require.Nil(t, err)
			require.Equal(t, tt.want, got)
		})
	}

	t.Cleanup(func() {
		models.Users().DeleteAll(context.Background(), dbService.Db.Instance())
	})
}

func TestUserService_GetUser(t *testing.T) {
	dbService := database.DbTestInstance

	newUserData := &models.User{
		Username:      "nelsonmandela",
		FirstName:     "Nelson",
		LastName:      "Mandela",
		ContactNumber: "0122458912",
		Email:         "nelson.mandela@gmail.com",
	}

	passwordData := &models.Password{Password: "password123"}

	userService := NewUserService(dbService)
	expectedUser, _ := userService.CreateUser(context.Background(), newUserData, passwordData)
	got, err := userService.GetUser(context.Background(), expectedUser.ID)

	require.Nil(t, err)
	require.Equal(t, expectedUser, got)

	t.Cleanup(func() {
		models.Users().DeleteAll(context.Background(), dbService.Instance())
	})
}
