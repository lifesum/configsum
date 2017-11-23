package config

import (
	"context"
	"encoding/json"
	"net/http"
	"sort"
	"time"

	"github.com/go-kit/kit/endpoint"
	"golang.org/x/text/language"

	"github.com/lifesum/configsum/pkg/auth"
	"github.com/lifesum/configsum/pkg/client"
	"github.com/lifesum/configsum/pkg/errors"
)

type baseCreateRequest struct {
	clientID string
	name     string
}

func baseCreateEndpoint(svc BaseService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(baseCreateRequest)

		c, err := svc.Create(req.clientID, req.name)
		if err != nil {
			return nil, err
		}

		return responseBaseConfig{config: c}, nil
	}
}

type baseGetRequest struct {
	id string
}

func baseGetEndpoint(svc BaseService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(baseGetRequest)

		c, err := svc.Get(req.id)
		if err != nil {
			return nil, err
		}

		return responseBaseConfig{config: c}, nil
	}
}

type baseListRequest struct{}

type baseListResponse struct {
	baseConfigs []BaseConfig
}

func (r baseListResponse) MarshalJSON() ([]byte, error) {
	cs := []responseBaseConfig{}

	for _, c := range r.baseConfigs {
		cs = append(cs, responseBaseConfig{config: c})
	}

	return json.Marshal(struct {
		BaseConfigs []responseBaseConfig `json:"base_configs"`
	}{
		BaseConfigs: cs,
	})
}

type responseParameter struct {
	Name  string      `json:"name"`
	Type  string      `json:"type"`
	Value interface{} `json:"value"`
}

type responseParameters []responseParameter

func (r responseParameters) Len() int {
	return len(r)
}

func (r responseParameters) Less(i, j int) bool {
	return r[i].Name > r[j].Name
}

func (r responseParameters) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

type responseBaseConfig struct {
	config BaseConfig
}

func (r responseBaseConfig) MarshalJSON() ([]byte, error) {
	v := struct {
		ClientID   string             `json:"client_id"`
		Deleted    bool               `json:"deleted"`
		ID         string             `json:"id"`
		Name       string             `json:"name"`
		Parameters responseParameters `json:"parameters"`
		CreatedAt  time.Time          `json:"created_at"`
		UpdatedAt  time.Time          `json:"updated_at"`
	}{
		ClientID:  r.config.ClientID,
		Deleted:   r.config.Deleted,
		ID:        r.config.ID,
		Name:      r.config.Name,
		CreatedAt: r.config.CreatedAt,
		UpdatedAt: r.config.UpdatedAt,
	}

	ps := responseParameters{}

	for k, val := range r.config.Parameters {
		p := responseParameter{
			Name:  k,
			Value: val,
		}

		switch val.(type) {
		case bool:
			p.Type = "bool"
		case float64:
			p.Type = "number"
		case string:
			p.Type = "string"
		default:
			p.Type = "unknown"
		}

		ps = append(ps, p)
	}

	sort.Sort(ps)

	v.Parameters = ps

	return json.Marshal(v)
}

func baseListEndpoint(svc BaseService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		cs, err := svc.List()
		if err != nil {
			return nil, err
		}

		return baseListResponse{baseConfigs: cs}, nil
	}
}

type baseUpdateRequest struct {
	id         string
	parameters rendered
}

func baseUpdateEndpoint(svc BaseService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(baseUpdateRequest)

		c, err := svc.Update(req.id, req.parameters)
		if err != nil {
			return nil, err
		}

		return responseBaseConfig{config: c}, nil
	}
}

type device struct {
	Location location `json:"location"`
}

type location struct {
	locale language.Tag
}

func (l *location) UnmarshalJSON(raw []byte) error {
	v := struct {
		Locale string `json:"locale"`
	}{}

	if err := json.Unmarshal(raw, &v); err != nil {
		return err
	}

	t, err := language.Parse(v.Locale)
	if err != nil {
		return errors.Wrapf(errors.ErrInvalidPayload, "invalid locale %s", v.Locale)
	}

	l.locale = t

	return nil
}

type userRenderContext struct {
	Device device `json:"device"`
}

type userRenderRequest struct {
	baseConfig string
	context    userRenderContext
}

type userRenderResponse struct {
	baseID    string
	baseName  string
	clientID  string
	id        string
	rendered  rendered
	createdAt time.Time
}

func (r userRenderResponse) StatusCode() int {
	return http.StatusCreated
}

func userRenderEndpoint(svc UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		var (
			req      = request.(userRenderRequest)
			clientID = ctx.Value(client.ContextKeyClientID).(string)
			userID   = ctx.Value(auth.ContextKeyUserID).(string)
		)

		c, err := svc.Render(clientID, req.baseConfig, userID)
		if err != nil {
			return nil, err
		}

		return userRenderResponse{
			baseID:    c.baseID,
			baseName:  req.baseConfig,
			clientID:  clientID,
			id:        c.id,
			rendered:  c.rendered,
			createdAt: c.createdAt,
		}, nil
	}
}
