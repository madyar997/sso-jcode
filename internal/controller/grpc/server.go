package grpc

import (
	"github.com/madyar997/sso-jcode/config"
	v1 "github.com/madyar997/sso-jcode/internal/controller/grpc/v1"
	"github.com/madyar997/sso-jcode/internal/usecase"
	"github.com/madyar997/user-client/protobuf"
	"google.golang.org/grpc"
	"net"
)

type GrpcServer struct {
	cfg     *config.Config
	Address string
	server  *grpc.Server

	userUseCase usecase.UserUseCase
}

func NewGrpcServer(address string, userUseCase usecase.UserUseCase, cfg *config.Config) *GrpcServer {
	return &GrpcServer{
		Address:     address,
		userUseCase: userUseCase,
		cfg:         cfg,
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

	return gs.server.Serve(lis)
}
