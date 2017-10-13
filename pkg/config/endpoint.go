package config

import (
	"context"
	"net/http"
	"time"

	"github.com/go-kit/kit/endpoint"

	"github.com/lifesum/configsum/pkg/client"
)

type userRequest struct {
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
		var (
			req      = request.(userRequest)
			clientID = ctx.Value(client.ContextKeyClientID).(string)
		)

		c, err := svc.Render(clientID, req.baseConfig, "id123")
		if err != nil {
			return nil, err
		}

		return userResponse{
			BaseConfig: c.baseID,
			CreatedAt:  c.createdAt,
		}, nil
	}
}
