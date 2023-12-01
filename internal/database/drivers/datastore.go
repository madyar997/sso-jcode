package drivers

import (
	"context"
	"github.com/madyar997/sso-jcode/internal/entity"
)

type DataStore interface {
	Name() string
	Close() error
	Connect() error
	UserRepo
}

type UserRepo interface {
	GetUsers(ctx context.Context) ([]*entity.User, error)
	GetUserByID(ctx context.Context, id int) (user *entity.User, err error)
	CreateUser(ctx context.Context, user *entity.User) (int, error)
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)
}
