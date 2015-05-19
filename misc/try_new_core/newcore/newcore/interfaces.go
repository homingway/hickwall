package newcore

import (
	"time"
)

type CollectResult struct {
	Collected *MultiDataPoint // slice of []*DataPoint
	Next      time.Time       // when the next run will take place
	Err       error
}

// A Collector collects Items and returns the time when the next collect should be
// attempted.  On failure, CollectOnce returns a non-nil error.
// recover from panic and timeout should be handled properly with CollectOnce.
type Collector interface {
	CollectOnce() *CollectResult
	Interval() time.Duration
	IsEnabled() bool
	Name() string
	ClassName() string
	Close() error
}

// A Subscription delivers Items over a channel.  Close cancels the
// subscription, closes the Updates channel, and returns the last collect error,
// if any.
type Subscription interface {
	Updates() <-chan *MultiDataPoint
	Name() string
	Close() error
}

type Publication interface {
	Updates() chan<- *MultiDataPoint
	Name() string
	Close() error
}

type PublicationSet interface {
	Close() error // close subscription, backends, and fanout
}
