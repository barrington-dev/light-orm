package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"light-orm/internal/database"
	"light-orm/internal/models"
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

func (us *UserService) GetUser(ctx context.Context, id int64) (*models.User, error) {
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

func (us *UserService) CreateUser(ctx context.Context, user *models.User, password *models.Password) (*models.User, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	err := user.Insert(ctx, us.Db.Instance(), boil.Infer())
	if err != nil {
		return nil, err
	}

	// Hash password
	hash, err := NewAuthService().HashPassword(password.Password)
	if err != nil {
		return nil, err
	}

	// Update password with hash
	password.Password = hash

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

func validateEmptyUserFieldsHook(ctx context.Context, exec boil.ContextExecutor, u *models.User) error {
	if u.Username == "" || u.FirstName == "" || u.LastName == "" || u.Email == "" || u.ContactNumber == "" {
		return errors.New("invalid input, empty fields not allowed")
	}
	return nil
}
