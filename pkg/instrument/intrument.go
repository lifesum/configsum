package instrument

import "time"

// CountRepoFunc wraps a counter to track vital repo information.
type CountRepoFunc func(store, repo, op string)

// ObserveRepoFunc wraps a histogram to track repo op latencies.
type ObserveRepoFunc func(store, repo, op string, begin time.Time)
