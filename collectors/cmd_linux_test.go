package collectors

import (
	"fmt"
	"github.com/oliveagle/hickwall/newcore"
	// "github.com/oliveagle/hickwall/config"
	// "log"
	// "os"
	// "path/filepath"
	"testing"
	"time"
)

func TestCmd(t *testing.T) {
	_ = fmt.Sprintf("")

	conf := config_command{
		Cmd: []string{
			"bash",
			"./tests/cmd_linux.sh",
		},
		Interval: "1s",
	}

	sub := newcore.Subscribe(NewCmdCollector("p1", conf), nil)

	time.AfterFunc(time.Second*1, func() {
		sub.Close()
	})

	timeout := time.After(time.Second * time.Duration(2))

main_loop:
	for {
		select {
		case md, openning := <-sub.Updates():
			if openning {
				if md == nil {
					fmt.Println("md is nil")
				} else {
					for _, dp := range *md {
						fmt.Println("dp: ---> ", dp)
						if dp.Value != 1.2 {
							t.Error("...")
						}
					}
				}
			} else {
				break main_loop
			}
		case <-timeout:
			t.Error("timed out! something is blocking")
			break main_loop
		}
	}

	// panic("")
}
