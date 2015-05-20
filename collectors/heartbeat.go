package collectors

import (
	"fmt"
	"github.com/oliveagle/hickwall/newcore"
	"time"
)

var (
	_ = fmt.Sprintf("")
)

type heartbeat struct {
	interval time.Duration
}

func NewHeartBeat(interval string) newcore.Collector {
	c := &heartbeat{
		interval: newcore.NewInterval(interval).MustDuration(time.Second),
	}
	return c
}

func (c *heartbeat) Name() string {
	return "heartbeat"
}

func (c *heartbeat) Close() error {
	return nil
}

func (c *heartbeat) ClassName() string {
	return "heartbeat"
}

func (c *heartbeat) IsEnabled() bool {
	return true
}

func (c *heartbeat) Interval() time.Duration {
	return c.interval
}

func (c *heartbeat) CollectOnce() *newcore.CollectResult {
	var items newcore.MultiDataPoint

	Add(&items, "hickwall.client.alive", 1, nil, "", "", "")

	return &newcore.CollectResult{
		Collected: &items,
		Next:      time.Now().Add(c.interval),
		Err:       nil,
	}
}
