package main

import (
	"flag"
	"fmt"
	"net/http"
	"net/http/pprof"
	"os"
	"os/user"
	"runtime"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/prometheus/client_golang/prometheus/promhttp"

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

// Tasks.
const (
	taskConfig = "config"
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

type runFunc func([]string, log.Logger) error

func main() {
	var (
		flagset = flag.NewFlagSet("configsum", flag.ExitOnError)

		debug = flagset.Bool("debug", false, "enable debug logging")
	)

	flagset.Usage = usage
	if err := flagset.Parse(os.Args[1:]); err != nil {
		panic(err)
	}

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
		)
		logger = level.NewFilter(logger, logLevel)
	}

	hostname, err := os.Hostname()
	if err != nil {
		abort(logger, err)
	}

	logger = log.With(logger, logHostname, hostname)

	if len(flagset.Args()) < 1 {
		usage()
		os.Exit(1)
	}

	var (
		task = strings.ToLower(flagset.Arg(0))

		run runFunc
	)

	switch task {
	case taskConfig:
		run = runConfig
	default:
		usage()
		os.Exit(1)
	}

	logger = log.With(logger, logTask, task)

	abort(logger, run(flagset.Args()[1:], logger))
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

func usage() {
	f := `USAGE
	%s [FLAGS] <cmd> [FLAGS]

COMMANDS
	config	API offering access to per user rendered configs

VERSION
	%s (%s)
`

	fmt.Fprintf(os.Stderr, f, os.Args[0], revision, runtime.Version())
}

func usageCmd(fs *flag.FlagSet, short string) func() {
	s := `USAGE
  configsum %s

FLAGS
`
	return func() {
		fmt.Fprintf(os.Stderr, s, short)

		w := tabwriter.NewWriter(os.Stderr, 0, 2, 2, ' ', 0)
		fs.VisitAll(func(f *flag.Flag) {
			fmt.Fprintf(w, "\t-%s\t%s\t%s\n", f.Name, f.DefValue, f.Usage)
		})
		_ = w.Flush()
	}
}

func init() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}

	defaultPostgresURI = fmt.Sprintf(pg.DefaultDevURI, user.Username)
}
