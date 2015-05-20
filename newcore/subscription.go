package newcore

import (
	"log"
	"time"
)

const (
	minimal_next_interval = time.Millisecond * 100
)

type SubOptions struct {
	MaxPending   int    // 1 is enough for most cases, if consumer is fast enough
	DelayOnError string // duration string, delay duration on collect error, minimal is 100ms
}

// Subscribe returns a new Subscription that uses collector to collected DataPoints.
func Subscribe(collector Collector, opt *SubOptions) Subscription {
	var delay time.Duration
	var err error

	if opt == nil {
		// log.Println("Subscribe: opt is nil, use default. 1, 5s")
		opt = &SubOptions{
			MaxPending:   1,
			DelayOnError: "5s",
		}
	} else if opt.MaxPending <= 0 {
		log.Println("opt.MaxPending is below 0, use default 1 instead")
		opt.MaxPending = 1
	}

	if delay, err = time.ParseDuration(opt.DelayOnError); delay < time.Millisecond*time.Duration(100) || err != nil {
		log.Println("opt.DelayOnError is too frequent, use default:  100ms ")
		delay = time.Duration(100) * time.Millisecond
	}

	s := &sub{
		collector:      collector,
		updates:        make(chan *MultiDataPoint), // for Updates
		closing:        make(chan chan error),      // for Close
		maxPending:     opt.MaxPending,             //
		delay_on_error: delay,                      // delay on collect error
	}
	go s.loop()
	return s
}

// sub implements the Subscription interface.
type sub struct {
	collector      Collector            // collected items
	updates        chan *MultiDataPoint // sends items to the user
	closing        chan chan error      // for Close
	maxPending     int
	delay_on_error time.Duration
}

func (s *sub) Updates() <-chan *MultiDataPoint {
	return s.updates
}

func (s *sub) Name() string {
	return s.collector.Name()
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

	var (
		collectDone  chan *CollectResult // if non-nil, CollectOnce is running
		pending      []*MultiDataPoint
		next         time.Time
		err          error
		first        *MultiDataPoint
		updates      chan *MultiDataPoint
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

		if s.collector.IsEnabled() && collectDone == nil && len(pending) < s.maxPending {
			startCollect = time.After(collectDelay) // enable collect case
		}

		if len(pending) > 0 {
			first = pending[0]
			updates = s.updates // enable send case
		}

		select {
		case <-startCollect:
			collectDone = make(chan *CollectResult, 1) // enable CollectOnce

			// TODO: add unittest for this.
			// collectOnce should be call async, otherwise, will block consuming result.
			go func() {
				collectDone <- s.collector.CollectOnce()
			}()
		// case collectDone <- s.collector.CollectOnce():
		// 	break
		case result := <-collectDone:
			log.Println("result := <- collectDone", result)
			collectDone = nil

			next, err = result.Next, result.Err
			if err != nil {
				// sub default delay if error happens while collecting data
				//TODO: add unittest for delay_on_error.
				log.Printf("ERROR: collector(%s) error: %v", s.collector.Name(), err)
				next = time.Now().Add(s.delay_on_error)
				break
			}

			//TODO: add unittest
			if next.Sub(time.Now()) < minimal_next_interval {
				next = time.Now().Add(minimal_next_interval)
			}

			//TODO: whatif result.Collected is nil ??
			if result.Collected != nil {
				pending = append(pending, result.Collected)
			}
		case errc := <-s.closing:
			// clean up collector resource.
			errc <- s.collector.Close()
			close(s.updates)
			return
		case updates <- first:
			pending = pending[1:]
		}
	}
}
