package main

import (
	"fmt"
	"time"
)

type Item struct {
	Title, Channel, GUID string
}

type Fetcher interface {
	Fetch() (items []Item, next time.Time, err error)
}

func Fetch(domain string) Fetcher {

}

type Subscription interface {
	Updates() <-chan Item // stream of Items
	Close() error         // shuts down the stream
}

// converts Fetches to a stream
func Subscribe(fetcher Fetcher) Subscription {
	s := &sub{
		fetcher: fetcher,
		updates: make(chan Item), // for Updates
	}

	go s.loop()
	return s
}

type sub struct {
	fetcher Fetcher   // fetches items
	updates chan Item // delivers items to the user
	closing chan chan error
}

type fetchResult struct {
	fetched []Item
	next    time.Time
	err     error
}

// loop fetches items using s.fetcher and sends them
// on s.updates.  loop exits when s.Close is called.
func (s *sub) loop() {
	var pending []Item //appened by fetch; consumed by send
	var next time.Time
	var err error
	var maxPending = 1000
	var fetchDone chan fetchResult // if non-nil, Fetch is running

	for {
		var first Item
		var updates chan Item
		if len(pending) > 0 {
			first = pending[0]
			updates = s.updates // enable send case
		}

		var fetchDelay time.Duration
		if now := time.Now(); next.After(now) {
			fetchDelay = next.Sub(now)
		}
		// startFetch := time.After(fetchDelay)
		var startFetch <-chan time.Time
		if fetchDone == nil && len(pending) < maxPending {
			startFetch = time.After(fetchDelay)
		}

		select {
		// closing ----------------------------------------
		case err_chan := <-s.closing:
			err_chan <- err
			close(s.updates) // tells receiver we're done
			return

		// fetching ---------------------------------------
		case <-startFetch:
			fetchDone = make(chan fetchResult, 1)

			go func() {
				fetched, next, err := s.fetcher.Fetch()
				fetchDone <- fetchResult{fetched, next, err}
			}()

			// if err != nil {
			// 	next = time.Now().Add(10 * time.Second)
			// 	break
			// }
			// pending = append(pending, fetched...)
		case result := <-fetchDone:
			fetchDone = nil
			pending = append(pending, fetched...)

		// sending ----------------------------------------
		case updates <- first:
			pending = pending[1:]
		}
	}

}

func (s *sub) Close() error {
	// TODO: make loop exit
	// TODO: find out about any error

	err_chan := make(chan error)
	s.closing <- err_chan
	return <-err_chan
}

func (s *sub) Updates() <-chan Item {
	return s.updates
}

// merge several streams
func Merge(subs ...Subscription) Subscription {

}

func main() {
	// Subscribe to some feeds, and create a merged update stream.
	merged := Merge(
		Subscribe(Fetch("blog.golang.org")),
		Subscribe(Fetch("www.baidu.com")),
	)

	// Close the subscriptions after some time.
	time.AfterFunc(3*time.Second, func() {
		fmt.Println("closed: ", merged.Close())
	})

	// Print the stream.
	for it := range merged.Updates() {
		fmt.Println(it.Channel, it.Title)
	}

	panic("show me the stacks")
}
