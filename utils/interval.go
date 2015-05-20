package utils

import (
	"time"
)

//FIXME: cannot explicit set "0" interval
func MustPositiveInterval(interval string, dft time.Duration) time.Duration {
	if interval == "" {
		return dft
	}

	d, err := time.ParseDuration(interval)
	if d < 0 || err != nil {
		return dft
	}
	return d
}
