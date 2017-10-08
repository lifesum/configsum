package main

import (
	"flag"
	"net/http"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"

	"github.com/lifesum/configsum/pkg/ui"
)

func runConsole(args []string, logger log.Logger) error {
	var (
		begin   = time.Now()
		flagset = flag.NewFlagSet("console", flag.ExitOnError)

		listenAddr  = flagset.String("listen.addr", ":8700", "HTTP API bind address")
		staticLocal = flagset.Bool("static.local", false, "Determines if static assets are loaded from the filesystem.")
	)

	if err := flagset.Parse(args); err != nil {
		return err
	}

	serveMux := http.NewServeMux()

	handler, err := ui.MakeHandler(logger, *staticLocal)
	if err != nil {
		abort(logger, err)
	}

	serveMux.Handle("/", handler)

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
