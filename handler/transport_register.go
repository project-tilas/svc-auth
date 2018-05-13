package handler

import (
	"context"
	"encoding/json"
	"net/http"
)

func decodeRegisterRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	return req, nil
}
