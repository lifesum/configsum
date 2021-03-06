package main

import (
	"flag"
	"fmt"
	"math/rand"
	"net/http"
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
	"github.com/lifesum/configsum/pkg/generate"
	"github.com/lifesum/configsum/pkg/instrument"
	"github.com/lifesum/configsum/pkg/rule"
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

		_ = logger.Log(
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

	var baseRepo config.BaseRepo
	baseRepo = config.NewPostgresBaseRepo(db)
	baseRepo = config.NewBaseRepoInstrumentMiddleware(
		instrument.ObserveRepo(instrumentNamespace, taskConsole),
		storeRepo,
	)(baseRepo)
	baseRepo = config.NewBaseRepoLogMiddleware(logger, storeRepo)(baseRepo)

	var userRepo config.UserRepo
	userRepo = config.NewPostgresUserRepo(db)
	userRepo = config.NewUserRepoInstrumentMiddleware(
		instrument.ObserveRepo(instrumentNamespace, taskConfig),
		storeRepo,
	)(userRepo)
	userRepo = config.NewUserRepoLogMiddleware(logger, storeRepo)(userRepo)

	var clientRepo client.Repo
	clientRepo = client.NewPostgresRepo(db)
	clientRepo = client.NewRepoInstrumentMiddleware(
		instrument.ObserveRepo(instrumentNamespace, taskConfig),
		storeRepo,
	)(clientRepo)
	clientRepo = client.NewRepoLogMiddleware(logger, storeRepo)(clientRepo)

	var tokenRepo client.TokenRepo
	tokenRepo = client.NewPostgresTokenRepo(db)
	tokenRepo = client.NewTokenRepoInstrumentMiddleware(
		instrument.ObserveRepo(instrumentNamespace, taskConfig),
		storeRepo,
	)(tokenRepo)
	tokenRepo = client.NewTokenRepoLogMiddleware(logger, storeRepo)(tokenRepo)

	ruleRepo := rule.NewPostgresRepo(db)
	ruleRepo = rule.NewRuleRepoInstrumentMiddleware(
		instrument.ObserveRepo(instrumentNamespace, taskConfig),
		storeRepo,
	)(ruleRepo)
	ruleRepo = rule.NewRuleRepoLogMiddleware(logger, storeRepo)(ruleRepo)

	// Setup service.
	var (
		seed         = rand.New(rand.NewSource(time.Now().UnixNano()))
		mux          = http.NewServeMux()
		prefixConfig = fmt.Sprintf(`/%s/config`, apiVersion)
		clientSVC    = client.NewService(clientRepo, tokenRepo)
		svc          = config.NewUserService(baseRepo, userRepo, ruleRepo, generate.RandPercentage(seed))
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
