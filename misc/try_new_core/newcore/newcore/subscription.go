package newcore

import (
	// "log"
	// "fmt"
	"time"
)

// Subscribe returns a new Subscription that uses collector to collected DataPoints.
func Subscribe(collector Collector) Subscription {
	s := &sub{
		collector:      collector,
		updates:        make(chan *DataPoint),           // for Updates
		closing:        make(chan chan error),           // for Close
		maxPending:     100,                             // 100 datapoing in mem
		delay_on_error: time.Duration(1) * time.Second,  // delay on collect error
		timeout:        time.Duration(10) * time.Second, // tiemout on CollectOnce
	}
	go s.loop()
	return s
}

// sub implements the Subscription interface.
type sub struct {
	collector      Collector       // collected items
	updates        chan *DataPoint // sends items to the user
	closing        chan chan error // for Close
	maxPending     int
	delay_on_error time.Duration
	timeout        time.Duration // timeout while CollectOnce
}

func (s *sub) Updates() <-chan *DataPoint {
	return s.updates
}

func (s *sub) Close() error {
	errc := make(chan error)
	s.closing <- errc
	return <-errc
}

// loop periodically fecthes Items, sends them on s.updates, and exits
// when Close is called.
// CollectOnce asynchronously.
func (s *sub) loop() {
	// type collectResult struct {
	// 	collected []*DataPoint
	// 	next      time.Time
	// 	err       error
	// }

	var (
		collectDone  chan *CollectResult // if non-nil, CollectOnce is running
		pending      MuliDataPoint
		next         time.Time
		err          error
		first        *DataPoint
		updates      chan *DataPoint
		startCollect <-chan time.Time
		collectDelay time.Duration
		now          = time.Now()
	)

	for {
		startCollect = nil
		first = nil
		updates = nil

		if now = time.Now(); next.After(now) {
			collectDelay = next.Sub(now)
		}
		if collectDone == nil && len(pending) < s.maxPending {
			startCollect = time.After(collectDelay) // enable collect case
		}

		if len(pending) > 0 {
			first = pending[0]
			updates = s.updates // enable send case
		}

		select {
		case <-startCollect:
			collectDone = make(chan *CollectResult, 1) // enable CollectOnce
		case collectDone <- s.collector.CollectOnce():
			break
		case result := <-collectDone:
			collectDone = nil

			next, err = result.Next, result.Err
			if err != nil {
				// sub default delay if error happens while collecting data
				next = time.Now().Add(s.delay_on_error)
				break
			}

			for _, item := range *result.Collected {
				pending = append(pending, item)
			}
		case errc := <-s.closing:
			errc <- err
			close(s.updates)
			return
		case updates <- first:
			pending = pending[1:]
		}
	}
}
