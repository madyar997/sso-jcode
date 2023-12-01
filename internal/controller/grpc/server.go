package grpc

import (
	"context"
	"github.com/madyar997/sso-jcode/config"
	v1 "github.com/madyar997/sso-jcode/internal/controller/grpc/v1"
	"github.com/madyar997/sso-jcode/internal/usecase"
	"github.com/madyar997/user-client/protobuf"
	"google.golang.org/grpc"
	"log"
	"net"
)

type GrpcServer struct {
	cfg     *config.Config
	Address string
	server  *grpc.Server

	idleConnsClosed chan struct{}
	masterCtx       context.Context

	userUseCase usecase.UserUseCase
}

func NewGrpcServer(ctx context.Context, address string, userUseCase usecase.UserUseCase, cfg *config.Config) *GrpcServer {
	return &GrpcServer{
		Address:         address,
		userUseCase:     userUseCase,
		cfg:             cfg,
		idleConnsClosed: make(chan struct{}),
		masterCtx:       ctx,
	}
}

func (gs *GrpcServer) Run() error {
	lis, err := net.Listen("tcp", gs.cfg.Grpc.Port)
	if err != nil {
		return err
	}

	gs.server = grpc.NewServer()

	resource := v1.NewUserServiceResource(gs.userUseCase)
	protobuf.RegisterUserServer(gs.server, resource)

	go gs.GracefulShutdown(gs.server)

	log.Printf("[INFO] serving GRPC on \"%s\"", gs.Address)
	if err = gs.server.Serve(lis); err != nil {
		return err
	}
	return nil
}

func (gs *GrpcServer) GracefulShutdown(grpcServer *grpc.Server) {
	<-gs.masterCtx.Done()
	log.Printf("[INFO] shutting down gRPC server")
	grpcServer.GracefulStop()
	close(gs.idleConnsClosed)
}

func (gs *GrpcServer) Wait() {
	<-gs.idleConnsClosed
	log.Println("[INFO] gRPC server has processed all idle connections")
}
