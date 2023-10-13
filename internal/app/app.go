// Package app configures and runs application.
package app

import (
	"fmt"
	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/evrone/go-clean-template/pkg/cache"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/evrone/go-clean-template/config"
	v1 "github.com/evrone/go-clean-template/internal/controller/http/v1"
	"github.com/evrone/go-clean-template/internal/usecase"
	"github.com/evrone/go-clean-template/internal/usecase/repo"
	"github.com/evrone/go-clean-template/pkg/httpserver"
	"github.com/evrone/go-clean-template/pkg/logger"
	"github.com/evrone/go-clean-template/pkg/postgres"
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
	v1.NewRouter(handler, l, userUseCase, userCache)
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
