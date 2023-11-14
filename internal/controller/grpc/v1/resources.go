package v1

import (
	"context"
	"github.com/madyar997/sso-jcode/internal/usecase"
	"github.com/madyar997/user-client/protobuf"
)

type UserServiceResources struct {
	protobuf.UnimplementedUserServer
	userUseCase usecase.UserUseCase
}

func NewUserServiceResource(userUseCase usecase.UserUseCase) *UserServiceResources {
	return &UserServiceResources{
		userUseCase: userUseCase,
	}
}

func (us *UserServiceResources) GetUserByID(ctx context.Context, req *protobuf.UserRequest) (*protobuf.UserResponse, error) {
	user, err := us.userUseCase.GetUserByID(ctx, int(req.Id))
	if err != nil {
		return nil, err
	}

	return &protobuf.UserResponse{
		Id:    int32(user.Id),
		Name:  user.Name,
		Email: user.Email,
		Age:   int32(user.Age),
	}, nil
}
