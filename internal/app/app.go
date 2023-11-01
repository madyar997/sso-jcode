// Package app configures and runs application.
package app

import (
	"github.com/gin-gonic/gin"
	"github.com/madyar997/sso-jcode/internal/entity"
	"github.com/madyar997/sso-jcode/pkg/cache"
	"github.com/madyar997/sso-jcode/pkg/jaeger"
	"github.com/madyar997/sso-jcode/pkg/logger"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/madyar997/sso-jcode/config"
	v1 "github.com/madyar997/sso-jcode/internal/controller/http/v1"
	"github.com/madyar997/sso-jcode/internal/usecase"
	"github.com/madyar997/sso-jcode/internal/usecase/repo"
	"github.com/madyar997/sso-jcode/pkg/httpserver"
	"github.com/madyar997/sso-jcode/pkg/postgres"
)

// Run creates objects via constructors.
func Run(cfg *config.Config) {
	//l := logger.New(cfg.Log.Level)

	l := logger.New()

	//tracing
	tracer, closer, _ := jaeger.InitJaeger()
	defer closer.Close()
	opentracing.SetGlobalTracer(tracer)

	// Repository
	pg, err := postgres.New(cfg.PG.URL)
	if err != nil {
		l.Logger.Fatal("app - Run - postgres.New: %w", zap.Error(err))
	}
	defer pg.Close()

	err = pg.DB.AutoMigrate(entity.User{})
	if err != nil {
		log.Fatalf("could not auto migrate: %s", err.Error())
	}

	redisClient, err := cache.NewRedisClient()
	if err != nil {
		return
	}

	userCache := cache.NewUserCache(redisClient, cache.UserCacheTimeout)

	userUseCase := usecase.NewUser(repo.NewUserRepo(pg), cfg, l)

	// HTTP Server
	handler := gin.New()
	v1.NewRouter(handler, l, userUseCase, userCache, cfg)
	httpServer := httpserver.New(handler, httpserver.Port(cfg.HTTP.Port))

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Logger.Fatal("app - Run - signal: " + s.String())
	case err = <-httpServer.Notify():
		l.Logger.Fatal("app - Run - httpServer.Notify")
	}

	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		l.Logger.Fatal("app - Run - httpServer.Shutdown: %w", zap.Error(err))
	}
}
