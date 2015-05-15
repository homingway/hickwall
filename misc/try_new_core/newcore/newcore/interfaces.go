package newcore

import (
	"time"
)

type CollectResult struct {
	Collected *MuliDataPoint // slice of []*DataPoint
	Next      time.Time
	Err       error
}

// A Collector collects Items and returns the time when the next collect should be
// attempted.  On failure, CollectOnce returns a non-nil error.
// recover from panic and timeout should be handled properly with CollectOnce
type Collector interface {
	// CollectOnce() (md MuliDataPoint, next time.Time, err error)
	CollectOnce() *CollectResult
}

// A Subscription delivers Items over a channel.  Close cancels the
// subscription, closes the Updates channel, and returns the last collect error,
// if any.
type Subscription interface {
	Updates() <-chan *DataPoint
	Close() error
}
