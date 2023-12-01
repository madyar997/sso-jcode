package mongo

import (
	"context"
	"github.com/madyar997/sso-jcode/internal/entity"
)

func (m *Mongo) GetUsers(ctx context.Context) ([]*entity.User, error) {
	//TODO implement me
	panic("implement me")
}

func (m *Mongo) GetUserByID(ctx context.Context, id int) (user *entity.User, err error) {
	//TODO implement me
	panic("implement me")
}

func (m *Mongo) CreateUser(ctx context.Context, user *entity.User) (int, error) {
	//TODO implement me
	panic("implement me")
}

func (m *Mongo) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	//TODO implement me
	panic("implement me")
}
