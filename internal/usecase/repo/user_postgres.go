package repo

import (
	"context"
	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/evrone/go-clean-template/pkg/postgres"
)

type UserRepo struct {
	*postgres.Postgres
}

func NewUserRepo(pg *postgres.Postgres) *UserRepo {
	return &UserRepo{pg}
}

func (ur *UserRepo) GetUsers(ctx context.Context) (users []*entity.User, err error) {
	res := ur.DB.WithContext(ctx).Find(&users)
	if res.Error != nil {
		return nil, res.Error
	}
	return users, nil
}

func (ur *UserRepo) CreateUser(ctx context.Context, user *entity.User) (int, error) {
	res := ur.DB.WithContext(ctx).Create(user)
	if res.Error != nil {
		return 0, res.Error
	}
	return user.Id, nil
}

func (ur *UserRepo) GetUserByEmail(ctx context.Context, email string) (user *entity.User, err error) {
	res := ur.DB.Where("email = ?", email).WithContext(ctx).Find(&user)
	if res.Error != nil {
		return nil, res.Error
	}
	return user, nil
}

func (ur *UserRepo) GetUserByID(ctx context.Context, id string) (user *entity.User, err error) {
	res := ur.DB.WithContext(ctx).Where("id = ?", id).Find(&user)
	if res.Error != nil {
		return nil, res.Error
	}
	return user, nil
}
