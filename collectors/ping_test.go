package collectors

import (
	"github.com/oliveagle/hickwall/newcore"

	"fmt"
	"testing"
	"time"
)

func TestPing(t *testing.T) {
	_ = fmt.Sprintf("")

	conf := config_single_pinger{
		Interval: "200ms",
		Metric:   "ping",
		Timeout:  "10ms",
		Target:   "www.baidu.com",
		Packets:  5,
	}

	sub := newcore.Subscribe(NewPingCollector("p1", conf), nil)

	time.AfterFunc(time.Second*3, func() {
		sub.Close()
	})

	timeout := time.After(time.Second * time.Duration(5))

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
