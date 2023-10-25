package usecase

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/jackc/pgx/v4"
	"github.com/madyar997/practice_7/config"
	"github.com/madyar997/practice_7/internal/controller/http/v1/dto"
	"github.com/madyar997/practice_7/internal/entity"
	"golang.org/x/crypto/bcrypt"
	"time"
)

const AccessTokenTTL = 900
const RefreshTokenTTL = 1800

type User struct {
	cfg  *config.Config
	repo UserRepo
}

func NewUser(repo UserRepo, cfg *config.Config) *User {
	return &User{repo: repo, cfg: cfg}
}

func (u *User) GetUserByID(ctx context.Context, id int) (*entity.User, error) {
	return u.repo.GetUserByID(ctx, id)
}

func (u *User) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	return u.repo.GetUserByEmail(ctx, email)
}

func (u *User) Users(ctx context.Context) ([]*entity.User, error) {
	return u.repo.GetUsers(ctx)
}

func (u *User) CreateUser(ctx context.Context, user *entity.User) (int, error) {
	return u.repo.CreateUser(ctx, user)
}

func (u *User) Register(ctx context.Context, email, password string) error {
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
	user, err := u.repo.GetUserByEmail(ctx, email)
	switch {
	case err == nil:
	case errors.Is(err, pgx.ErrNoRows):
		return nil, errors.New("user is not exist")
	default:
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New(fmt.Sprintf("passwords do not match %v", err))
	}

	accessTokenClaims := jwt.MapClaims{
		"user_id":   user.Id,
		"email":     user.Email,
		"name":      user.Name,
		"ExpiresAt": time.Now().Add(time.Hour * 1).Unix(),
	}

	accessToken := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), accessTokenClaims)

	accessTokenString, err := accessToken.SignedString([]byte(u.cfg.SecretKey))
	if err != nil {
		return nil, err
	}

	refreshTokenClaims := jwt.MapClaims{
		"user_id":   user.Id,
		"ExpiresAt": time.Now().Add(time.Hour * 1),
	}

	refreshToken := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), refreshTokenClaims)

	resfreshTokenString, err := refreshToken.SignedString([]byte(u.cfg.SecretKey))
	if err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		Name:         user.Name,
		Email:        user.Email,
		AccessToken:  accessTokenString,
		RefreshToken: resfreshTokenString,
	}, nil
}
