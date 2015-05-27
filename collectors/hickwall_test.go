package collectors

import (
	"fmt"
	"github.com/oliveagle/hickwall/newcore"
	"strings"
	"testing"
	"time"
)

func TestHickwallCollector(t *testing.T) {
	sub := newcore.Subscribe(NewHickwallCollector("500ms"), nil)

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
						if !strings.HasPrefix(dp.Metric.Clean(), "hickwall.client.") {
							t.Error("metric wrong")
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
