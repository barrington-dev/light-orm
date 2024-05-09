package services

import (
	"context"
	"github.com/stretchr/testify/suite"
	"light-orm/internal/database"
	"light-orm/internal/models"
	"testing"
)

type UserTestSuite struct {
	suite.Suite
	postgresTestingService *database.PostgresTestingService
}

func TestUserTestSuite(t *testing.T) {
	suite.Run(t, new(UserTestSuite))
}

func (suite *UserTestSuite) SetupSuite() {
	suite.T().Log("setting up postgres testing service")
	suite.postgresTestingService = database.NewPostgresTestingService() // setup testing database service connection
}

func (suite *UserTestSuite) TearDownSuite() {
	suite.T().Log("tearing down postgres testing service")
	if err := suite.postgresTestingService.Close(); err != nil {
		suite.FailNow("failed to close postgres testing service", "error", err.Error())
	}
}

func (suite *UserTestSuite) AfterTest(suiteName, testName string) {
	if testName == "TestNewUserService" {
		return
	}

	suite.T().Logf("truncating users and related tables for %s...", testName)
	_, err := suite.postgresTestingService.Instance().Exec("TRUNCATE TABLE users RESTART IDENTITY CASCADE")
	suite.Require().Nil(err)
	suite.T().Log("truncation for users and related tables complete")
}

func (suite *UserTestSuite) TestNewUserService() {
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
		suite.Run(tt.name, func() {
			suite.Require().Equal(tt.want, NewUserService(tt.args.db))
		})
	}
}

func (suite *UserTestSuite) TestUserService_CreateUser() {
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
		suite.Run(tt.name, func() {
			us := &UserService{
				Db: tt.fields.Db,
			}
			got, err := us.CreateUser(tt.args.ctx, tt.args.user, tt.args.password)

			if (err != nil) && err.Error() == "invalid input, empty fields not allowed" {
				return
			}

			suite.Require().Nil(err)
			suite.Require().Equal(tt.want, got)
		})
	}
}

func (suite *UserTestSuite) TestUserService_GetUser() {
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
	got, err := userService.GetUserById(context.Background(), expectedUser.ID)

	suite.Require().Nil(err)
	suite.Require().Equal(expectedUser, got)
}
