package v1

import (
	"github.com/evrone/go-clean-template/internal/usecase"
	"github.com/evrone/go-clean-template/pkg/logger"
	"github.com/gin-gonic/gin"
	"net/http"
)

type userRoutes struct {
	u usecase.UserUseCase
	l logger.Interface
}

func newUserRoutes(handler *gin.RouterGroup, u usecase.UserUseCase, l logger.Interface) {
	r := &userRoutes{u, l}

	h := handler.Group("/user")
	{
		h.GET("/all", r.GetUsers)
		//h.POST("/do-translate", r.doTranslate
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
