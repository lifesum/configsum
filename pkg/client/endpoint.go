package client

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-kit/kit/endpoint"
)

type createRequest struct {
	name string
}

type createResponse struct {
	CreatedAt time.Time `json:"created_at"`
	Deleted   bool      `json:"deleted"`
	ID        string    `json:"id"`
	Name      string    `json:"name"`
}

func (r createResponse) StatusCode() int {
	return http.StatusCreated
}

type listResponse struct {
	clientTokens clientTokens
}

func (r listResponse) MarshalJSON() ([]byte, error) {
	type client struct {
		CreatedAt time.Time `json:"created_at"`
		Deleted   bool      `json:"deleted"`
		ID        string    `json:"id"`
		Name      string    `json:"name"`
		Token     string    `json:"token"`
	}

	cs := []client{}

	for c, t := range r.clientTokens {
		cs = append(cs, client{
			CreatedAt: c.createdAt,
			Deleted:   c.deleted,
			ID:        c.id,
			Name:      c.name,
			Token:     t.secret,
		})
	}

	return json.Marshal(struct {
		Clients []client `json:"clients"`
	}{
		Clients: cs,
	})
}

func createEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		r := request.(createRequest)

		c, err := svc.Create(r.name)
		if err != nil {
			return nil, err
		}

		return createResponse{
			CreatedAt: c.createdAt,
			Deleted:   c.deleted,
			ID:        c.id,
			Name:      c.name,
		}, nil
	}
}

func listEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		ct, err := svc.ListWithToken()
		if err != nil {
			return nil, err
		}

		return listResponse{
			clientTokens: ct,
		}, nil
	}
}
