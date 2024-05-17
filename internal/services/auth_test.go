package services

import (
	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/suite"
	"light-orm/internal/config"
	"strconv"
	"testing"
	"time"
)

type AuthTestSuite struct {
	suite.Suite
}

func TestAuthTestSuite(t *testing.T) {
	suite.Run(t, new(AuthTestSuite))
}

func (suite *AuthTestSuite) TestNewAuthService() {
	tests := []struct {
		name string
		want *Auth
	}{
		{
			name: "can create new auth service",
			want: &Auth{},
		},
	}
	for _, tt := range tests {
		suite.Run(tt.name, func() {
			suite.Require().Equal(tt.want, NewAuthService())
		})
	}
}

func (suite *AuthTestSuite) TestAuth_HashPassword() {
	type args struct {
		password string
	}
	tests := []struct {
		name string
		args args
		want error
	}{
		{
			name: "can hash new password and verify",
			args: args{
				password: "testNewPassword123",
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		suite.Run(tt.name, func() {
			a := &Auth{}
			hash, err := a.HashPassword(tt.args.password)
			suite.Require().Nil(err)

			result := a.CheckPasswordHash(tt.args.password, hash)
			suite.Require().Equal(tt.want, result)
		})
	}
}

func (suite *AuthTestSuite) TestAuth_NewJWTAccessToken() {
	userId := int64(1234567890)
	authService := NewAuthService()
	issuedAtTime := time.Now()

	suite.Run("can create a new jwt access token", func() {
		accessToken, err := authService.NewJWTAccessToken(config.UserClaims{
			StandardClaims: jwt.StandardClaims{
				Audience:  "",
				ExpiresAt: issuedAtTime.Add(time.Minute * 10).Unix(),
				Id:        "",
				IssuedAt:  issuedAtTime.Unix(),
				Issuer:    "",
				NotBefore: 0,
				Subject:   strconv.FormatInt(userId, 10),
			},
			UserName: "banana",
		})

		suite.Require().Nil(err)
		suite.Require().NotEmpty(accessToken)
	})
}

func (suite *AuthTestSuite) TestAuth_NewJWTRefreshToken() {
	userId := int64(1234567890)
	authService := NewAuthService()
	issuedAtTime := time.Now()

	suite.Run("can create refresh access token", func() {
		refreshToken, err := authService.NewJWTRefreshToken(jwt.StandardClaims{
			Audience:  "",
			ExpiresAt: issuedAtTime.Add(time.Hour * 24 * 7).Unix(),
			Id:        "",
			IssuedAt:  issuedAtTime.Unix(),
			Issuer:    "",
			NotBefore: 0,
			Subject:   strconv.FormatInt(userId, 10),
		})

		suite.Require().Nil(err)
		suite.Require().NotEmpty(refreshToken)
	})
}

func (suite *AuthTestSuite) TestAuth_ParseJWTAccessToken() {
	userId := int64(1234567890)
	authService := NewAuthService()
	issuedAtTime := time.Now()

	expectedClaims := &config.UserClaims{
		StandardClaims: jwt.StandardClaims{
			Audience:  "",
			ExpiresAt: issuedAtTime.Add(time.Minute * 10).Unix(),
			Id:        "",
			IssuedAt:  issuedAtTime.Unix(),
			Issuer:    "",
			NotBefore: 0,
			Subject:   strconv.FormatInt(userId, 10),
		},
		UserName: "honey",
	}

	accessToken, err := authService.NewJWTAccessToken(expectedClaims)
	suite.Require().Nil(err)

	suite.Run("can parse access token", func() {
		claims, err := authService.ParseJWTAccessToken(accessToken, &config.UserClaims{})

		suite.Require().Nil(err)
		suite.Require().Equal(expectedClaims, claims)
	})
}
