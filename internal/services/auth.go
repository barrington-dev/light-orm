package services

import (
	"errors"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"light-orm/internal/config"
)

type Auth struct{}

func NewAuthService() *Auth {
	return &Auth{}
}

func (a *Auth) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func (a *Auth) CheckPasswordHash(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func (a *Auth) NewJWTAccessToken(claims jwt.Claims) (string, error) {
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return accessToken.SignedString(config.NewJWTConfig().ConvertSecretKeyToBytes())
}

func (a *Auth) NewJWTRefreshToken(claims jwt.Claims) (string, error) {
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return refreshToken.SignedString(config.NewJWTConfig().ConvertSecretKeyToBytes())
}

func (a *Auth) ParseJWTAccessToken(tokenString string, emptyClaims jwt.Claims) (jwt.Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, emptyClaims, func(token *jwt.Token) (interface{}, error) {
		return config.NewJWTConfig().ConvertSecretKeyToBytes(), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.Claims)
	if !ok && token.Valid {
		return nil, errors.New("invalid token claims")
	}

	err = claims.Valid()
	if err != nil {
		return nil, err
	}

	return claims, nil
}
