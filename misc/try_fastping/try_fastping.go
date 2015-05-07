package main

import (
	"fmt"
	"github.com/tatsushid/go-fastping"
	"net"
	"os"
	// "fmt"
	"github.com/oliveagle/hickwall/collectorlib"
	. "github.com/oliveagle/hickwall/collectors"
	"github.com/oliveagle/hickwall/config"
	"time"

	"github.com/kr/pretty"
)

// func main() {
// 	p := fastping.NewPinger()
// 	ra, err := net.ResolveIPAddr("ip4:icmp", os.Args[1])
// 	if err != nil {
// 		fmt.Println(err)
// 		os.Exit(1)
// 	}
// 	p.AddIPAddr(ra)
// 	p.OnRecv = func(addr *net.IPAddr, rtt time.Duration) {
// 		fmt.Printf("IP Addr: %s receive, RTT: %v\n", addr.String(), rtt)
// 	}
// 	p.OnIdle = func() {
// 		fmt.Println("finish")
// 	}
// 	err = p.Run()
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// }

func main() {
	pretty.Println("")

	runtime_conf := config.GetRuntimeConf()

	cs := GetBuiltinCollectors()

	AddCustomizedCollectorByName("ping", "ping", runtime_conf.Collector_ping[0])
	cc := GetCustomizedCollectors()

	fmt.Println(" ++ customized_collectors:  ", cc)
	fmt.Println(" ++ builtin_collectors: ", cs)

	ch := make(chan collectorlib.MultiDataPoint)

	go cs[0].Run(ch)
	go cc[0].Run(ch)

	done := time.After(time.Second * 3)
	delay := time.After(time.Second * 1)
loop:
	for {
		select {
		case dp, err := <-ch:
			fmt.Println("MultiDataPoint: ", dp, err)
			// case <-ch:
			// fmt.Println(" point ---> ", dp, err)
			// fmt.Println("-------------------")
			// pretty.Println(dp)
			for _, p := range dp {
				fmt.Println(" point ---> ", p)
			}
		case <-delay:
			// change config on the fly
			// cs[0].Init()
			// cs[0].(*IntervalCollector).SetInterval(time.Millisecond * 200)
		case <-done:
			break loop
		}
	}
}
