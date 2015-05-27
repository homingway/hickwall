package newcore

import (
	"fmt"
	"time"
)

// dummyCollectorFactory returns a collector
func dummyCollectorFactory(name string) Collector {
	return newCollector(name, time.Duration(10)*time.Millisecond)
}

type dummy_collector struct {
	name     string // collector name
	interval time.Duration
	enabled  bool
}

// newCollector returns a Collector for uri.
func newCollector(name string, interval time.Duration) Collector {
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
	var items MultiDataPoint

	for i := 0; i < 100; i++ {
		items = append(items, &DataPoint{
			Metric:    fmt.Sprintf("metric.%s", f.name),
			Timestamp: time.Now(),
			Value:     1,
			Tags:      nil,
			Meta:      nil,
		})
	}

	return &CollectResult{
		Collected: &items,
		Next:      time.Now().Add(f.interval),
		Err:       nil,
	}
}
