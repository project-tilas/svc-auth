package handler

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/project-tilas/svc-auth/domain"
	"github.com/project-tilas/svc-auth/service"
)

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type loginResponse struct {
	User  *domain.User  `json:"user,omitempty"`
	Token *domain.Token `json:"token,omitempty"`
	Err   error         `json:"error,omitempty"`
}

func (r *loginResponse) error() error { return r.Err }

func makeLoginEndpoint(s service.AuthService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(loginRequest)
		user, token, err := s.Login(ctx, req.Username, req.Password)
		if err != nil {
			return loginResponse{Err: err}, nil
		}
		return loginResponse{User: user, Token: token}, nil
	}
}
