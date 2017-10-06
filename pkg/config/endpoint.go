package config

import (
	"context"
	"net/http"
	"time"

	"github.com/go-kit/kit/endpoint"
)

type userRequest struct {
	appID      string
	baseConfig string
}

type userResponse struct {
	BaseConfig string    `json:"baseConfig"`
	CreatedAt  time.Time `json:"createdAt"`
}

func (r userResponse) StatusCode() int {
	return http.StatusCreated
}

func userEndpoint(svc ServiceUser) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(userRequest)

		c, err := svc.Get(req.baseConfig, "id123")
		if err != nil {
			return nil, err
		}

		return userResponse{
			BaseConfig: c.baseID,
			CreatedAt:  c.createdAt,
		}, nil
	}
}
