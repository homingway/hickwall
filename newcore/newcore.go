package newcore

import (
	"os"
	"strings"
	"time"
)

var (
	timestamp time.Time
	hostname  string
)

func init() {
	go func() {
		// first tick
		timestamp = time.Now()

		for t := range time.Tick(time.Millisecond * 1000) {
			// for t := range time.Tick(time.Millisecond * 998) {
			timestamp = t
		}
	}()

	get_os_hostname()
}

func Now() time.Time {
	return timestamp
}

func GetHostname() string {
	return hostname
}

// Clean cleans a hostname based on the current FullHostname setting.
func clean_hostname(s string, full bool) string {
	if !full {
		s = strings.SplitN(s, ".", 2)[0]
	}
	return strings.ToLower(s)
}

// Set sets Hostntame based on the current preferences.
func get_os_hostname() {
	h, err := os.Hostname()
	if err != nil {
		h = "unknown"
	}
	hostname = clean_hostname(h, false)
}
