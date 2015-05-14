// +build OMIT

// realmain runs the Subscribe example with a real RSS fetcher.
package main

import (
	"fmt"
	"math/rand"
	"time"
)

// An Item is a stripped-down RSS item.
type Item struct{ Title, Channel, GUID string }

// A Fetcher fetches Items and returns the time when the next fetch should be
// attempted.  On failure, Fetch returns a non-nil error.
type Fetcher interface {
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
func Subscribe(fetcher Fetcher) Subscription {
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
	fetcher Fetcher         // fetches items
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

// FakeFetch causes Fetch to use a fake fetcher instead of the real
// one.
var FakeFetch bool

// Fetch returns a Fetcher for Items from domain.
func Fetch(domain string) Fetcher {
	if FakeFetch {
		return fakeFetch(domain)
	}
	return realFetch(domain)
}

func fakeFetch(domain string) Fetcher {
	return &fakeFetcher{channel: domain}
}

type fakeFetcher struct {
	channel string
	items   []Item
}

// FakeDuplicates causes the fake fetcher to return duplicate items.
var FakeDuplicates bool

func (f *fakeFetcher) Fetch() (items []Item, next time.Time, err error) {
	now := time.Now()
	next = now.Add(time.Duration(rand.Intn(5)) * 500 * time.Millisecond)
	item := Item{
		Channel: f.channel,
		Title:   fmt.Sprintf("Item %d", len(f.items)),
	}
	item.GUID = item.Channel + "/" + item.Title
	f.items = append(f.items, item)
	if FakeDuplicates {
		items = f.items
	} else {
		items = []Item{item}
	}
	return
}

// realFetch returns a fetcher for the specified blogger domain.
func realFetch(domain string) Fetcher {
	return NewFetcher(fmt.Sprintf("http://%s/feeds/posts/default?alt=rss", domain))
}

type fetcher struct {
	uri      string
	items    []Item
	interval time.Duration
	fetch    func(uri string) error
}

// NewFetcher returns a Fetcher for uri.
func NewFetcher(uri string) Fetcher {
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

	f.interval = time.Duration(1) * time.Second
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

// TODO: in a longer talk: move the Subscribe function onto a Reader type, to
// support dynamically adding and removing Subscriptions.  Reader should dedupe.

// TODO: in a longer talk: make successive Subscribe calls for the same uri
// share the same underlying Subscription, but provide duplicate streams.

func init() {
	rand.Seed(time.Now().UnixNano())
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
	}

	// Uncomment the panic below to dump the stack traces.  This
	// will show several stacks for persistent HTTP connections
	// created by the real RSS client.  To clean these up, we'll
	// need to extend Fetcher with a Close method and plumb this
	// through the RSS client implementation.
	//
	// panic("show me the stacks")
}
