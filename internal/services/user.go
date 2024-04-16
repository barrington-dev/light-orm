package services

import (
	"context"
	"database/sql"
	"errors"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"light-orm/internal/database"
	"light-orm/internal/models"

	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type UserService struct {
	Db database.Service
}

func NewUserService(db database.Service) *UserService {
	return &UserService{
		Db: db,
	}
}

func (us *UserService) GetUser(ctx context.Context, id int64) (*models.User, error) {
	user, err := models.Users(qm.Where("id=?", id)).One(ctx, us.Db.Instance())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return user, nil
}

func (us *UserService) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	err := user.Insert(ctx, us.Db.Instance(), boil.Infer())
	if err != nil {
		return nil, err
	}

	return user, nil
}
