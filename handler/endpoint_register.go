package handler

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/project-tilas/svc-auth/domain"
	"github.com/project-tilas/svc-auth/service"
)

type registerRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type registerResponse struct {
	User  *domain.User  `json:"user,omitempty"`
	Token *domain.Token `json:"token,omitempty"`
	Err   error         `json:"error,omitempty"`
}

func (r *registerResponse) error() error { return r.Err }

func makeRegisterEndpoint(s service.AuthService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(registerRequest)
		user, token, err := s.Register(ctx, domain.User{
			Username: req.Username,
			Password: req.Password,
		})
		if err != nil {
			return registerResponse{Err: err}, nil
		}
		return registerResponse{User: user, Token: token}, nil
	}
}
