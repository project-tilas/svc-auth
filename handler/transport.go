package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/project-tilas/svc-auth/domain"
)

// errorer is implemented by all concrete response types that may contain
// errors. It allows us to change the HTTP response code without needing to
// trigger an endpoint (transport-level) error. For more information, read the
// big comment in endpoints.go.
type errorer interface {
	error() error
}

// statuser allows response types to implement a custom
// http success status - default is 200 OK
type statuser interface {
	status() int
}

// encodeResponse is the common method to encode all response types to the
// client. I chose to do it this way because, since we're using JSON, there's no
// reason to provide anything more specific. It's certainly possible to
// specialize on a per-response (per-method) basis.
func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if e, ok := response.(errorer); ok && e.error() != nil {
		// Not a Go kit transport error, but a business-logic error.
		// Provide those as HTTP errors.
		encodeError(ctx, e.error(), w)
		return nil
	}

	if e, ok := response.(statuser); ok {
		w.WriteHeader(e.status())
		if e.status() == http.StatusNoContent {
			return nil
		}
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(response)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	if err == nil {
		panic("encodeError with nil error")
	}

	w.WriteHeader(codeFrom(err))
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}

func codeFrom(err error) int {
	switch err.(type) {
	case *domain.NotFoundError:
		return http.StatusNotFound
	case *domain.AlreadyExistsError:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
