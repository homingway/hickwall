package utils

import (
	"testing"
)

func Test_RecoverAndLog(t *testing.T) {
	defer Recover_and_log()
	panic("haha")
	t.Log("good")
}
