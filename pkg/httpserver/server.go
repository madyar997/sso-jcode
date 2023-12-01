// Package httpserver implements HTTP server.
package httpserver

import (
	"context"
	"github.com/madyar997/sso-jcode/config"
	"github.com/madyar997/sso-jcode/internal/usecase"
	"github.com/pkg/errors"
	"log"
	"net/http"
	"time"
)

const (
	_defaultReadTimeout  = 5 * time.Second
	_defaultWriteTimeout = 5 * time.Second
)

// Server -.
type Server struct {
	Address         string
	UserService     usecase.User
	Handler         http.Handler
	idleConnsClosed chan struct{}
	masterCtx       context.Context
}

// New -.
func New(ctx context.Context, cfg *config.Config, handler http.Handler) *Server {
	httpSrv := &Server{
		Address:         cfg.HTTP.Port,
		idleConnsClosed: make(chan struct{}),
		masterCtx:       ctx,
		Handler:         handler,
	}

	return httpSrv
}

func (s *Server) Run() error {

	srv := &http.Server{
		Addr:         s.Address,
		Handler:      s.Handler,
		ReadTimeout:  _defaultReadTimeout,
		WriteTimeout: _defaultWriteTimeout,
	}

	go s.GracefulShutdown(srv)
	log.Printf("[INFO] serving HTTP on \"%s\"", s.Address)

	if err := srv.ListenAndServe(); err != nil {
		return errors.WithMessage(err, "error when starting the http server")
	}

	return nil
}

func (srv *Server) GracefulShutdown(httpSrv *http.Server) {
	<-srv.masterCtx.Done()

	if err := httpSrv.Shutdown(context.Background()); err != nil {
		log.Printf("[ERROR] HTTP server Shutdown: %v", err)
	}

	log.Println("[INFO] HTTP server has processed all idle connections")
	close(srv.idleConnsClosed)
}

func (srv *Server) Wait() {
	<-srv.idleConnsClosed
}
