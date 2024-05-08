package services

import (
	"github.com/stretchr/testify/suite"
	"testing"
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
		want bool
	}{
		{
			name: "can hash new password and verify",
			args: args{
				password: "testNewPassword123",
			},
			want: true,
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
