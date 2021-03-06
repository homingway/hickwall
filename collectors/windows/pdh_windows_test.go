package windows

import (
	"fmt"
	"github.com/oliveagle/hickwall/collectors/config"
	"github.com/oliveagle/hickwall/newcore"
	"strings"
	"testing"
	"time"
)

func Test_pdh_windows_QueryToTag_1(t *testing.T) {
	opts := config.Config_win_pdh_collector{
		Interval: "100ms",
		Queries: []config.Config_win_pdh_query{
			{
				Query:        "\\System\\Processes",
				Metric:       "processes.1",
				Query_to_tag: true,
			}, {
				Query:  "\\System\\Processor Queue Length",
				Metric: "win.pdh.processor_queue_length",
			}},
	}

	c := MustNewWinPdhCollector("c1", "prefix", opts)
	q, _ := c.map_queries["\\System\\Processes"]
	if _, ok := q.Tags["query"]; ok == false {
		t.Error("query is not in tag")
		t.Logf("%+v", q)
	}
	q, _ = c.map_queries["\\System\\Processor Queue Length"]
	if _, ok := q.Tags["query"]; ok != false {
		t.Error("query should not in tag")
		t.Logf("%+v", q)
	}
}

func Test_pdh_windows_QueryToTag_2(t *testing.T) {
	opts := config.Config_win_pdh_collector{
		Interval:     "100ms",
		Query_to_tag: true,
		Queries: []config.Config_win_pdh_query{
			{
				Query:  "\\System\\Processes",
				Metric: "processes.1",
			}},
	}

	c := MustNewWinPdhCollector("c1", "prefix", opts)
	q, _ := c.map_queries["\\System\\Processes"]
	if _, ok := q.Tags["query"]; ok == false {
		t.Error("query is not in tag")
		t.Logf("%+v", q)
	}
}

func TestWinPdhCollector(t *testing.T) {
	opts := config.Config_win_pdh_collector{
		Interval: "100ms",
		Queries: []config.Config_win_pdh_query{
			{
				Query:  "\\System\\Processes",
				Metric: "processes.1",
			}, {
				Query:  "\\System\\Processes",
				Metric: "processes.2",
			}},
	}

	sub := newcore.Subscribe(MustNewWinPdhCollector("c1", "prefix", opts), nil)

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
					for _, dp := range md {
						fmt.Println("dp: ---> ", dp)
						if _, ok := dp.Tags["host"]; ok == false {
							t.Error("host is not in tags")
							return
						}
						if !strings.HasPrefix(dp.Metric.Clean(), "prefix.processes.") {
							t.Error("metric wrong")
							return
						}
						if dp.Value.(float64) < 10 {
							t.Error("processes count less than 10 ?")
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
