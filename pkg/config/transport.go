package config

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"

	"github.com/lifesum/configsum/pkg/client"
)

// Headers.
const (
	headerContentType = "Content-Type"
)

// MakeHandler returns an http.Handler for the config service.
func MakeHandler(
	logger log.Logger,
	svc ServiceUser,
	clientSVC client.Service,
) http.Handler {
	r := mux.NewRouter()
	r.StrictSlash(true)

	r.Methods("PUT").Path(`/{baseConfig:[a-z0-9]+}`).Name("configUser").Handler(
		kithttp.NewServer(
			client.AuthMiddleware(clientSVC)(
				userEndpoint(svc),
			),
			decodeUserRequest,
			kithttp.EncodeJSONResponse,
			kithttp.ServerBefore(client.HTTPToContext),
			kithttp.ServerErrorEncoder(encodeError),
			kithttp.ServerErrorLogger(log.With(logger, "route", "configUser")),
		),
	)

	return r
}

func decodeUserRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	c, ok := mux.Vars(r)["baseConfig"]
	if !ok {
		return nil, fmt.Errorf("Baseconfig missing")
	}

	return userRequest{baseConfig: c}, nil
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	switch errors.Cause(err) {
	case ErrNotFound:
		w.WriteHeader(http.StatusNotFound)
	case client.ErrNotFound, client.ErrSecretMissing:
		w.WriteHeader(http.StatusUnauthorized)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.Header().Set(headerContentType, "application/json; charset=utf-8")

	_ = json.NewEncoder(w).Encode(struct {
		Reason string `json:"reason"`
	}{
		Reason: err.Error(),
	})
}
