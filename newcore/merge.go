package newcore

// import (
// 	"github.com/oliveagle/hickwall/logging"
// )

type merge struct {
	subs    []Subscription
	updates chan MultiDataPoint
	quit    chan struct{}
	errs    chan error
}

// Merge returns a Subscription that merges the item streams from subs.
// Closing the merged subscription closes subs.
func Merge(subs ...Subscription) Subscription {

	m := &merge{
		subs:    subs,
		updates: make(chan MultiDataPoint),
		quit:    make(chan struct{}),
		errs:    make(chan error),
	}

	for _, sub := range subs {
		go func(s Subscription) {
			for {
				// var it *DataPoint
				var it MultiDataPoint
				select {
				case it = <-s.Updates():
				case <-m.quit:
					m.errs <- s.Close()
					return
				}
				select {
				case m.updates <- it:
					// for _, dp := range it {
					// 	logging.Tracef("merged datapoing: %+v", dp)
					// }
				case <-m.quit:
					m.errs <- s.Close()
					return
				}
			}
		}(sub)
	}

	return m
}

func (m *merge) Updates() <-chan MultiDataPoint {
	return m.updates
}

func (m *merge) Close() (err error) {
	close(m.quit)
	for _ = range m.subs {
		if e := <-m.errs; e != nil {
			err = e
		}
	}
	close(m.updates)
	return
}

func (m *merge) Name() string {
	return "merge"
}
