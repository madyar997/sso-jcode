package postgres

import (
	"context"
	"github.com/madyar997/sso-jcode/internal/entity"
	"github.com/opentracing/opentracing-go"
)

type SearchQuery struct {
	Limit     int
	Offset    int
	SortBy    string
	SortOrder string
}

func (ur *Postgres) GetUsers(ctx context.Context) (users []*entity.User, err error) {
	res := ur.client.WithContext(ctx).Find(&users)
	if res.Error != nil {
		return nil, res.Error
	}
	return users, nil
}

func (ur *Postgres) CreateUser(ctx context.Context, user *entity.User) (int, error) {
	res := ur.client.WithContext(ctx).Create(user)
	if res.Error != nil {
		return 0, res.Error
	}
	return user.Id, nil
}

func (ur *Postgres) GetUserByEmail(ctx context.Context, email string) (user *entity.User, err error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "get user by email repo")
	defer span.Finish()

	res := ur.client.Where("email = ?", email).WithContext(ctx).Find(&user)
	if res.Error != nil {
		return nil, res.Error
	}
	return user, nil
}

func (ur *Postgres) GetUserByID(ctx context.Context, id int) (user *entity.User, err error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "get user by id - repo")
	defer span.Finish()

	res := ur.client.WithContext(ctx).Where("id = ?", id).Find(&user)
	if res.Error != nil {
		return nil, res.Error
	}
	return user, nil
}
