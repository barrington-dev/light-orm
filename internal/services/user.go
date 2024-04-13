package services

import (
	"context"
	"database/sql"
	"errors"
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
