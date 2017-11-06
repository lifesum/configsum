package client

import (
	"context"
	"encoding/json"
	"net/http"
	"sort"
	"time"

	"github.com/go-kit/kit/endpoint"
)

type createRequest struct {
	name string
}

type createResponse struct {
	client Client
	token  string
}

func (r createResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal(responseClient{
		CreatedAt: r.client.createdAt,
		Deleted:   r.client.deleted,
		ID:        r.client.id,
		Name:      r.client.name,
		Token:     r.token,
	})
}

func (r createResponse) StatusCode() int {
	return http.StatusCreated
}

type listResponse struct {
	clientTokens clientTokens
}

func (r listResponse) MarshalJSON() ([]byte, error) {
	cs := responseClientList{}

	for c, t := range r.clientTokens {
		cs = append(cs, responseClient{
			CreatedAt: c.createdAt,
			Deleted:   c.deleted,
			ID:        c.id,
			Name:      c.name,
			Token:     t.secret,
		})
	}

	sort.Sort(cs)

	return json.Marshal(struct {
		Clients responseClientList `json:"clients"`
	}{
		Clients: cs,
	})
}

type responseClient struct {
	CreatedAt time.Time `json:"created_at"`
	Deleted   bool      `json:"deleted"`
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Token     string    `json:"token"`
}

type responseClientList []responseClient

func (l responseClientList) Len() int {
	return len(l)
}

func (l responseClientList) Less(i, j int) bool {
	return l[i].CreatedAt.After(l[j].CreatedAt)
}

func (l responseClientList) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

func createEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		r := request.(createRequest)

		c, secret, err := svc.Create(r.name)
		if err != nil {
			return nil, err
		}

		return createResponse{
			client: c,
			token:  secret,
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
