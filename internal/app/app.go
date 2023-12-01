// Package app configures and runs application.
package app

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/madyar997/sso-jcode/internal/controller/grpc"
	"github.com/madyar997/sso-jcode/internal/database"
	"github.com/madyar997/sso-jcode/pkg/cache"
	"github.com/madyar997/sso-jcode/pkg/jaeger"
	"github.com/madyar997/sso-jcode/pkg/logger"
	"github.com/opentracing/opentracing-go"
	"golang.org/x/sync/errgroup"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/madyar997/sso-jcode/config"
	v1 "github.com/madyar997/sso-jcode/internal/controller/http/v1"
	"github.com/madyar997/sso-jcode/internal/usecase"
	"github.com/madyar997/sso-jcode/pkg/httpserver"
)

// Run creates objects via constructors.
func Run(cfg *config.Config) {
	appCtx, appCtxCancel := context.WithCancel(context.Background())
	defer appCtxCancel()

	l := logger.New()

	//tracing
	tracer, closer, _ := jaeger.InitJaeger()
	defer closer.Close()
	opentracing.SetGlobalTracer(tracer)

	// Repository
	ds, err := database.Connect(map[string]string{
		"datastore": cfg.PG.Name,
		"url":       cfg.PG.URL,
	})
	if err != nil {
		log.Printf("[ERROR] cannot connect to datastore: %v", err)
		return
	}
	log.Printf("[INFO] connected to %s", ds.Name())
	defer ds.Close()

	redisClient, err := cache.NewRedisClient()
	if err != nil {
		return
	}

	userCache := cache.NewUserCache(redisClient, cache.UserCacheTimeout)
	userUseCase := usecase.NewUser(ds, cfg, l)

	go signalHandler(appCtxCancel)

	g, gCtx := errgroup.WithContext(appCtx)

	g.Go(func() error {
		handler := gin.New()
		v1.NewRouter(handler, l, userUseCase, userCache, cfg)
		httpServer := httpserver.New(gCtx, cfg, handler)

		err = httpServer.Run()
		if err != nil {
			return fmt.Errorf("HTTP server: %v", err)
		}

		httpServer.Wait()
		return nil
	})

	g.Go(func() error {
		grpcServer := grpc.NewGrpcServer(gCtx,
			cfg.Grpc.Port,
			userUseCase,
			cfg)

		err = grpcServer.Run()
		if err != nil {
			log.Fatalf(err.Error())
		}

		grpcServer.Wait()
		return nil
	})

	// Ждем пока все горутины не будут завершены
	if err = g.Wait(); err != nil {
		log.Printf("[INFO] process terminated, %s", err)
		return
	}
}

// signalHandler обработка сигнала SIGTERM с остановкой контекста
func signalHandler(cancel context.CancelFunc) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	<-stop
	log.Print("[WARN] interrupt signal")
	cancel()
}
