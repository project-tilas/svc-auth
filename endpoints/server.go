package endpoints

import (
	"context"

	"github.com/project-tilas/svc-auth/domain"

	"github.com/go-kit/kit/endpoint"
	"github.com/project-tilas/svc-auth/service"
)

func MakeServerEndpoints(s service.AuthService) Endpoints {
	return Endpoints{
		Login:    MakeLoginEndpoint(s),
		Register: MakeRegisterEndpoint(s),
	}
}

func MakeLoginEndpoint(s service.AuthService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(loginRequest)
		user, token, e := s.Login(ctx, req.Username, req.Password)
		return loginResponse{User: user, Token: token, Err: e}, nil
	}
}

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type loginResponse struct {
	User  domain.User  `json:"user,omitempty"`
	Token domain.Token `json:"token,omitempty"`
	Err   error        `json:"error,omitempty"`
}

func MakeRegisterEndpoint(s service.AuthService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(registerRequest)
		user, token, e := s.Register(ctx, domain.User{
			Username: req.Username,
			Password: req.Password,
		})
		return registerResponse{User: user, Token: token, Err: e}, nil
	}
}

type registerRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type registerResponse struct {
	User  domain.User  `json:"user,omitempty"`
	Token domain.Token `json:"token,omitempty"`
	Err   error        `json:"error,omitempty"`
}
