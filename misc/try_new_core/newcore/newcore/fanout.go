package newcore

import (
	"fmt"
)

// TODO: Close All Subscriptions at once.
// TODO: Close one of Subscriptions, remaining still working.
// TODO: one Subscription is not been consuming(bloking), others still working.

const (
	maxPending = 10
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
	fmt.Println("fanout Close")
	errc := make(chan error)
	f.closing <- errc
	return <-errc
}

func (f *fanout) cosuming(idx int, closing chan chan error) {
	var (
		first *MultiDataPoint
		pub   chan<- *MultiDataPoint
	)

	for {
		first = nil

		if len(f.pending[idx]) > 0 {
			first = f.pending[idx][0]
			pub = f.pubs[idx] // enable pub channel
		}
		select {
		case errc := <-closing:
			pub = nil         // nil pub channel
			f.pubs[idx] = nil // nil pub channel
			errc <- nil       // response to closing channel
			return
		case pub <- first: // send to pub channel
			f.pending[idx] = f.pending[idx][1:]
		}
	}
}

func (f *fanout) Run() {
	fmt.Println("fanout Run")
	var (
		startConsuming <-chan *MultiDataPoint
	)

	startConsuming = f.sub.Updates()

	for idx, _ := range f.pubs {
		closing := make(chan chan error)
		f.closing_list = append(f.closing_list, closing)
		go f.cosuming(idx, closing)
	}

main_loop:
	for {
		select {
		case md, opening := <-startConsuming:
			if opening == false {
				fmt.Println("consuming channel closed, stopping")
				f.Close()
				break main_loop
			}
			fmt.Println("md received: ", md)
			for idx, _ := range f.pending {
				if len(f.pending[idx]) < maxPending {
					f.pending[idx] = append(f.pending[idx], md)
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
				bk.Close()
			}
			errc <- nil
			break main_loop
		}
	}

	fmt.Println("Run Closed")
}

func FanOut(sub Subscription, bks ...Publication) PublicationSet {
	var f fanout
	f.bks = append(f.bks, bks...)
	f.sub = sub

	for _, pub := range bks {
		f.pubs = append(f.pubs, pub.Updates())
		f.pending = append(f.pending, make([]*MultiDataPoint, 0))
	}

	// go f.loop()
	return &f
}

type stdoutBackend struct {
	name    string
	closing chan chan error // for Close
	updates chan *MultiDataPoint
}

func NewStdoutBackend(name string) Publication {
	return &stdoutBackend{
		name: name,
	}
}

func (b *stdoutBackend) loop() {
	var (
		startConsuming <-chan *MultiDataPoint
	)

	startConsuming = b.updates

	for {
		select {
		case md := <-startConsuming:
			fmt.Println(md)
		case errc := <-b.closing:
			startConsuming = nil // stop comsuming
			errc <- nil
			close(b.updates)
			return
		}
	}
}

func (b *stdoutBackend) Updates() chan<- *MultiDataPoint {
	return b.updates
}

func (b *stdoutBackend) Close() error {
	errc := make(chan error)
	b.closing <- errc
	return <-errc
}
