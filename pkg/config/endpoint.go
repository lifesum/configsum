package config

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

type userRequest struct{}

func userEndpoint(svc ServiceUser) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		err := svc.Get()

		return nil, err
	}
}
