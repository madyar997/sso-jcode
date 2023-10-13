package usecase

import (
	"context"
	"errors"
	"fmt"
	"github.com/evrone/go-clean-template/internal/controller/http/v1/dto"
	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/golang-jwt/jwt"
	"github.com/jackc/pgx/v4"
	"golang.org/x/crypto/bcrypt"
	"time"
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

func (u *User) CreateUser(ctx context.Context, user *entity.User) (int, error) {
	return u.repo.CreateUser(ctx, user)
}

func (u *User) Register(ctx context.Context, email, password string) error {
	//email password
	//is email exists return with message "go to login"

	generatedHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_, err = u.repo.CreateUser(ctx, &entity.User{
		Email:    email,
		Password: string(generatedHash),
	})
	if err != nil {
		return err
	}

	return nil
}

func (u *User) Login(ctx context.Context, email, password string) (*dto.LoginResponse, error) {
	//есть ли такой аккаунт с email =  email
	user, err := u.repo.GetUserByEmail(ctx, email)
	switch {
	case err == nil:
	case err == pgx.ErrNoRows:
		return nil, errors.New("user is not exist")
	default:
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New(fmt.Sprintf("passwords do not match %v", err))
	}

	expiresAt := time.Now().Add(time.Hour * 1).Unix()

	tk := &entity.Token{
		Name:  user.Name,
		Email: user.Email,
		StandardClaims: &jwt.StandardClaims{
			Audience:  user.Name,
			ExpiresAt: expiresAt,
		},
	}

	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk.StandardClaims)

	tokenString, err := token.SignedString([]byte("practice_7"))
	if err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		Name:  user.Name,
		Email: user.Email,
		Token: tokenString,
	}, nil
}

func (u *User) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	time.Sleep(2 * time.Second)

	return u.repo.GetUserByEmail(ctx, email)
}
