package config

import (
	"context"
	"net/http"
	"time"

	"github.com/go-kit/kit/endpoint"
)

type userRequest struct{}

type userResponse struct {
	CreatedAt time.Time `json:"createdAt"`
}

func (r userResponse) StatusCode() int {
	return http.StatusCreated
}

func userEndpoint(svc ServiceUser) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		c, err := svc.Get()
		if err != nil {
			return nil, err
		}

		return userResponse{CreatedAt: c.createdAt}, nil
	}
}
