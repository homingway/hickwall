package newcore

import (
	"fmt"
	"time"
)

type DummyCollector struct {
	name     string // collector name
	interval time.Duration
	enabled  bool

	// dummy_collector specific attributes
	points int
}

func dummyCollectorFactory(name string) Collector {
	return NewDummyCollector(name, time.Millisecond*100, 10)
}

func NewDummyCollector(name string, interval time.Duration, points int) *DummyCollector {
	if points <= 0 {
		points = 1
	}
	c := &DummyCollector{
		name:     name,
		enabled:  true,
		interval: interval,
		points:   points,
	}
	return c
}

func (c *DummyCollector) Name() string {
	return c.name
}

func (c *DummyCollector) Close() error {
	return nil
}

func (c *DummyCollector) ClassName() string {
	return "dummy_collector"
}

func (c *DummyCollector) IsEnabled() bool {
	return c.enabled
}

func (c *DummyCollector) Interval() time.Duration {
	return c.interval
}

func (c *DummyCollector) CollectOnce() CollectResult {
	var items MultiDataPoint

	for i := 0; i < c.points; i++ {
		items = append(items, DataPoint{
			Metric:    Metric(fmt.Sprintf("metric.%s", c.name)),
			Timestamp: time.Now(),
			Value:     1,
			Tags:      nil,
			Meta:      nil,
		})
	}

	return CollectResult{
		Collected: items,
		Next:      time.Now().Add(c.interval),
		Err:       nil,
	}
}
