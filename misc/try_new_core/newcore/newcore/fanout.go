package newcore

import (
	"fmt"
)

// TODO: Close All Subscriptions at once.
// TODO: Close one of Subscriptions, remaining still working.
// TODO: one Subscription is not been consuming(bloking), others still working.

type fanout struct {
	subs []Subscription
}

type Backend interface {
	Updates() chan<- *MultiDataPoint
	Close() error
}

type FanOutSet interface {
	Close() error
}

func FanOut(sub Subscription, backs ...Backend) FanoutSet {

}

type stdoutBackend struct {
	name    string
	closing chan chan error // for Close
	updates chan *MultiDataPoint
}

func NewStdoutBackend(name string) Backend {
	return &stdoutBackend{
		name: name,
	}
}

func (b *stdoutBackend) loop() {

}

func (b *stdoutBackend) Updates() chan<- *MultiDataPoint {
	fmt.Println(md)
	return nil
}

func (b *stdoutBackend) Close() error {
	errc := make(chan error)
	b.closing <- errc
	return <-errc
}
