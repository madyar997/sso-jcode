// Package v1 implements routing paths. Each services in own file.
package v1

import (
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/madyar997/sso-jcode/config"
	"github.com/madyar997/sso-jcode/pkg/cache"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	_ "github.com/santosh/gingo/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"

	// Swagger docs.
	"github.com/madyar997/sso-jcode/internal/usecase"
	"github.com/madyar997/sso-jcode/pkg/logger"
)

// NewRouter -.
// Swagger spec:
// @title       SSO API
// @description single sign on service for users
// @version     1.0
// @host        localhost:8080
// @BasePath    /api/v1
func NewRouter(handler *gin.Engine, l *logger.Logger, u usecase.UserUseCase, uc cache.User, cfg *config.Config) {
	// Options
	handler.Use(gin.Recovery())

	pprof.Register(handler)
	handler.Static("/assets", "./docs")
	//handler.StaticFS("/files/*any", gin.Dir("sso-jcode", true))
	// Swagger
	//handler.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler,
	//	ginSwagger.URL("http://localhost:8080/docs/swagger.json"),
	//	ginSwagger.DefaultModelsExpandDepth(-1)))
	handler.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// K8s probe
	handler.GET("/healthz", func(c *gin.Context) { c.Status(http.StatusOK) })

	// Prometheus metrics
	handler.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Routers
	h := handler.Group("/api/v1")
	{
		newUserRoutes(h, u, l, uc, cfg)
	}
}
