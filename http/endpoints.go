package http

import (
	"github.com/go-kit/kit/endpoint"
	"github.com/project-tilas/svc-auth/domain"
)

type AuthEndpoints struct {
	LoginEndpoint    endpoint.Endpoint
	RegisterEndpoint endpoint.Endpoint
}

// REGISTER

type registerRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type registerResponse struct {
	User  domain.User  `json:"user,omitempty"`
	Token domain.Token `json:"token,omitempty"`
	Err   error        `json:"error,omitempty"`
}

// LOGIN

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type loginResponse struct {
	User  domain.User  `json:"user,omitempty"`
	Token domain.Token `json:"token,omitempty"`
	Err   error        `json:"error,omitempty"`
}
