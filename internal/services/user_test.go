package services

import (
	"context"
	"github.com/stretchr/testify/suite"
	"light-orm/internal/config"
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
	type args struct {
		ctx      context.Context
		user     *models.User
		password string
	}

	dbService := suite.postgresTestingService

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

	password := "password123"

	tests := []struct {
		name string
		db   database.Service
		args args
		want *models.User
	}{
		{
			name: "can create a new user",
			db:   dbService,
			args: args{
				context.Background(),
				users[0],
				password,
			},
			want: users[0],
		},
		{
			name: "can't create a new user with incomplete fields e.g first and last name",
			db:   dbService,
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
				Db: tt.db,
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
	dbService := suite.postgresTestingService

	newUserData := &models.User{
		Username:      "nelsonmandela",
		FirstName:     "Nelson",
		LastName:      "Mandela",
		ContactNumber: "0122458912",
		Email:         "nelson.mandela@gmail.com",
	}

	passwordData := "password123"

	userService := NewUserService(dbService)
	expectedUser, _ := userService.CreateUser(context.Background(), newUserData, passwordData)
	got, err := userService.GetUserById(context.Background(), expectedUser.ID)

	suite.Require().Nil(err)
	suite.Require().Equal(expectedUser, got)
}

func (suite *UserTestSuite) TestUserService_createJWTAccessTokens() {
	userData := &models.User{
		Username:      "mandylee",
		FirstName:     "Mandy",
		LastName:      "Lee",
		ContactNumber: "01234567890",
		Email:         "mandy.lee@gmail.com",
	}
	password := "password123"
	ctx := context.Background()

	userService := NewUserService(suite.postgresTestingService)
	user, err := userService.CreateUser(ctx, userData, password)
	suite.Require().Nil(err)

	suite.Run("test create jwt access tokens", func() {
		accessToken, err := userService.createJWTAccessTokens(ctx, user, NewAuthService())

		suite.Require().Nil(err)
		suite.Require().NotNil(accessToken)
	})
}

func (suite *UserTestSuite) TestUserService_Login() {
	userData := &models.User{
		Username:      "tommylee",
		FirstName:     "tommy",
		LastName:      "jones",
		ContactNumber: "01234567890",
		Email:         "tommylee.jones@gmail.com",
	}
	password := "password123"
	ctx := context.Background()

	userService := NewUserService(suite.postgresTestingService)
	user, err := userService.CreateUser(ctx, userData, password)
	suite.Require().Nil(err)

	suite.Run("test user can login", func() {
		accessToken, err := userService.Login(ctx, user, password)
		suite.Require().Nil(err)
		suite.Require().NotNil(accessToken)

		authService := NewAuthService()
		claims, err := authService.ParseJWTAccessToken(accessToken, &config.UserClaims{})
		suite.Require().Nil(err)

		err = claims.Valid()
		suite.Require().Nil(err)
	})
}
