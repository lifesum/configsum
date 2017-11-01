package instrument

import (
	"fmt"
	"strconv"
	"time"

	kitprom "github.com/go-kit/kit/metrics/prometheus"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/lifesum/configsum/pkg/errors"
)

// Labels.
const (
	labelErr        = "err"
	labelHost       = "host"
	labelMethod     = "method"
	labelOp         = "op"
	labelProto      = "proto"
	labelRepo       = "repo"
	labelRoute      = "route"
	labelStatusCode = "statusCode"
	labelStore      = "store"
)

var (
	repoLatencies    = map[string]*kitprom.Histogram{}
	requestLatencies = map[string]*kitprom.Histogram{}
)

// ObserveRepoFunc wraps a histogram to track repo op latencies.
type ObserveRepoFunc func(store, repo, op string, begin time.Time, err error)

// ObserveRepo wraps a histogram to track repo op latencies.
func ObserveRepo(namespace, subsystem string) ObserveRepoFunc {
	key := fmt.Sprintf("%s-%s", namespace, subsystem)

	_, ok := repoLatencies[key]
	if !ok {
		repoLatencies[key] = kitprom.NewHistogramFrom(
			prometheus.HistogramOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "op_latency_seconds",
				Help:      "Latency of repo operations.",
			},
			[]string{
				labelErr,
				labelOp,
				labelRepo,
				labelStore,
			},
		)
	}

	return func(store, repo, op string, begin time.Time, err error) {
		errVal := ""

		if e := errors.Cause(err); e != nil {
			errVal = e.Error()
		}

		repoLatencies[key].With(
			labelErr, errVal,
			labelOp, op,
			labelRepo, repo,
			labelStore, store,
		).Observe(time.Since(begin).Seconds())
	}
}

// ObserveRequestFunc wraps a histogram to track request latencies.
type ObserveRequestFunc func(
	code int,
	host, method, proto, route string,
	begin time.Time,
)

// ObserveRequest wraps a histogram to track request latencies.
func ObserveRequest(namespace, subsystem string) ObserveRequestFunc {
	key := fmt.Sprintf("%s-%s", namespace, subsystem)

	_, ok := requestLatencies[key]
	if !ok {
		requestLatencies[key] = kitprom.NewHistogramFrom(
			prometheus.HistogramOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "transport_http_latency_seconds",
				Help:      "Total duration of requests in seconds",
			},
			[]string{
				labelHost,
				labelMethod,
				labelProto,
				labelStatusCode,
			},
		)
	}

	return func(code int, host, method, proto, route string, begin time.Time) {
		requestLatencies[key].With(
			labelStatusCode, strconv.Itoa(code),
			labelHost, host,
			labelMethod, method,
			labelProto, proto,
			labelRoute, route,
		).Observe(time.Since(begin).Seconds())
	}
}
