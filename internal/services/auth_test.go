package services

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewAuthService(t *testing.T) {
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
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, NewAuthService())
		})
	}
}

func TestAuth_HashPassword(t *testing.T) {
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
		t.Run(tt.name, func(t *testing.T) {
			a := &Auth{}
			hash, err := a.HashPassword(tt.args.password)
			require.Nil(t, err)

			result := a.CheckPasswordHash(tt.args.password, hash)
			require.Equal(t, tt.want, result)
		})
	}
}
