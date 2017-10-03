package main

import (
	"flag"
	"fmt"
	"net/http"
	"net/http/pprof"
	"os"
	"os/user"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	kitprom "github.com/go-kit/kit/metrics/prometheus"
	"github.com/jmoiron/sqlx"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/lifesum/configsum/pkg/config"
	"github.com/lifesum/configsum/pkg/pg"
)

// Versions.
const (
	apiVersion = "v1"
)

// Log fields.
const (
	logCaller    = "caller"
	logDuration  = "duration"
	logError     = "err"
	logHostname  = "hostname"
	logJob       = "job"
	logLifecycle = "lifecycle"
	logListen    = "listen"
	logNow       = "now"
	logRevision  = "revision"
	logService   = "service"
	logTask      = "task"
)

// Instrument labels.
const (
	labelOp    = "op"
	labelRepo  = "repo"
	labelStore = "store"
)

// Instrument fields.
const (
	instrumentNamespace = "configsum"
	instrumentSubsystem = "config_api"
)

// Lifecycles.
const (
	lifecycleAbort = "abort"
	lifecycleStart = "start"
)

// Timeouts.
const (
	defaultTimeoutRead  = 1 * time.Second
	defaultTimeoutWrite = 1 * time.Second
)

const storeRepo = "postgres"

// Buildtime vars.
var revision = "0000000-dev"

// Default vars.
var defaultPostgresURI string

func main() {
	var (
		begin = time.Now()

		debug         = flag.Bool("debug", false, "enable debug logging")
		intrumentAddr = flag.String("instrument.addir", ":8701", "Listen address for instrumentation")
		listenAddr    = flag.String("listen.addr", ":8700", "Listen address for HTTP API")
		postgresURI   = flag.String("postgres.uri", defaultPostgresURI, "URI for Posgres connection")
	)
	flag.Parse()

	// Setup logging.
	var logger log.Logger
	{
		logLevel := level.AllowInfo()
		if *debug {
			logLevel = level.AllowAll()
		}
		logger = log.With(
			log.NewJSONLogger(os.Stdout),
			logCaller, log.Caller(5),
			logJob, "configsum",
			logNow, log.DefaultTimestampUTC,
			logRevision, revision,
			logTask, "config",
		)
		logger = level.NewFilter(logger, logLevel)
	}

	hostname, err := os.Hostname()
	if err != nil {
		abort(logger, err)
	}

	logger = log.With(logger, logHostname, hostname)

	// Setup instrunentation.
	go func(logger log.Logger, addr string) {
		mux := http.NewServeMux()

		registerMetrics(mux)
		registerProfile(mux)

		logger.Log(
			logDuration, time.Since(begin).Nanoseconds(),
			logLifecycle, lifecycleStart,
			logListen, addr,
			logService, "instrument",
		)

		abort(logger, http.ListenAndServe(addr, mux))
	}(logger, *intrumentAddr)

	repoLabels := []string{
		labelOp,
		labelRepo,
		labelStore,
	}

	repoErrCount := kitprom.NewCounterFrom(
		prometheus.CounterOpts{
			Namespace: instrumentNamespace,
			Subsystem: instrumentSubsystem,
			Name:      "err_count",
			Help:      "Amount of failed repo operations.",
		},
		repoLabels,
	)

	repoErrCountFunc := func(store, repo, op string) {
		repoErrCount.With(
			labelOp, op,
			labelRepo, repo,
			labelStore, store,
		).Add(1)
	}

	repoOpCount := kitprom.NewCounterFrom(
		prometheus.CounterOpts{
			Namespace: instrumentNamespace,
			Subsystem: instrumentSubsystem,
			Name:      "op_count",
			Help:      "Amount of successful repo operations.",
		},
		repoLabels,
	)

	repoOpCountFunc := func(store, repo, op string) {
		repoOpCount.With(
			labelOp, op,
			labelRepo, repo,
			labelStore, store,
		).Add(1)
	}

	repoOpLatency := kitprom.NewHistogramFrom(
		prometheus.HistogramOpts{
			Namespace: instrumentNamespace,
			Subsystem: instrumentSubsystem,
			Name:      "op_latency_seconds",
			Help:      "Latency of successful repo operations.",
		},
		repoLabels,
	)

	repoOpLatencyFunc := func(store, repo, op string, begin time.Time) {
		repoOpLatency.With(
			labelOp, op,
			labelRepo, repo,
			labelStore, store,
		).Observe(time.Since(begin).Seconds())
	}

	// Setup clients.
	db, err := sqlx.Connect(storeRepo, *postgresURI)
	if err != nil {
		abort(logger, err)
	}

	// Setup repos.
	baseRepo, err := config.NewInmemBaseRepo()
	if err != nil {
		abort(logger, err)
	}

	userRepo, err := config.NewPostgresUserRepo(db)
	if err != nil {
		abort(logger, err)
	}
	userRepo = config.NewUserRepoInstrumentMiddleware(
		repoErrCountFunc,
		repoOpCountFunc,
		repoOpLatencyFunc,
		storeRepo,
	)(userRepo)
	userRepo = config.NewUserRepoLogMiddleware(logger, storeRepo)(userRepo)

	// Setup service.
	var (
		mux          = http.NewServeMux()
		prefixConfig = fmt.Sprintf(`/%s/config`, apiVersion)
		svc          = config.NewServiceUser(baseRepo, userRepo)
	)

	mux.Handle(
		fmt.Sprintf(`%s/`, prefixConfig),
		http.StripPrefix(
			prefixConfig,
			config.MakeHandler(logger, svc),
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
		logService, "api",
	)

	abort(logger, srv.ListenAndServe())
}

func abort(logger log.Logger, err error) {
	if err != nil {
		return
	}

	_ = logger.Log(logError, err, logLifecycle, lifecycleAbort)
	os.Exit(1)
}

func registerMetrics(mux *http.ServeMux) {
	mux.Handle("/metrics", promhttp.Handler())
}

func registerProfile(mux *http.ServeMux) {
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	mux.Handle("/debug/pprof/block", pprof.Handler("block"))
	mux.Handle("/debug/pprof/goroutine", pprof.Handler("goroutine"))
	mux.Handle("/debug/pprof/heap", pprof.Handler("heap"))
	mux.Handle("/debug/pprof/mutex", pprof.Handler("mutex"))
	mux.Handle("/debug/pprof/threadcreate", pprof.Handler("threadcreate"))
}

func init() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}

	defaultPostgresURI = fmt.Sprintf(pg.DefaultDevURI, user.Username)
}
