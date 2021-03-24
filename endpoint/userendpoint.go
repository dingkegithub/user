package endpoint

import (
	"context"

	"github.com/dingkegithub/user/service"
	"github.com/go-kit/kit/endpoint"
)

type UserEndpoints struct {
	LoginEndpoint    endpoint.Endpoint
	RegisterEndpoint endpoint.Endpoint
}

func NewUserEndpoints(userService service.UserService) *UserEndpoints {
	return &UserEndpoints{
		LoginEndpoint:    MakeLoginEndpoint(userService),
		RegisterEndpoint: MakeRegisterEndpoint(userService),
	}
}

type LoginRequest struct {
	Email    string
	Password string
}

type LoginResponse struct {
	UserInfo *service.UserInfoDto `json:"user_info"`
}

func MakeLoginEndpoint(userService service.UserService) endpoint.Endpoint {
	//type Endpoint func(ctx context.Context, request interface{}) (response interface{}, err error)
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*LoginRequest)
		dto, err := userService.Login(ctx, req.Email, req.Password)
		return &LoginResponse{UserInfo: dto}, err
	}
}

type RegisterRequest struct {
	UserName string
	Email    string
	Password string
}

type RegisterResponse struct {
	UserInfo *service.UserInfoDto `json:"user_info"`
}

func MakeRegisterEndpoint(userService service.UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*RegisterRequest)
		vo := &service.RegisterUserVo{
			UserName: req.UserName,
			Email:    req.Email,
			Password: req.Password,
		}

		dto, err := userService.Register(ctx, vo)

		return &RegisterResponse{UserInfo: dto}, err
	}
}
