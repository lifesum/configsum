package main

import (
	"flag"
	"net/http"
	"net/http/pprof"
	"os"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Log fields.
const (
	logCaller    = "caller"
	logDuration  = "duration"
	logHostname  = "hostname"
	logJob       = "job"
	logLifecycle = "lifecycle"
	logListen    = "listen"
	logNow       = "now"
	logRevision  = "revision"
	logService   = "service"
	logTask      = "task"
)

// Lifecycles.
const (
	lifecycleAbort = "abort"
	lifecycleStart = "start"
)

// Buildtime vars.
var revision = "0000000-dev"

func main() {
	var (
		begin = time.Now()

		debug               = flag.Bool("debug", false, "enable debug logging")
		instrumentationAddr = flag.String("instrumentation.addir", ":8701", "Listen address for instrumentation")
		listenAddr          = flag.String("listen.addr", ":8700", "Listen address for HTTP API")
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
			logService, "instrumentation",
		)

		abort(logger, http.ListenAndServe(addr, mux))
	}(logger, *instrumentationAddr)

	srv := &http.Server{
		Addr:         *listenAddr,
		Handler:      http.NewServeMux(),
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
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

	_ = logger.Log("err", err, logLifecycle, lifecycleAbort)
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
