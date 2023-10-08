package v1

import (
	"github.com/evrone/go-clean-template/internal/controller/http/middleware"
	"github.com/evrone/go-clean-template/internal/controller/http/v1/dto"
	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/evrone/go-clean-template/internal/usecase"
	"github.com/evrone/go-clean-template/pkg/logger"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type userRoutes struct {
	u usecase.UserUseCase
	l logger.Interface
}

func newUserRoutes(handler *gin.RouterGroup, u usecase.UserUseCase, l logger.Interface) {
	r := &userRoutes{u, l}

	adminHandler := handler.Group("/admin/user")
	{
		adminHandler.Use(middleware.CustomLogger())
		adminHandler.Use(middleware.JwtVerify())
		adminHandler.GET("/all", r.GetUsers)
		adminHandler.POST("/", r.CreateUser)
		adminHandler.GET("/test", func(ctx *gin.Context) {
			log.Println("hello from controller")
			ctx.JSON(http.StatusOK, "test")
		})
		adminHandler.GET("/test2", func(ctx *gin.Context) {
			log.Println("hello from controller2")
			ctx.JSON(http.StatusOK, "test")
		})
	}

	userHandler := handler.Group("/user")
	{
		userHandler.POST("/register", r.Register)
		userHandler.POST("/login", r.Login)
	}

}

func (u *userRoutes) GetUsers(ctx *gin.Context) {
	users, err := u.u.Users(ctx)
	if err != nil {
		u.l.Error(err, "http - v1 - user - all")
		errorResponse(ctx, http.StatusInternalServerError, "database problems")

		return
	}

	ctx.JSON(http.StatusOK, users)
}

func (u *userRoutes) CreateUser(ctx *gin.Context) {
	var user *entity.User

	err := ctx.ShouldBindJSON(&user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	insertedID, err := u.u.CreateUser(ctx, user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, insertedID)
}

func (u *userRoutes) Register(ctx *gin.Context) {
	var registerRequest dto.RegisterRequest

	err := ctx.ShouldBindJSON(&registerRequest)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	err = u.u.Register(ctx, registerRequest.Email, registerRequest.Password)
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
