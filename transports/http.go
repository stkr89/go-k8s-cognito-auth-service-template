package transport

import (
	"context"
	"encoding/json"
	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/stkr89/go-auth-service-template/common"
	"github.com/stkr89/go-auth-service-template/endpoints"
	"github.com/stkr89/go-auth-service-template/middleware"
	"github.com/stkr89/go-auth-service-template/types"
	"net/http"
)

type errorWrapper struct {
	Error string `json:"error"`
}

func NewHTTPHandler(endpoints endpoints.Endpoints) http.Handler {
	m := mux.NewRouter()
	m.Handle("/api/auth/v1/signup", httptransport.NewServer(
		endpoint.Chain(
			middleware.ValidateSignUpInput(),
			middleware.ConformSignUpInput(),
		)(endpoints.SignUp),
		decodeHTTPSignUpRequest,
		encodeHTTPGenericResponse,
	)).Methods(http.MethodPost)
	m.Handle("/api/auth/v1/signin", httptransport.NewServer(
		endpoint.Chain(
			middleware.ValidateSignInInput(),
			middleware.ConformSignInInput(),
		)(endpoints.SignIn),
		decodeHTTPSignInRequest,
		encodeHTTPGenericResponse,
	)).Methods(http.MethodPost)

	return m
}

func err2code(err *common.Error) int {
	switch err.Key {
	case common.InvalidRequestBody:
		return http.StatusBadRequest
	case common.Unauthorized:
		return http.StatusUnauthorized
	}
	return http.StatusInternalServerError
}

func errorEncoder(_ context.Context, err error, w http.ResponseWriter) {
	w.WriteHeader(err2code(err.(*common.Error)))
	json.NewEncoder(w).Encode(errorWrapper{Error: err.Error()})
}

func encodeHTTPGenericResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if f, ok := response.(endpoint.Failer); ok && f.Failed() != nil {
		errorEncoder(ctx, f.Failed(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func decodeHTTPSignInRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req types.SignInRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	return &req, err
}

func decodeHTTPSignUpRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req types.SignUpRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	return &req, err
}
