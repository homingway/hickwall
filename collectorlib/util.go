package collectorlib

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	// "strconv"
	"strings"
	"sync"
	"time"
	"unicode"
	"unicode/utf8"
)

var (
	// Hostname is the machine's hostname.
	Hostname string
	// FullHostname will, if false, uses the hostname upto the first ".". Run Set()
	// manually after changing.
	FullHostname bool
	timestamp    = time.Now().Unix()
	tlock        sync.Mutex
)

// Clean cleans a hostname based on the current FullHostname setting.
func Clean(s string) string {
	if !FullHostname {
		s = strings.SplitN(s, ".", 2)[0]
	}
	return strings.ToLower(s)
}

// Set sets Hostntame based on the current preferences.
func Set() {
	h, err := os.Hostname()
	if err == nil {
		if !FullHostname {
			h = strings.SplitN(h, ".", 2)[0]
		}
	} else {
		h = "unknown"
	}
	Hostname = Clean(h)
}

func init() {
	Set()

	// timestamp
	go func() {
		for t := range time.Tick(time.Second) {
			tlock.Lock()
			timestamp = t.Unix()
			tlock.Unlock()
		}
	}()
}

func Now() (t int64) {
	tlock.Lock()
	t = timestamp
	tlock.Unlock()
	return
}

// IsDigit returns true if s consists of decimal digits.
func IsDigit(s string) bool {
	r := strings.NewReader(s)
	for {
		ch, _, err := r.ReadRune()
		if ch == 0 || err != nil {
			break
		} else if ch == utf8.RuneError {
			return false
		} else if !unicode.IsDigit(ch) {
			return false
		}
	}
	return true
}

// IsAlNum returns true if s is alphanumeric.
func IsAlNum(s string) bool {
	r := strings.NewReader(s)
	for {
		ch, _, err := r.ReadRune()
		if ch == 0 || err != nil {
			break
		} else if ch == utf8.RuneError {
			return false
		} else if !unicode.IsDigit(ch) && !unicode.IsLetter(ch) {
			return false
		}
	}
	return true
}

func TSys100NStoEpoch(nsec uint64) int64 {
	nsec -= 116444736000000000
	seconds := nsec / 1e7
	return int64(seconds)
}

func ReadLine(fname string, line func(string) error) error {
	f, err := os.Open(fname)
	if err != nil {
		return err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if err := line(scanner.Text()); err != nil {
			return err
		}
	}
	return scanner.Err()
}

var (
	pat_parse_interval = regexp.MustCompile(`(\d+)(\w+)`)
)

func ParseInterval(interval string) (d time.Duration, err error) {
	d, err = time.ParseDuration(interval)
	if d <= 0 {
		err = fmt.Errorf("interval should greater than zero: %s", interval)
	}
	return
}
