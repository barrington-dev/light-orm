package config

import (
	"github.com/golang-jwt/jwt"
	"os"
)

type JWTConfig struct {
	SecretKey string
}

type UserClaims struct {
	jwt.StandardClaims
	PreferredUsername string `json:"preferred_username"`
}

func NewJWTConfig() *JWTConfig {
	return &JWTConfig{
		SecretKey: os.Getenv("JWT_SECRET_KEY"),
	}
}

func (c *JWTConfig) ConvertSecretKeyToBytes() []byte {
	return []byte(c.SecretKey)
}
