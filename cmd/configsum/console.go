package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/jmoiron/sqlx"

	"github.com/lifesum/configsum/pkg/client"
	"github.com/lifesum/configsum/pkg/config"
	"github.com/lifesum/configsum/pkg/instrument"
	confhttp "github.com/lifesum/configsum/pkg/transport/http"
	"github.com/lifesum/configsum/pkg/ui"
)

func runConsole(args []string, logger log.Logger) error {
	var (
		begin   = time.Now()
		flagset = flag.NewFlagSet("console", flag.ExitOnError)

		instrumentAddr = flagset.String("instrument.addr", ":8711", "Listen address for instrumenation")
		listenAddr     = flagset.String("listen.addr", ":8710", "HTTP API bind address")
		postgresURI    = flagset.String("postgres.uri", defaultPostgresURI, "URI for Posgres connection")
		uiBase         = flagset.String("ui.base", "/", "Base URI to use for path based mounting")
		uiLocal        = flagset.Bool("ui.local", false, "Load static assets from the filesystem")
	)

	flagset.Usage = usageCmd(flagset, "console [flags]")
	if err := flagset.Parse(args); err != nil {
		return err
	}

	go func(lgoger log.Logger, addr string) {
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
	}(logger, *instrumentAddr)

	db, err := sqlx.Connect(storeRepo, *postgresURI)
	if err != nil {
		return err
	}

	baseRepo := config.NewPostgresBaseRepo(db)
	baseRepo = config.NewBaseRepoInstrumentMiddleware(
		instrument.ObserveRepo(instrumentNamespace, taskConsole),
		storeRepo,
	)(baseRepo)
	baseRepo = config.NewBaseRepoLogMiddleware(logger, storeRepo)(baseRepo)

	clientRepo := client.NewPostgresRepo(db)
	clientRepo = client.NewRepoInstrumentMiddleware(
		instrument.ObserveRepo(instrumentNamespace, taskConsole),
		storeRepo,
	)(clientRepo)
	clientRepo = client.NewRepoLogMiddleware(logger, storeRepo)(clientRepo)

	tokenRepo := client.NewPostgresTokenRepo(db)
	tokenRepo = client.NewTokenRepoInstrumentMiddleware(
		instrument.ObserveRepo(instrumentNamespace, taskConsole),
		storeRepo,
	)(tokenRepo)
	tokenRepo = client.NewTokenRepoLogMiddleware(logger, storeRepo)(tokenRepo)

	var (
		baseConfigSVC    = config.NewBaseService(baseRepo, clientRepo)
		clientSVC        = client.NewService(clientRepo, tokenRepo)
		prefixBaseConfig = "/api/configs/base"
		prefixClient     = "/api/clients"
		serveMux         = http.NewServeMux()
		opts             = []kithttp.ServerOption{
			kithttp.ServerBefore(kithttp.PopulateRequestContext),
			kithttp.ServerBefore(confhttp.PopulateRequestContext),
			kithttp.ServerErrorEncoder(confhttp.ErrorEncoder),
			kithttp.ServerFinalizer(
				confhttp.ServerFinalizer(
					logger,
					instrument.ObserveRequest(instrumentNamespace, taskConsole),
				),
			),
		}
	)

	serveMux.Handle(
		fmt.Sprintf("%s/", prefixBaseConfig),
		http.StripPrefix(
			prefixBaseConfig,
			config.MakeBaseHandler(baseConfigSVC, opts...),
		),
	)
	serveMux.Handle(
		fmt.Sprintf("%s/", prefixClient),
		http.StripPrefix(
			prefixClient,
			client.MakeHandler(clientSVC, opts...),
		),
	)
	serveMux.Handle("/", ui.MakeHandler(logger, *uiBase, *uiLocal))

	srv := &http.Server{
		Addr:         *listenAddr,
		Handler:      serveMux,
		ReadTimeout:  defaultTimeoutRead,
		WriteTimeout: defaultTimeoutWrite,
	}

	_ = level.Info(logger).Log(
		logDuration, time.Since(begin).Nanoseconds(),
		logLifecycle, lifecycleStart,
		logListen, *listenAddr,
		logService, serviceAPI,
	)

	return srv.ListenAndServe()
}
