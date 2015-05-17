package newcore

import (
	"fmt"
	"time"
)

// CollectorFactory returns a collector
func CollectorFactory(name string) Collector {

	return NewCollector(name)
}

type collector struct {
	basename string // which collector type is
	name     string // collector name

	items    MultiDataPoint
	interval time.Duration
	enabled  bool
}

// this is the function actually collector data
func (c *collector) collect() error {
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
func NewCollector(name string) Collector {
	f := &collector{
		basename: "first_collector_type",
		name:     name,
		enabled:  true,
		// collect:  collect,
		interval: time.Duration(1) * time.Millisecond,
		// interval: time.Duration(1) * time.Millisecond,
		// interval: time.Duration(100) * time.Microsecond,
	}
	return f
}

func (f *collector) Name() string {
	return f.name
}

func (f *collector) Close() error {
	return nil
}

func (f *collector) BaseName() string {
	return f.basename
}

func (f *collector) IsEnabled() bool {
	return f.enabled
}

func (f *collector) Interval() time.Duration {
	return f.interval
}

func (f *collector) CollectOnce() *CollectResult {
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
