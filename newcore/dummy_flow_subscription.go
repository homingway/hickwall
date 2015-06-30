package newcore

import (
	"fmt"
	"time"
)

type DummyFlowSubscriber struct {
	name     string // collector name
	interval time.Duration
	enabled  bool

	// dummy_flow_collector specific attributes
	points  int
	closing chan chan error
	updates chan MultiDataPoint
}

func NewDummyFlowSubscription(name string, interval time.Duration, points int) *DummyFlowSubscriber {
	if points <= 0 {
		points = 1
	}
	c := &DummyFlowSubscriber{
		name:     name,
		enabled:  true,
		interval: interval,
		points:   points,
		updates:  make(chan MultiDataPoint), // for Updates
		closing:  make(chan chan error),     // for Close
	}
	go c.loop()
	return c
}

func (c *DummyFlowSubscriber) Name() string {
	return c.name
}

func (c *DummyFlowSubscriber) Close() error {
	errc := make(chan error)
	c.closing <- errc
	return <-errc
}

func (c *DummyFlowSubscriber) loop() {
	var items MultiDataPoint
	var tick = time.Tick(c.interval)
	var out = c.updates

	for {
		select {
		case <-tick:
			if out != nil {
				for i := 0; i < c.points; i++ {
					items = append(items, DataPoint{
						Metric:    Metric(fmt.Sprintf("metric.%s", c.name)),
						Timestamp: time.Now(),
						Value:     1,
						Tags:      nil,
						Meta:      nil,
					})
				}
			}
		case out <- items:
			// c.updates <- items
			items = nil
		case errc := <-c.closing:
			// clean up collector resource.
			out = nil
			close(c.updates)
			errc <- nil
			return
		}
	}
}

func (c *DummyFlowSubscriber) Updates() <-chan MultiDataPoint {
	return c.updates
}
