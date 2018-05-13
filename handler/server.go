package handler

import (
	"net/http"

	"github.com/go-kit/kit/endpoint"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/project-tilas/svc-auth/service"
)

func makeAuthServerEndpoints(s service.AuthService) AuthEndpoints {
	return AuthEndpoints{
		LoginEndpoint:    makeLoginEndpoint(s),
		RegisterEndpoint: makeRegisterEndpoint(s),
	}
}

type authHandlers struct {
	Login    http.Handler
	Register http.Handler
}

func makeAuthKitHandlers(s AuthEndpoints, opts []kithttp.ServerOption) *authHandlers {
	newServer := func(e endpoint.Endpoint, decodeFn kithttp.DecodeRequestFunc) http.Handler {
		return kithttp.NewServer(e, decodeFn, encodeResponse, opts...)
	}

	return &authHandlers{
		Login:    newServer(s.LoginEndpoint, decodeLoginRequest),
		Register: newServer(s.RegisterEndpoint, decodeRegisterRequest),
	}
}

func attachAuthKitRoutes(r *mux.Router, h *authHandlers) {
	r.Handle("/login", h.Login).Methods("POST").Name("login")
	r.Handle("/register", h.Register).Methods("POST").Name("login")
}

// MakeServerHandler creates an HTTP handler all of the endpoints.
func MakeServerHandler(svc service.AuthService, r *mux.Router) http.Handler {
	kolideAPIOptions := []kithttp.ServerOption{
		kithttp.ServerBefore(
			kithttp.PopulateRequestContext, // populate the request context with common fields
		),
		kithttp.ServerErrorEncoder(encodeError),
		kithttp.ServerAfter(
			kithttp.SetContentType("application/json; charset=utf-8"),
		),
	}

	authEndpoints := makeAuthServerEndpoints(svc)
	authHandlers := makeAuthKitHandlers(authEndpoints, kolideAPIOptions)

	attachAuthKitRoutes(r, authHandlers)
	return r
}
