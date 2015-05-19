package backends

import (
	"fmt"
	"github.com/oliveagle/hickwall/newcore"
	// "time"
)

var (
	_ = fmt.Sprintf("")
)

type fileBackend struct {
	name    string
	closing chan chan error              // for Close
	updates chan *newcore.MultiDataPoint // for receive updates
}

func NewFileBackend(name string) newcore.Publication {
	s := &fileBackend{
		name:    name,
		closing: make(chan chan error),
		updates: make(chan *newcore.MultiDataPoint),
	}
	go s.loop()
	return s
}

func (b *fileBackend) loop() {
	var (
		startConsuming <-chan *newcore.MultiDataPoint
	)

	startConsuming = b.updates

	for {
		select {
		case md := <-startConsuming:
			fmt.Printf("fileBackend.loop name:%s, consuming md: 0x%X \n", b.name, &md)
		case errc := <-b.closing:
			// fmt.Println("errc <- b.closing")
			startConsuming = nil // stop comsuming
			errc <- nil
			close(b.updates)
			return
		}
	}
}

func (b *fileBackend) Updates() chan<- *newcore.MultiDataPoint {
	return b.updates
}

func (b *fileBackend) Close() error {
	// fmt.Println("bk.Close() start")
	errc := make(chan error)
	b.closing <- errc
	// fmt.Println("bk.Closed() finished")
	return <-errc
}

func (b *fileBackend) Name() string {
	return b.name
}
