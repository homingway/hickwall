package utils

import (
	"strconv"
	"time"
)

func UTCTimeFromUnixStr(timestamp string) (time.Time, error) {
	i, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		return time.Unix(0, 0), err
	}
	return time.Unix(i, 0), nil
}
