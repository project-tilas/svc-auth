package handler

import (
	"github.com/go-kit/kit/endpoint"
)

type AuthEndpoints struct {
	LoginEndpoint    endpoint.Endpoint
	RegisterEndpoint endpoint.Endpoint
}
