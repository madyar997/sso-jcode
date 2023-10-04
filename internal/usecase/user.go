package usecase

import (
	"context"
	"github.com/evrone/go-clean-template/internal/entity"
)

type User struct {
	repo UserRepo
}

func NewUser(repo UserRepo) *User {
	return &User{repo: repo}
}

func (u *User) Users(ctx context.Context) ([]*entity.User, error) {
	return u.repo.GetUsers(ctx)
}
