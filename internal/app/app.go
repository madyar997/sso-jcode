// Package app configures and runs application.
package app

import (
	"fmt"
	"github.com/madyar997/practice_7/internal/entity"
	"github.com/madyar997/practice_7/pkg/cache"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/madyar997/practice_7/config"
	v1 "github.com/madyar997/practice_7/internal/controller/http/v1"
	"github.com/madyar997/practice_7/internal/usecase"
	"github.com/madyar997/practice_7/internal/usecase/repo"
	"github.com/madyar997/practice_7/pkg/httpserver"
	"github.com/madyar997/practice_7/pkg/logger"
	"github.com/madyar997/practice_7/pkg/postgres"
)

// Run creates objects via constructors.
func Run(cfg *config.Config) {
	l := logger.New(cfg.Log.Level)

	// Repository
	pg, err := postgres.New(cfg.PG.URL)
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - postgres.New: %w", err))
	}
	defer pg.Close()

	err = pg.DB.AutoMigrate(entity.User{}, entity.Token{})
	if err != nil {
		log.Fatalf("could not auto migrate: %s", err.Error())
	}

	redisClient, err := cache.NewRedisClient()
	if err != nil {
		return
	}

	userCache := cache.NewUserCache(redisClient, 10*time.Minute)

	userUseCase := usecase.NewUser(repo.NewUserRepo(pg))

	// HTTP Server
	handler := gin.New()
	v1.NewRouter(handler, l, userUseCase, userCache, cfg)
	httpServer := httpserver.New(handler, httpserver.Port(cfg.HTTP.Port))

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("app - Run - signal: " + s.String())
	case err = <-httpServer.Notify():
		l.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	}

	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}
}
