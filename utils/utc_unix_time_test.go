package utils

import (
	"fmt"
	"testing"
	"time"
)

func Test_UTCTimeFromUnix(t *testing.T) {
	now := time.Now().Add(time.Second).UTC()
	//	unix_str := string(now.Unix())
	unix_str := fmt.Sprintf("%d", now.Unix())
	t.Log("unix_str: ", unix_str)
	res, err := UTCTimeFromUnixStr(unix_str)
	if err != nil {
		t.Errorf("failed to parse unix time: %v", err)
		return
	}
	if res.Unix() != now.Unix() {
		t.Error("incorrect result.")
	}
}
