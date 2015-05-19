package newcore

import (
	"fmt"
	"time"
)

// CollectorFactory returns a collector
func CollectorFactory(name string) Collector {

	return NewCollector(name, time.Duration(10)*time.Millisecond)
}

type dummy_collector struct {
	name string // collector name

	items    MultiDataPoint
	interval time.Duration
	enabled  bool
}

// this is the function actually collector data
func (c *dummy_collector) collect() error {
	for i := 0; i < 100; i++ {
		c.items = append(c.items, &DataPoint{
			Metric:    fmt.Sprintf("metric.%s", c.name),
			Timestamp: time.Now(),
			Value:     1,
			Tags:      nil,
			Meta:      nil,
		})
	}
	return nil
}

// NewCollector returns a Collector for uri.
func NewCollector(name string, interval time.Duration) Collector {
	f := &dummy_collector{
		name:     name,
		enabled:  true,
		interval: interval,
	}
	return f
}

func (f *dummy_collector) Name() string {
	return f.name
}

func (f *dummy_collector) Close() error {
	return nil
}

func (f *dummy_collector) ClassName() string {
	return "dummy_collector"
}

func (f *dummy_collector) IsEnabled() bool {
	return f.enabled
}

func (f *dummy_collector) Interval() time.Duration {
	return f.interval
}

func (f *dummy_collector) CollectOnce() *CollectResult {
	if err := f.collect(); err != nil {
		return &CollectResult{
			Collected: nil,
			Next:      time.Now().Add(f.interval),
			Err:       err,
		}
	}

	items := f.items
	f.items = nil

	return &CollectResult{
		Collected: &items,
		Next:      time.Now().Add(f.interval),
		Err:       nil,
	}
}
