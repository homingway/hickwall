package collectors

import (
	"fmt"
	"github.com/oliveagle/hickwall/misc/try_new_core/newcore/newcore"
	"time"
)

type dummy_collector struct {
	name     string // collector name
	interval time.Duration
	enabled  bool
}

// newCollector returns a Collector for uri.
func NewDummyCollector(name string, interval time.Duration) newcore.Collector {
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

func (f *dummy_collector) CollectOnce() *newcore.CollectResult {
	var items newcore.MultiDataPoint

	for i := 0; i < 100; i++ {
		items = append(items, &newcore.DataPoint{
			Metric:    fmt.Sprintf("metric.%s", f.name),
			Timestamp: time.Now(),
			Value:     1,
			Tags:      nil,
			Meta:      nil,
		})
	}

	return &newcore.CollectResult{
		Collected: &items,
		Next:      time.Now().Add(f.interval),
		Err:       nil,
	}
}
