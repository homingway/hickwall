package main

import (
	"fmt"
	"time"
)

// An Item is a stripped-down RSS item.
type Item struct{ Title, Channel, GUID string }

// A Collector fetches Items and returns the time when the next fetch should be
// attempted.  On failure, Fetch returns a non-nil error.
type Collector interface {
	Fetch() (items []Item, next time.Time, err error)
}

// A Subscription delivers Items over a channel.  Close cancels the
// subscription, closes the Updates channel, and returns the last fetch error,
// if any.
type Subscription interface {
	Updates() <-chan Item
	Close() error
}

// Subscribe returns a new Subscription that uses fetcher to fetch Items.
func Subscribe(fetcher Collector) Subscription {
	s := &sub{
		fetcher: fetcher,
		updates: make(chan Item),       // for Updates
		closing: make(chan chan error), // for Close
	}
	go s.loop()
	return s
}

// sub implements the Subscription interface.
type sub struct {
	fetcher Collector       // fetches items
	updates chan Item       // sends items to the user
	closing chan chan error // for Close
}

func (s *sub) Updates() <-chan Item {
	return s.updates
}

func (s *sub) Close() error {

	errc := make(chan error)
	s.closing <- errc // HLchan
	return <-errc     // HLchan
}

// loop periodically fecthes Items, sends them on s.updates, and exits
// when Close is called.
// Fetch asynchronously.
func (s *sub) loop() {
	const maxPending = 10
	type fetchResult struct {
		fetched []Item
		next    time.Time
		err     error
	}

	var fetchDone chan fetchResult // if non-nil, Fetch is running // HL

	var pending []Item
	var next time.Time
	var err error
	for {
		var fetchDelay time.Duration
		if now := time.Now(); next.After(now) {
			fetchDelay = next.Sub(now)
		}

		var startFetch <-chan time.Time
		if fetchDone == nil && len(pending) < maxPending { // HLfetch
			startFetch = time.After(fetchDelay) // enable fetch case
		}

		var first Item
		var updates chan Item
		if len(pending) > 0 {
			first = pending[0]
			updates = s.updates // enable send case
		}

		select {
		case <-startFetch: // HLfetch
			fetchDone = make(chan fetchResult, 1) // HLfetch
			go func() {
				fetched, next, err := s.fetcher.Fetch()
				fetchDone <- fetchResult{fetched, next, err}
			}()
		case result := <-fetchDone: // HLfetch
			fetchDone = nil // HLfetch
			// Use result.fetched, result.next, result.err

			fetched := result.fetched
			next, err = result.next, result.err
			if err != nil {
				next = time.Now().Add(10 * time.Second)
				break
			}
			for _, item := range fetched {
				pending = append(pending, item)
			}
		case errc := <-s.closing:
			errc <- err
			close(s.updates)
			return
		case updates <- first:
			pending = pending[1:]
		}
	}
}

type merge struct {
	subs    []Subscription
	updates chan Item
	quit    chan struct{}
	errs    chan error
}

// Merge returns a Subscription that merges the item streams from subs.
// Closing the merged subscription closes subs.
func Merge(subs ...Subscription) Subscription {

	m := &merge{
		subs:    subs,
		updates: make(chan Item),
		quit:    make(chan struct{}),
		errs:    make(chan error),
	}

	for _, sub := range subs {
		go func(s Subscription) {
			for {
				var it Item
				select {
				case it = <-s.Updates():
				case <-m.quit: // HL
					m.errs <- s.Close() // HL
					return              // HL
				}
				select {
				case m.updates <- it:
				case <-m.quit: // HL
					m.errs <- s.Close() // HL
					return              // HL
				}
			}
		}(sub)
	}

	return m
}

func (m *merge) Updates() <-chan Item {
	return m.updates
}

func (m *merge) Close() (err error) {
	close(m.quit) // HL
	for _ = range m.subs {
		if e := <-m.errs; e != nil { // HL
			err = e
		}
	}
	close(m.updates) // HL
	return
}

// Fetch returns a Collector for Items from domain.
func Fetch(domain string) Collector {
	return realFetch(domain)
}

// realFetch returns a fetcher for the specified blogger domain.
func realFetch(domain string) Collector {
	return NewFetcher(fmt.Sprintf("http://%s/feeds/posts/default?alt=rss", domain))
}

type fetcher struct {
	uri      string
	items    []Item
	interval time.Duration
	fetch    func(uri string) error
}

// NewFetcher returns a Collector for uri.
func NewFetcher(uri string) Collector {
	f := &fetcher{
		uri: uri,
	}

	f.fetch = func(uri string) error {
		for i := 0; i < 10; i++ {
			f.items = append(f.items, Item{
				Channel: uri,
				GUID:    "guid",
				Title:   "title",
			})
		}
		return nil
	}

	// f.interval = time.Duration(1) * time.Second
	f.interval = time.Duration(100) * time.Millisecond
	return f
}

func (f *fetcher) Fetch() (items []Item, next time.Time, err error) {
	// fmt.Println("fetching", f.uri)
	if err = f.fetch(f.uri); err != nil {
		return
	}
	items = f.items
	f.items = nil

	next = time.Now().Add(f.interval)
	return
}

func main() {

	// Subscribe to some feeds, and create a merged update stream.
	merged := Merge(
		Subscribe(Fetch("blog.golang.org")),
		Subscribe(Fetch("googleblog.blogspot.com")),
		Subscribe(Fetch("googledevelopers.blogspot.com")))

	// Close the subscriptions after some time.
	time.AfterFunc(3*time.Second, func() {
		fmt.Println("closed:", merged.Close())
	})

	// Print the stream.
	for it := range merged.Updates() {
		fmt.Println(it.Channel, it.Title)
		// send data to backend writer.
	}

	panic("show me the stacks")
}
