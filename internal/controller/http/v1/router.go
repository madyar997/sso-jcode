// Package v1 implements routing paths. Each services in own file.
package v1

import (
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/madyar997/sso-jcode/config"
	"github.com/madyar997/sso-jcode/pkg/cache"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"

	// Swagger docs.
	_ "github.com/madyar997/sso-jcode/docs"
	"github.com/madyar997/sso-jcode/internal/usecase"
	"github.com/madyar997/sso-jcode/pkg/logger"
)

// NewRouter -.
// Swagger spec:
// @title       Go Clean Template API
// @description Using a translation service as an example
// @version     1.0
// @host        localhost:8080
// @BasePath    /v1
func NewRouter(handler *gin.Engine, l *logger.Logger, u usecase.UserUseCase, uc cache.User, cfg *config.Config) {
	// Options
	handler.Use(gin.Recovery())

	pprof.Register(handler)

	// Swagger
	swaggerHandler := ginSwagger.DisablingWrapHandler(swaggerFiles.Handler, "DISABLE_SWAGGER_HTTP_HANDLER")
	handler.GET("/swagger/*any", swaggerHandler)

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
