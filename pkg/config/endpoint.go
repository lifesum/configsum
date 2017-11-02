package config

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-kit/kit/endpoint"
	"golang.org/x/text/language"

	"github.com/lifesum/configsum/pkg/auth"
	"github.com/lifesum/configsum/pkg/client"
	"github.com/lifesum/configsum/pkg/errors"
)

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

type userContext struct {
	Device device `json:"device"`
}

type userRequest struct {
	baseConfig string
	context    userContext
}

type userResponse struct {
	baseID    string
	baseName  string
	clientID  string
	id        string
	rendered  rendered
	createdAt time.Time
}

func (r userResponse) StatusCode() int {
	return http.StatusCreated
}

func userEndpoint(svc ServiceUser) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		var (
			req      = request.(userRequest)
			clientID = ctx.Value(client.ContextKeyClientID).(string)
			userID   = ctx.Value(auth.ContextKeyUserID).(string)
		)

		c, err := svc.Render(clientID, req.baseConfig, userID)
		if err != nil {
			return nil, err
		}

		return userResponse{
			baseID:    c.baseID,
			baseName:  req.baseConfig,
			clientID:  clientID,
			id:        c.id,
			rendered:  c.rendered,
			createdAt: c.createdAt,
		}, nil
	}
}
