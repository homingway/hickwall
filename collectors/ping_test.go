package collectors

import (
	"github.com/oliveagle/hickwall/collectors/config"
	"github.com/oliveagle/hickwall/newcore"

	"fmt"
	"strings"
	"testing"
	"time"
)

func TestNewPingCollectors(t *testing.T) {
	conf := config.Config_Ping{
		Interval: "200ms",
		Metric:   "ping",
		Timeout:  "100ms",
		Targets:  []string{"www.baidu.com", "www.123.com"},
		Packets:  5,
	}

	cs := NewPingCollectors("test", "prefix", conf)
	if len(cs) != 2 {
		t.Error("")
	}
}

func TestPing(t *testing.T) {
	_ = fmt.Sprintf("")

	conf := config.Config_single_pinger{
		Interval: "200ms",
		Metric:   "ping",
		Timeout:  "100ms",
		Target:   "www.baidu.com",
		Packets:  5,
	}

	sub := newcore.Subscribe(NewSinglePingCollector("p1", "prefix", conf), nil)

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
						if _, ok := dp.Tags["host"]; ok == false {
							t.Error("host is not in tags")
							return
						}
						if _, ok := dp.Tags["target"]; ok == false {
							t.Error("target is not in tags")
							return
						}

						if !strings.HasPrefix(dp.Metric.Clean(), "prefix.ping.") {
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

	// panic("")
}
