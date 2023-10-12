// Package usecase implements application business logic. Each logic group in own file.
package usecase

import (
	"context"
	"github.com/evrone/go-clean-template/internal/controller/http/v1/dto"

	"github.com/evrone/go-clean-template/internal/entity"
)

//go:generate mockgen -source=interfaces.go -destination=./mocks_test.go -package=usecase_test

type (

	// User
	UserUseCase interface {
		Users(ctx context.Context) ([]*entity.User, error)
		CreateUser(ctx context.Context, user *entity.User) (int, error)
		GetUserByEmail(ctx context.Context, id string) (*entity.User, error)

		Register(ctx context.Context, email, password string) error
		Login(ctx context.Context, email, password string) (*dto.LoginResponse, error)
	}

	//UserRepo
	UserRepo interface {
		GetUsers(ctx context.Context) ([]*entity.User, error)
		GetUserByID(ctx context.Context, id string) (user *entity.User, err error)
		CreateUser(ctx context.Context, user *entity.User) (int, error)

		GetUserByEmail(ctx context.Context, email string) (*entity.User, error)
	}
)
