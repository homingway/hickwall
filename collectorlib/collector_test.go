package collectorlib

import (
	"fmt"
	"github.com/oliveagle/go-collectors/datapoint"
	// "github.com/oliveagle/go-collectors/metadata"
	"testing"
	"time"
)

func MockF() (datapoint.MultiDataPoint, error) {
	var md datapoint.MultiDataPoint
	Add(&md, "test_collector.metric1", 1, nil, "", "", "")
	return md, nil
}

func MockInit(c *IntervalCollector, config interface{}) {
	c.Interval = time.Second * 1
	c.name = "hahahah"
}

func Test_Collector(t *testing.T) {
	// add into global collectors list without customized config
	ic := IntervalCollector{
		F:        MockF,
		Interval: time.Millisecond * 200,
		name:     "test_collector",
		init:     MockInit,
		Enable:   nil,
	}

	// Init with default config or customized configuratioin later on
	ic.Init(nil)

	fmt.Println(ic.Name())

	ch := make(chan *datapoint.DataPoint)
	go ic.Run(ch)

	done := time.After(time.Second * 3)
loop:
	for {
		select {
		case dp, err := <-ch:
			fmt.Println(dp, err)
		case <-done:
			break loop
		}
	}
}
