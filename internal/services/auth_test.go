package services

import (
	"reflect"
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
			if got := NewAuthService(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewAuthService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAuth_HashPassword(t *testing.T) {
	type args struct {
		password string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "can hash new password and verify",
			args: args{
				password: "testNewPassword123",
			},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Auth{}
			hash, err := a.HashPassword(tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("HashPassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			result := a.CheckPasswordHash(tt.args.password, hash)
			if result != tt.want {
				t.Errorf("HashPassword() got = %v, want %v", result, tt.want)
			}
		})
	}
}
