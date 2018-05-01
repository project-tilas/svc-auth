package http

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	httptransport "github.com/go-kit/kit/transport/http"
)

func MakeClientEndpoints(instance string) (AuthEndpoints, error) {
	if !strings.HasPrefix(instance, "http") {
		instance = "http://" + instance
	}
	tgt, err := url.Parse(instance)
	if err != nil {
		return AuthEndpoints{}, err
	}
	tgt.Path = ""

	options := []httptransport.ClientOption{}
	return AuthEndpoints{
		LoginEndpoint:    httptransport.NewClient("POST", tgt, encodeLoginRequest, decodeLoginResponse, options...).Endpoint(),
		RegisterEndpoint: httptransport.NewClient("POST", tgt, encodeRegisterRequest, decodeRegisterResponse, options...).Endpoint(),
	}, nil
}

func encodeRegisterRequest(ctx context.Context, req *http.Request, request interface{}) error {
	req.Method, req.URL.Path = "POST", "/register/"
	return encodeRequest(ctx, req, request)
}

func decodeRegisterResponse(_ context.Context, resp *http.Response) (interface{}, error) {
	var response registerResponse
	err := json.NewDecoder(resp.Body).Decode(&response)
	return response, err
}

func encodeLoginRequest(ctx context.Context, req *http.Request, request interface{}) error {
	req.Method, req.URL.Path = "POST", "/login/"
	return encodeRequest(ctx, req, request)
}

func decodeLoginResponse(_ context.Context, resp *http.Response) (interface{}, error) {
	var response loginResponse
	err := json.NewDecoder(resp.Body).Decode(&response)
	return response, err
}

// encodeRequest likewise JSON-encodes the request to the HTTP request body.
// Don't use it directly as a transport/http.Client EncodeRequestFunc:
// profilesvc endpoints require mutating the HTTP method and request path.
func encodeRequest(_ context.Context, req *http.Request, request interface{}) error {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(request)
	if err != nil {
		return err
	}
	req.Body = ioutil.NopCloser(&buf)
	return nil
}
