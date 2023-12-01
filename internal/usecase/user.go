package usecase

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/jackc/pgx/v4"
	"github.com/madyar997/sso-jcode/config"
	"github.com/madyar997/sso-jcode/internal/controller/http/v1/dto"
	"github.com/madyar997/sso-jcode/internal/database/drivers"
	"github.com/madyar997/sso-jcode/internal/entity"
	"github.com/madyar997/sso-jcode/pkg/logger"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"time"
)

const AccessTokenTTL = 900
const RefreshTokenTTL = 1800

type User struct {
	cfg    *config.Config
	repo   drivers.DataStore
	logger *logger.Logger
}

func NewUser(repo drivers.DataStore, cfg *config.Config, logger *logger.Logger) *User {
	return &User{repo: repo, cfg: cfg, logger: logger}
}

func (u *User) GetUserByID(ctx context.Context, id int) (*entity.User, error) {
	span, spanCtx := opentracing.StartSpanFromContext(ctx, "get user by id - use case")
	defer span.Finish()

	return u.repo.GetUserByID(spanCtx, id)
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
	span, spanCtx := opentracing.StartSpanFromContext(ctx, "login use case")
	defer span.Finish()

	user, err := u.repo.GetUserByEmail(spanCtx, email)
	switch {
	case err == nil:
	case errors.Is(err, pgx.ErrNoRows):
		u.logger.Warn("user not found", zap.Error(err))
		return nil, errors.New("user is not exist")
	default:
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		u.logger.Error("passwords not match", zap.Error(err))
		return nil, errors.New(fmt.Sprintf("passwords do not match %v", err))
	}

	u.logger.Info("generating access and refresh tokens ...")
	accessTokenClaims := jwt.MapClaims{
		"user_id": user.Id,
		"email":   user.Email,
		"name":    user.Name,
		"exp":     time.Now().Add(time.Hour * 1).Unix(),
	}

	accessToken := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), accessTokenClaims)

	accessTokenString, err := accessToken.SignedString([]byte(u.cfg.SecretKey))
	if err != nil {
		return nil, err
	}

	refreshTokenClaims := jwt.MapClaims{
		"user_id": user.Id,
		"exp":     time.Now().Add(time.Hour * 1),
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
