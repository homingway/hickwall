package main

import (
	"fmt"
	"github.com/oliveagle/hickwall/collectorlib"
	"time"
)

func MockF() (collectorlib.MultiDataPoint, error) {
	var md collectorlib.MultiDataPoint
	collectorlib.Add(&md, "test_collector.metric1", 1, nil, "", "", "")
	return md, nil
}

type MockConfig struct {
	Interval time.Duration
}

func MockInit(c *collectorlib.IntervalCollector, config interface{}) {
	c.Interval = config.(MockConfig).Interval
}

func main() {
	// add into global collectors list without customized config
	ic := collectorlib.NewIntervalCollector(
		"test_collector",
		MockInit,
		MockF,
		nil,
	)

	// Init with default config or customized configuratioin later on
	conf := MockConfig{Interval: time.Second * 1}
	ic.Init(conf)

	fmt.Println("ic.Name: ", ic.Name())

	ch := make(chan collectorlib.DataPoint)
	go ic.Run(ch)

	conf2 := MockConfig{Interval: time.Millisecond * 200}

	done := time.After(time.Second * 3)
	delay := time.After(time.Second * 1)
loop:
	for {
		select {
		case dp, err := <-ch:
			fmt.Println(dp, err)
		case <-delay:
			// change config on the fly
			ic.Init(conf2)
		case <-done:
			break loop
		}
	}
}
