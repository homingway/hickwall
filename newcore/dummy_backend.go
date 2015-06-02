package newcore

import (
	"fmt"
	"time"
)

var (
	_ = fmt.Sprintf("")
)

type dummyBackend struct {
	name      string
	closing   chan chan error     // for Close
	updates   chan MultiDataPoint // for receive updates
	jamming   time.Duration       // jamming a little period of time while comsuming, 0 duration disable it
	printting bool                // print consuming md to stdout
	detail    bool                // if true print every datapoing
}

func MustNewDummyBackend(name string, jamming Interval, printting bool, detail bool) Publication {
	s := &dummyBackend{
		name:      name,
		closing:   make(chan chan error),
		updates:   make(chan MultiDataPoint),
		jamming:   jamming.MustDuration(time.Second),
		printting: printting,
		detail:    detail,
	}
	go s.loop()
	return s
}

func (b *dummyBackend) loop() {
	var (
		startConsuming <-chan MultiDataPoint
	)

	startConsuming = b.updates

	for {
		select {
		case md := <-startConsuming:
			if b.printting {
				if b.detail == true {
					for _, dp := range md {
						fmt.Printf("dummy(%s) --> %+v \n", b.name, dp)
					}
				} else {
					fmt.Printf("dummyBackend.loop name:%s, consuming md: 0x%X \n", b.name, &md)
				}

			}
			if b.jamming > 0 {
				fmt.Println("jamming: ", b.jamming)
				time.Sleep(b.jamming)
			}
		case errc := <-b.closing:
			// fmt.Println("errc <- b.closing")
			startConsuming = nil // stop comsuming
			errc <- nil
			close(b.updates)
			return
		}
	}
}

func (b *dummyBackend) Updates() chan<- MultiDataPoint {
	return b.updates
}

func (b *dummyBackend) Close() error {
	// fmt.Println("bk.Close() start")
	errc := make(chan error)
	b.closing <- errc
	// fmt.Println("bk.Closed() finished")
	return <-errc
}

func (b *dummyBackend) Name() string {
	return b.name
}
