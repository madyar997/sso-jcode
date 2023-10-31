package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/madyar997/sso-jcode/config"
	"github.com/madyar997/sso-jcode/internal/controller/http/v1/dto"
	"github.com/madyar997/sso-jcode/internal/entity"
	"github.com/madyar997/sso-jcode/internal/usecase"
	"github.com/madyar997/sso-jcode/pkg/cache"
	"github.com/madyar997/sso-jcode/pkg/jaeger"
	"github.com/madyar997/sso-jcode/pkg/logger"
	"github.com/opentracing/opentracing-go"
	"log"
	"net/http"
	"strconv"
	"time"
)

type userRoutes struct {
	u         usecase.UserUseCase
	l         logger.Interface
	userCache cache.User
	cfg       *config.Config
}

func newUserRoutes(handler *gin.RouterGroup, u usecase.UserUseCase, l logger.Interface, uc cache.User, cfg *config.Config) {
	r := &userRoutes{u, l, uc, cfg}

	adminHandler := handler.Group("/admin/user")
	{
		adminHandler.GET("/:id", r.GetUserByID)
		adminHandler.GET("/all", r.GetUsers)
		adminHandler.POST("/", r.CreateUser)
		adminHandler.GET("/", r.GetUserByEmail)
	}

	userHandler := handler.Group("/user")
	{
		userHandler.POST("/register", r.Register)
		userHandler.POST("/login", r.Login)
		userHandler.POST("/refresh", r.Refresh)
	}
}

func (ur *userRoutes) GetUsers(ctx *gin.Context) {
	users, err := ur.u.Users(ctx)
	if err != nil {
		ur.l.Error(err, "http - v1 - user - all")
		errorResponse(ctx, http.StatusInternalServerError, "database problems")

		return
	}

	ctx.JSON(http.StatusOK, users)
}

func (ur *userRoutes) CreateUser(ctx *gin.Context) {
	var user *entity.User

	err := ctx.ShouldBindJSON(&user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)

		return
	}

	insertedID, err := ur.u.CreateUser(ctx, user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)

		return
	}

	ctx.JSON(http.StatusOK, insertedID)
}

func (ur *userRoutes) Register(ctx *gin.Context) {
	var registerRequest dto.RegisterRequest

	err := ctx.ShouldBindJSON(&registerRequest)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)

		return
	}

	err = ur.u.Register(ctx, registerRequest.Email, registerRequest.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)

		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "user successfully registered"})
}

func (ur *userRoutes) Login(ctx *gin.Context) {
	span := opentracing.StartSpan("auth service, handler /login")
	defer span.Finish()

	context := opentracing.ContextWithSpan(ctx.Request.Context(), span)

	var loginRequest dto.LoginRequest

	err := ctx.ShouldBindJSON(&loginRequest)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)

		return
	}

	token, err := ur.u.Login(context, loginRequest.Email, loginRequest.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)

		return
	}

	ctx.SetCookie("access_token", token.AccessToken, 3600, "/", "localhost", false, true)
	ctx.SetCookie("refresh_token", token.RefreshToken, 3600, "/", "localhost", false, true)

	ctx.JSON(http.StatusOK, token)
}

func (ur *userRoutes) GetUserByEmail(ctx *gin.Context) {
	email := ctx.Query("email")

	user, err := ur.userCache.Get(ctx, email)
	if err != nil {

		return
	}

	if user == nil {
		user, err = ur.u.GetUserByEmail(ctx, email)
		if err != nil {
			ur.l.Error(err, "http - v1 - user - all")
			errorResponse(ctx, http.StatusInternalServerError, "database problems")

			return
		}

		err = ur.userCache.Set(ctx, email, user)
		if err != nil {
			log.Printf("could not cache user with email %s: %v", email, err)
		}
	}

	ctx.JSON(http.StatusOK, user)
}

func (ur *userRoutes) Refresh(ctx *gin.Context) {
	userID, ok := ctx.Get("user_id")
	if !ok {
		ctx.JSON(http.StatusBadRequest, fmt.Errorf("could not get user id from token"))

		return
	}

	user, err := ur.u.GetUserByID(ctx, int(userID.(float64)))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)

		return
	}
	//не изменились ли роли
	if user == nil {
		ctx.JSON(http.StatusInternalServerError, err)

		return
	}

	accessTokenClaims := jwt.MapClaims{
		"user_id": user.Id,
		"email":   user.Email,
		"name":    user.Name,
		"exp":     time.Now().Add(time.Second * usecase.AccessTokenTTL).Unix(),
	}

	accessToken := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), accessTokenClaims)

	accessTokenString, err := accessToken.SignedString([]byte(ur.cfg.SecretKey))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)

		return
	}

	refreshTokenClaims := jwt.MapClaims{
		"user_id": user.Id,
		"exp":     time.Now().Add(time.Second * usecase.RefreshTokenTTL),
	}

	refreshToken := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), refreshTokenClaims)

	refreshTokenString, err := refreshToken.SignedString([]byte(ur.cfg.SecretKey))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)

		return
	}

	ctx.JSON(http.StatusOK, dto.LoginResponse{
		Name:         user.Name,
		Email:        user.Email,
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
	})
}

func (ur *userRoutes) GetUserByID(ctx *gin.Context) {
	span := jaeger.StartSpanFromRequest(jaeger.Tracer, ctx.Request, "sso /getUserByID handler method")
	defer span.Finish()

	idQueryParam := ctx.Param("id")

	span.LogKV("id", idQueryParam)

	id, err := strconv.Atoi(idQueryParam)
	if err != nil {
		ur.l.Error(err, "http - v1 - user - get by id ")
		errorResponse(ctx, http.StatusBadRequest, "id is incorrect")

		return
	}

	context := opentracing.ContextWithSpan(ctx.Request.Context(), span)

	user, err := ur.u.GetUserByID(context, id)
	if err != nil {
		ur.l.Error(err, "http - v1 - user - all")
		errorResponse(ctx, http.StatusInternalServerError, "database problems")

		return
	}

	userDto := dto.UserInfo{
		Id:    user.Id,
		Name:  user.Name,
		Email: user.Email,
		Age:   user.Age,
	}

	ctx.JSON(http.StatusOK, userDto)
}
