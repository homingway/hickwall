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
	name     string
	items    MuliDataPoint
	interval time.Duration

	// this is the function actually collector data
	collect func() error
}

// NewCollector returns a Collector for uri.
func NewCollector(name string) Collector {
	f := &collector{
		name: name,
	}

	f.collect = func() error {
		for i := 0; i < 1; i++ {
			f.items = append(f.items, &DataPoint{
				Metric:    fmt.Sprintf("metric.%s", f.name),
				Timestamp: time.Now(),
				Value:     1,
				Tags:      nil,
				Meta:      nil,
			})
		}
		return nil
	}

	// f.interval = time.Duration(1) * time.Second
	f.interval = time.Duration(1) * time.Millisecond
	return f
}

// func (f *collector) CollectOnce() (items MuliDataPoint, next time.Time, err error) {
// 	if err = f.collect(); err != nil {
// 		return
// 	}
// 	items = f.items
// 	f.items = nil

// 	next = time.Now().Add(f.interval)
// 	return
// }

func (f *collector) IsEnabled() bool {
	return true
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
