package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"light-orm/internal/config"
	"light-orm/internal/database"
	"light-orm/internal/models"
	"strconv"
	"time"
)

type UserService struct {
	Db database.Service
}

func NewUserService(db database.Service) *UserService {
	models.AddUserHook(boil.BeforeInsertHook, validateEmptyUserFieldsHook)
	return &UserService{
		Db: db,
	}
}

func (us *UserService) GetUserById(ctx context.Context, id int64) (*models.User, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	user, err := models.Users(qm.Where("id=?", id)).One(ctx, us.Db.Instance())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New(fmt.Sprintf("user %d not found", id))
		}
		return nil, err
	}
	return user, nil
}

func (us *UserService) CreateUser(ctx context.Context, user *models.User, clearPassword string) (*models.User, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	err := user.Insert(ctx, us.Db.Instance(), boil.Infer())
	if err != nil {
		return nil, err
	}

	// Hash password
	hash, err := NewAuthService().HashPassword(clearPassword)
	if err != nil {
		return nil, err
	}

	// Update password with hash
	password := &models.Password{
		Hash: hash,
	}

	// Store password in db
	err = user.AddPasswords(ctx, us.Db.Instance(), true, password)
	if err != nil {
		return nil, err
	}

	if err = user.Reload(ctx, us.Db.Instance()); err != nil {
		return nil, err
	}

	return user, nil
}

func (us *UserService) Login(ctx context.Context, user *models.User, clearPassword string) (string, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	password, err := user.Passwords(qm.OrderBy("created_at DESC")).One(ctx, us.Db.Instance())
	if err != nil {
		return "", err
	}

	authService := NewAuthService()
	err = authService.CheckPasswordHash(clearPassword, password.Hash)
	if err != nil {
		return "", err
	}

	issuedAtTime := time.Now()
	accessToken, err := authService.NewJWTAccessToken(config.UserClaims{
		StandardClaims: jwt.StandardClaims{
			Audience:  "",
			ExpiresAt: issuedAtTime.Add(time.Minute * 15).Unix(),
			Id:        "",
			IssuedAt:  issuedAtTime.Unix(),
			Issuer:    "",
			NotBefore: 0,
			Subject:   strconv.FormatInt(user.ID, 10),
		},
		UserName: user.Username,
	})

	if err != nil {
		return "", err
	}

	return accessToken, nil
}

func validateEmptyUserFieldsHook(ctx context.Context, exec boil.ContextExecutor, u *models.User) error {
	if u.Username == "" || u.FirstName == "" || u.LastName == "" || u.Email == "" || u.ContactNumber == "" {
		return errors.New("invalid input, empty fields not allowed")
	}
	return nil
}
