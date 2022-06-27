package middleware

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/leebenson/conform"
	"github.com/stkr89/go-auth-service-template/types"
)

func ConformSignInInput() endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (interface{}, error) {
			req := request.(*types.SignInRequest)
			err := conform.Strings(req)
			if err != nil {
				return nil, err
			}
			return next(ctx, req)
		}
	}
}

func ConformSignUpInput() endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (interface{}, error) {
			req := request.(*types.SignUpRequest)
			err := conform.Strings(req)
			if err != nil {
				return nil, err
			}
			return next(ctx, req)
		}
	}
}
