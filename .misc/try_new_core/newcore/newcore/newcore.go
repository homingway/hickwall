package newcore

import "time"

var (
	timestamp time.Time
)

func init() {
	go func() {
		for t := range time.Tick(time.Second) {
			timestamp = t
		}
	}()
}

func Now() time.Time {
	return timestamp
}
