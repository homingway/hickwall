package newcore

import (
	"fmt"
)

var (
	_ = fmt.Sprintf("")
)

type stdoutBackend struct {
	name    string
	closing chan chan error      // for Close
	updates chan *MultiDataPoint // for receive updates
}

func NewStdoutBackend(name string) Publication {
	s := &stdoutBackend{
		name:    name,
		closing: make(chan chan error),
		updates: make(chan *MultiDataPoint),
	}
	go s.loop()
	return s
}

func (b *stdoutBackend) loop() {
	var (
		startConsuming <-chan *MultiDataPoint
	)

	startConsuming = b.updates

	for {
		select {
		case md := <-startConsuming:
			fmt.Printf("stdoutBackend.loop name:%s, consuming md: 0x%X \n", b.name, &md)
		// case _ = <-startConsuming:
		// 	break
		// fmt.Printf("stdoutBackend.loop name:%s, consuming md: 0x%X \n", b.name, &md)
		case errc := <-b.closing:
			// fmt.Println("errc <- b.closing")
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
	// fmt.Println("bk.Close() start")
	errc := make(chan error)
	b.closing <- errc
	// fmt.Println("bk.Closed() finished")
	return <-errc
}

func (b *stdoutBackend) Name() string {
	return b.name
}
