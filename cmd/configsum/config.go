package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"github.com/lifesum/configsum/pkg/auth/dory"
	"github.com/lifesum/configsum/pkg/auth/simple"
	"github.com/lifesum/configsum/pkg/client"
	"github.com/lifesum/configsum/pkg/config"
	"github.com/lifesum/configsum/pkg/instrument"
	confhttp "github.com/lifesum/configsum/pkg/transport/http"
)

const (
	authDory   = "dory"
	authSimple = "simple"
)

func runConfig(args []string, logger log.Logger) error {
	var (
		begin   = time.Now()
		flagset = flag.NewFlagSet("config", flag.ExitOnError)

		authMethod    = flagset.String("auth", authSimple, "User authenticaiton method to use (dory, simple)")
		dorySecret    = flagset.String("dory.secret", "", "Shared secret for Dory Authentication middleware")
		baseState     = flagset.String("base.state", "", "Initial base_config repo state")
		intrumentAddr = flagset.String("instrument.addir", ":8701", "Listen address for instrumentation")
		listenAddr    = flagset.String("listen.addr", ":8700", "Listen address for HTTP API")
		postgresURI   = flagset.String("postgres.uri", defaultPostgresURI, "URI for Posgres connection")
	)

	flagset.Usage = usageCmd(flagset, "config [flags]")
	if err := flagset.Parse(args); err != nil {
		return err
	}

	// Setup instrunentation.
	go func(logger log.Logger, addr string) {
		mux := http.NewServeMux()

		registerMetrics(mux)
		registerProfile(mux)

		logger.Log(
			logDuration, time.Since(begin).Nanoseconds(),
			logLifecycle, lifecycleStart,
			logListen, addr,
			logService, serviceInstrument,
		)

		abort(logger, http.ListenAndServe(addr, mux))
	}(logger, *intrumentAddr)

	logger = log.With(logger, logService, serviceAPI)

	// Setup clients.
	db, err := sqlx.Connect(storeRepo, *postgresURI)
	if err != nil {
		return err
	}

	// Setup repos.
	var state config.InmemBaseState

	// TODO(xla): Temporary solution to set up base configs without proper repo
	// integration.
	if *baseState != "" {
		f, err := os.Open(*baseState)
		if err != nil {
			return err
		}

		t := struct {
			ClientID   string                 `json:"clientID"`
			ID         string                 `json:"id"`
			Name       string                 `json:"name"`
			Parameters map[string]interface{} `json:"rendered"`
		}{}

		err = json.NewDecoder(f).Decode(&t)
		if err != nil {
			return err
		}

		state = config.InmemBaseState{
			t.ClientID: map[string]config.BaseConfig{
				t.Name: config.BaseConfig{
					ClientID:   t.ClientID,
					ID:         t.ID,
					Name:       t.Name,
					Parameters: t.Parameters,
				},
			},
		}
	}

	baseRepo := config.NewInmemBaseRepo(state)

	userRepo := config.NewPostgresUserRepo(db)
	userRepo = config.NewUserRepoInstrumentMiddleware(
		instrument.ObserveRepo(instrumentNamespace, taskConfig),
		storeRepo,
	)(userRepo)
	userRepo = config.NewUserRepoLogMiddleware(logger, storeRepo)(userRepo)

	clientRepo := client.NewPostgresRepo(db)
	clientRepo = client.NewRepoInstrumentMiddleware(
		instrument.ObserveRepo(instrumentNamespace, taskConfig),
		storeRepo,
	)(clientRepo)
	clientRepo = client.NewRepoLogMiddleware(logger, storeRepo)(clientRepo)

	tokenRepo := client.NewPostgresTokenRepo(db)
	tokenRepo = client.NewTokenRepoInstrumentMiddleware(
		instrument.ObserveRepo(instrumentNamespace, taskConfig),
		storeRepo,
	)(tokenRepo)
	tokenRepo = client.NewTokenRepoLogMiddleware(logger, storeRepo)(tokenRepo)

	// Setup service.
	var (
		mux          = http.NewServeMux()
		prefixConfig = fmt.Sprintf(`/%s/config`, apiVersion)
		clientSVC    = client.NewService(clientRepo, tokenRepo)
		svc          = config.NewUserService(baseRepo, userRepo)
		opts         = []kithttp.ServerOption{
			kithttp.ServerBefore(kithttp.PopulateRequestContext),
			kithttp.ServerBefore(confhttp.PopulateRequestContext),
			kithttp.ServerErrorEncoder(confhttp.ErrorEncoder),
			kithttp.ServerFinalizer(
				confhttp.ServerFinalizer(
					logger,
					instrument.ObserveRequest(instrumentNamespace, taskConfig),
				),
			),
		}

		auth endpoint.Middleware
	)

	auth = endpoint.Chain(client.AuthMiddleware(clientSVC))
	opts = append(opts, kithttp.ServerBefore(client.HTTPToContext))

	switch *authMethod {
	case authDory:
		auth = endpoint.Chain(auth, dory.AuthMiddleware(*dorySecret))
		opts = append(opts, kithttp.ServerBefore(dory.HTTPToContext))
	case authSimple:
		auth = endpoint.Chain(auth, simple.AuthMiddleware())
		opts = append(opts, kithttp.ServerBefore(simple.HTTPToContext))
	default:
		return errors.Errorf("unsupported auth: '%s'", *authMethod)
	}

	mux.Handle(
		fmt.Sprintf(`%s/`, prefixConfig),
		http.StripPrefix(
			prefixConfig,
			config.MakeHandler(
				svc,
				auth,
				opts...,
			),
		),
	)

	// Setup server.
	srv := &http.Server{
		Addr:         *listenAddr,
		Handler:      mux,
		ReadTimeout:  defaultTimeoutRead,
		WriteTimeout: defaultTimeoutWrite,
	}

	_ = level.Info(logger).Log(
		logDuration, time.Since(begin).Nanoseconds(),
		logLifecycle, lifecycleStart,
		logListen, *listenAddr,
	)

	return srv.ListenAndServe()
}
