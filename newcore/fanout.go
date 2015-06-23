package newcore

import (
	"github.com/oliveagle/hickwall/logging"
	"time"
)

// TODO: Close All Subscriptions at once.
// TODO: Close one of Subscriptions, remaining still working.
// TODO: one Subscription is not been consuming(bloking), others still working.

const (
	maxPending = 10
)

var (
	_ = time.Now()
)

type fanout struct {
	bks          []Publication           // backends
	sub          Subscription            // subscription
	chan_pubs    []chan<- MultiDataPoint // publication channels from backends
	closing      chan chan error         // for closing
	pending      [](chan MultiDataPoint) // pending channels
	closing_list [](chan chan error)     // closing channels for backends
}

func (f *fanout) Close() error {
	errc := make(chan error)
	f.closing <- errc
	return <-errc
}

func (f *fanout) Count() int {
	return len(f.bks)
}

func (f *fanout) cosuming(idx int, closing chan chan error) {
	var (
		first   MultiDataPoint
		pub     chan<- MultiDataPoint
		pending <-chan MultiDataPoint
	)

	first = nil
	pending = nil
	pub = nil

	logging.Tracef("fanout.consuming: -started- idx: %d, closing: 0x%X\n", idx, closing)

	for {
		if pending == nil && pub == nil {
			pending = f.pending[idx] // enable read from pending chan
		}
		logging.Tracef("fanout.consuming -1- idx: %d, first: %x, pending: %x, pub: %x\n", idx, &first, pending, pub)

		select {
		case first = <-pending:
			logging.Tracef("fanout.consuming -2- idx: %d, first: %x, pending: %x, pub: %x\n", idx, &first, pending, pub)
			pending = nil          // disable read from pending chan
			pub = f.chan_pubs[idx] // enable send to pub chan
		case pub <- first:
			logging.Debugf("fanout.consuming -3- send data Finished: %s idx: %d, sent cnt: %d, pub: %x\n", f.bks[idx].Name(), idx, len(first), pub)
			pub = nil   // disable send to pub chan
			first = nil // clear first
		case errc := <-closing:
			logging.Tracef("fanout.consuming -4.Start- closing idx: %d, first: %x, pending: %x, pub: %x\n", idx, &first, pending, pub)

			pending = nil // nil startSend channel
			pub = nil

			f.chan_pubs[idx] = nil // nil pub channel

			f.pending[idx] = nil

			errc <- nil // response to closing channel

			logging.Tracef("fanout.consuming -4.End- closing idx: %d, first: %x, pending: %x, pub: %x\n", idx, &first, pending, pub)
			return
		}
	}
}

func (f *fanout) loop() {
	logging.Debug("fanout.loop() started")
	var (
		startConsuming <-chan MultiDataPoint
	)

	startConsuming = f.sub.Updates()

	for idx, _ := range f.chan_pubs {
		closing := make(chan chan error)
		f.closing_list = append(f.closing_list, closing)
		go f.cosuming(idx, closing)
	}

main_loop:
	for {
		select {
		case md, opening := <-startConsuming:
			if opening == false {
				f.Close()
				break main_loop
			}
			for idx, p := range f.pending {
				_ = idx
				if len(p) < maxPending {
					p <- md
				} else {
					logging.Warnf("fanout.loop.main_loop: pending channel is jamming: bkname: %s\n", f.bks[idx].Name())
				}
			}
		case errc := <-f.closing:
			startConsuming = nil // stop consuming from sub

			for idx, bk := range f.bks {
				// closing consuming of each backend
				consuming_errc := make(chan error)
				f.closing_list[idx] <- consuming_errc
				<-consuming_errc

				// close backend.
				go func() {
					consuming_errc <- bk.Close()
				}()
				timeout := time.After(time.Duration(1) * time.Second)
			wait_bk_close:
				for {
					select {
					case <-consuming_errc:
						break wait_bk_close
					case <-timeout:
						logging.Errorf("backend(%s) is blocking the fanout closing process!\n", bk.Name())
						break wait_bk_close
					}
				}

			}
			logging.Debug("fanout.loop() closed all consuming backends")
			errc <- nil
			break main_loop
		}
	}

	logging.Debug("fanout.loop() exit main_loop")

	timeout := time.After(time.Duration(1) * time.Second)
	closing_sub := make(chan error)
	go func() {
		closing_sub <- f.sub.Close()
	}()
	for {
		select {
		case <-closing_sub:
			logging.Debug("fanout.loop() returned")
			return
		case <-timeout:
			logging.Errorf("Subscription(%s) is blocking the fanout closing process! forced return with timeout\n", f.sub.Name())
			return
		}
	}
}

func FanOut(sub Subscription, bks ...Publication) PublicationSet {
	f := &fanout{
		sub:     sub,
		closing: make(chan chan error),
	}

	f.bks = append(f.bks, bks...)

	for _, pub := range bks {
		f.chan_pubs = append(f.chan_pubs, pub.Updates())
		f.pending = append(f.pending, make(chan MultiDataPoint, maxPending))
	}
	go f.loop()
	return f
}
