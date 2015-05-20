package newcore

import (
	"time"
)

type Interval string

func NewInterval(v string) *Interval {
	intv := Interval(v)
	return &intv
}

func (i *Interval) MustDuration(dft time.Duration) time.Duration {
	res := dft
	if res < 0 {
		res = time.Duration(0)
	}

	if *i == "" {
		return res
	}

	d, err := time.ParseDuration(string(*i))
	// interval must be positive
	if d < 0 || err != nil {
		return res
	}
	return d
}

func (i *Interval) Duration() (time.Duration, error) {
	return time.ParseDuration(string(*i))
}
