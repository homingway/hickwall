package collectors

import (
	"fmt"
	"github.com/oliveagle/hickwall/newcore"
	"time"
)

type dummy_collector struct {
	name     string // collector name
	interval time.Duration
	enabled  bool
	points   int
}

func NewDummyCollector(name string, interval time.Duration, points int) newcore.Collector {
	if points <= 0 {
		points = 1
	}
	c := &dummy_collector{
		name:     name,
		enabled:  true,
		interval: interval,
		points:   points,
	}
	return c
}

func (c *dummy_collector) Name() string {
	return c.name
}

func (c *dummy_collector) Close() error {
	return nil
}

func (c *dummy_collector) ClassName() string {
	return "dummy_collector"
}

func (c *dummy_collector) IsEnabled() bool {
	return c.enabled
}

func (c *dummy_collector) Interval() time.Duration {
	return c.interval
}

func (c *dummy_collector) CollectOnce() *newcore.CollectResult {
	var items newcore.MultiDataPoint

	for i := 0; i < c.points; i++ {
		items = append(items, &newcore.DataPoint{
			Metric:    newcore.Metric(fmt.Sprintf("metric.%s", c.name)),
			Timestamp: time.Now(),
			Value:     1,
			Tags:      nil,
			Meta:      nil,
		})
	}

	return &newcore.CollectResult{
		Collected: &items,
		Next:      time.Now().Add(c.interval),
		Err:       nil,
	}
}
