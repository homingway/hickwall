package newcore

import (
	"fmt"
)

var (
	_ = fmt.Sprintf("")
)

type HookBackend struct {
	closing chan chan error      // for Close
	updates chan *MultiDataPoint // for receive updates
	hook    chan *MultiDataPoint // for hook
}

func NewHookBackend() *HookBackend {
	s := &HookBackend{
		closing: make(chan chan error),
		updates: make(chan *MultiDataPoint),
		hook:    make(chan *MultiDataPoint),
	}
	go s.loop()
	return s
}

func (b *HookBackend) loop() {
	var (
		startConsuming <-chan *MultiDataPoint
	)

	startConsuming = b.updates

	for {
		select {
		case md := <-startConsuming:
			b.hook <- md
		case errc := <-b.closing:
			startConsuming = nil // stop comsuming
			errc <- nil
			close(b.updates)
			return
		}
	}
}

func (b *HookBackend) Updates() chan<- *MultiDataPoint {
	return b.updates
}

// Should consume all data within Hook(), Don't block it.
func (b *HookBackend) Hook() <-chan *MultiDataPoint {
	return b.hook
}

func (b *HookBackend) Close() error {
	errc := make(chan error)
	b.closing <- errc
	return <-errc
}

func (b *HookBackend) Name() string {
	return "hook_backend"
}
