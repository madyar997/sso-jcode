package v1

import (
	"github.com/evrone/go-clean-template/internal/controller/http/v1/dto"
	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/evrone/go-clean-template/internal/usecase"
	"github.com/evrone/go-clean-template/pkg/cache"
	"github.com/evrone/go-clean-template/pkg/logger"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type userRoutes struct {
	u         usecase.UserUseCase
	l         logger.Interface
	userCache cache.User
}

func newUserRoutes(handler *gin.RouterGroup, u usecase.UserUseCase, l logger.Interface, uc cache.User) {
	r := &userRoutes{u, l, uc}

	adminHandler := handler.Group("/admin/user")
	{
		//adminHandler.Use(middleware.CustomLogger())
		//adminHandler.Use(middleware.JwtVerify())
		adminHandler.GET("/all", r.GetUsers)
		adminHandler.POST("/", r.CreateUser)
		adminHandler.GET("/", r.GetUserByEmail)
	}

	userHandler := handler.Group("/user")
	{
		userHandler.POST("/register", r.Register)
		userHandler.POST("/login", r.Login)
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
	var loginRequest dto.LoginRequest

	err := ctx.ShouldBindJSON(&loginRequest)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	token, err := ur.u.Login(ctx, loginRequest.Email, loginRequest.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

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
