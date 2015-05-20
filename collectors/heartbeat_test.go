package collectors

import (
	"fmt"
	"github.com/oliveagle/hickwall/newcore"
	"testing"
	"time"
)

var (
	_ = fmt.Sprintf("")
)

func TestHeartBeat(t *testing.T) {

	sub := newcore.Subscribe(NewHeartBeat("1s"), nil)

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
						if dp.Value != 1 || dp.Metric != "hickwall.client.alive" {
							t.Error("heartbeat is broken")
							return
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
}
