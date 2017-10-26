package instrument

import "time"

// CountRepoFunc wraps a counter to track vital repo information.
type CountRepoFunc func(store, repo, op string)

// ObserveRepoFunc wraps a histogram to track repo op latencies.
type ObserveRepoFunc func(store, repo, op string, begin time.Time)

// CountRequestFunc wraps a counter to track number of received requests.
type CountRequestFunc func(statusCode int, host, method, proto string)

// ObserveRequestFunc wraps a histogram to track request latencies.
type ObserveRequestFunc func(statusCode int, host, method, proto string, begin time.Time)
