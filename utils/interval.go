package utils

import (
	"time"
)

//FIXME: dft should also > 0
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
