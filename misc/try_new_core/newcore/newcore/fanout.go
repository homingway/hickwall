package newcore

import (
	"fmt"
	"time"
)

// TODO: Close All Subscriptions at once.
// TODO: Close one of Subscriptions, remaining still working.
// TODO: one Subscription is not been consuming(bloking), others still working.

const (
	maxPending = 10
)

var (
	_ = fmt.Sprintf("")
	_ = time.Now()
)

type fanout struct {
	bks []Publication
	sub Subscription

	pubs []chan<- *MultiDataPoint

	closing chan chan error
	pending [][]*MultiDataPoint

	closing_list [](chan chan error)
}

func (f *fanout) Close() error {
	// fmt.Println("fanout Close");
	errc := make(chan error)
	f.closing <- errc
	// fmt.Println("fanout.Close() finished, wait to return")
	return <-errc
}

func (f *fanout) cosuming(idx int, closing chan chan error) {
	var (
		first *MultiDataPoint
		pub   chan<- *MultiDataPoint
	)
	// fmt.Printf("fanout.consuming: idx: %d, closing: 0x%X\n", idx, closing)
	for {
		first = nil

		if len(f.pending[idx]) > 0 {
			// fmt.Printf("fanout.consuming: pending[%d] len: %d\n", idx, len(f.pending[idx]))
			first = f.pending[idx][0]
			pub = f.pubs[idx] // enable pub channel
		}

		select {
		case errc := <-closing:
			// fmt.Printf("fanout.consuming, errc := <-closing: idx: %d, pub: 0x%X\n", idx, pub)
			pub = nil         // nil pub channel
			f.pubs[idx] = nil // nil pub channel
			errc <- nil       // response to closing channel
			return
		case pub <- first: // send to pub channel
			// fmt.Printf("fanout.consuming: idx: %d, first: 0x%X\n", idx, &first)
			f.pending[idx] = f.pending[idx][1:]
			pub = nil
			first = nil
		default:
			// maybe this is the performance bottleneck
			break
		}
	}
}

func (f *fanout) loop() {
	// fmt.Println("fanout.loop() start")
	var (
		startConsuming <-chan *MultiDataPoint
	)

	startConsuming = f.sub.Updates()

	for idx, _ := range f.pubs {
		closing := make(chan chan error)
		f.closing_list = append(f.closing_list, closing)
		go f.cosuming(idx, closing)
	}

	// fmt.Println("fanout.main loop")
main_loop:
	for {
		select {
		case md, opening := <-startConsuming:
			if opening == false {
				// fmt.Println("consuming channel closed, stopping")
				f.Close()
				break main_loop
			}
			// fmt.Printf("faout.loop.main_loop: md received: 0x%X\n", &md)
			for idx, _ := range f.pending {
				if len(f.pending[idx]) < maxPending {
					// fmt.Printf("faout.loop.main_loop: pending md: 0x%X\n", &md)
					f.pending[idx] = append(f.pending[idx], md)
					// fmt.Printf("faout.loop.main_loop: len: f.pending[%d] : %d\n", idx, len(f.pending[idx]))
				}
			}
		case errc := <-f.closing:
			// fmt.Println("errc := <- f.closing")

			startConsuming = nil // stop consuming from sub

			// fmt.Println("consuming bks: ", f.bks)
			for idx, bk := range f.bks {
				// fmt.Printf("in closing: closing_list[%d]: 0x%X\n", idx, f.closing_list[idx])

				// closing consuming of each backend
				consuming_errc := make(chan error)
				f.closing_list[idx] <- consuming_errc
				<-consuming_errc
				// fmt.Println("<-consuming_errc")

				// close backend.
				go func() {
					consuming_errc <- bk.Close()
				}()
				timeout := time.After(time.Duration(2) * time.Second)
			wait_bk_close:
				for {
					select {
					case <-consuming_errc:
						// fmt.Printf("INFO: backend(%s) closed. \n", bk.Name())
						break wait_bk_close
					case <-timeout:
						fmt.Printf("CRITICAL: backend(%s) is blocking the fanout closing process!\n", bk.Name())
						break wait_bk_close
					}
				}

			}
			// fmt.Println("closed all consuming bks")
			errc <- nil
			// fmt.Println("errc <- nil, break main_loop")
			break main_loop
		}
	}

	timeout := time.After(time.Duration(2) * time.Second)
	closing_sub := make(chan error)
	go func() {
		closing_sub <- f.sub.Close()
	}()
	for {
		select {
		case <-closing_sub:
			return
		case <-timeout:
			fmt.Printf("CRITICAL: Subscription(%s) is blocking the fanout closing process!\n", f.sub.Name())
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
		f.pubs = append(f.pubs, pub.Updates())
		f.pending = append(f.pending, make([]*MultiDataPoint, 0))
	}

	// fmt.Println("f.pubs: ", f.pubs)

	go f.loop()
	return f
}
