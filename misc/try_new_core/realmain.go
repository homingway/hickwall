package main

import (
	"fmt"
	"runtime"
	// "runtime/debug"
	"time"
)

/*
#cgo CFLAGS: -I../..
#cgo LDFLAGS: -L. -lPsapi -lpdh

#include "getPeakRSS.c"
*/
import "C"

type DataPoint struct {
	Metric    string            `json:"metric"`
	Timestamp time.Time         `json:"timestamp"`
	Value     interface{}       `json:"value"`
	Tags      map[string]string `json:"tags"`
	Meta      map[string]string `json:"meta,omitempty"`
}

type MuliDataPoint []*DataPoint

// A Collector collects Items and returns the time when the next collect should be
// attempted.  On failure, CollectOnce returns a non-nil error.
type Collector interface {
	CollectOnce() (md MuliDataPoint, next time.Time, err error)
}

// A Subscription delivers Items over a channel.  Close cancels the
// subscription, closes the Updates channel, and returns the last collect error,
// if any.
type Subscription interface {
	Updates() <-chan *DataPoint
	Close() error
}

// Subscribe returns a new Subscription that uses collector to collected DataPoints.
func Subscribe(collector Collector) Subscription {
	s := &sub{
		collector: collector,
		updates:   make(chan *DataPoint), // for Updates
		closing:   make(chan chan error), // for Close
	}
	go s.loop()
	return s
}

// sub implements the Subscription interface.
type sub struct {
	collector Collector       // collected items
	updates   chan *DataPoint // sends items to the user
	closing   chan chan error // for Close
}

func (s *sub) Updates() <-chan *DataPoint {
	return s.updates
}

func (s *sub) Close() error {

	errc := make(chan error)
	s.closing <- errc // HLchan
	return <-errc     // HLchan
}

// loop periodically fecthes Items, sends them on s.updates, and exits
// when Close is called.
// CollectOnce asynchronously.
func (s *sub) loop() {
	const maxPending = 10
	type collectResult struct {
		collected []*DataPoint
		next      time.Time
		err       error
	}

	var collectDone chan collectResult // if non-nil, CollectOnce is running // HL

	var pending []*DataPoint
	var next time.Time
	var err error
	for {
		var collectDelay time.Duration
		if now := time.Now(); next.After(now) {
			collectDelay = next.Sub(now)
		}

		var startCollect <-chan time.Time
		if collectDone == nil && len(pending) < maxPending {
			startCollect = time.After(collectDelay) // enable collect case
		}

		var first *DataPoint
		var updates chan *DataPoint
		if len(pending) > 0 {
			first = pending[0]
			updates = s.updates // enable send case
		}

		select {
		case <-startCollect:
			collectDone = make(chan collectResult, 1)
			go func() {
				collected, next, err := s.collector.CollectOnce()
				collectDone <- collectResult{collected, next, err}
			}()
		case result := <-collectDone:
			collectDone = nil

			collected := result.collected
			next, err = result.next, result.err
			if err != nil {
				next = time.Now().Add(10 * time.Second)
				break
			}
			for _, item := range collected {
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
	updates chan *DataPoint
	quit    chan struct{}
	errs    chan error
}

// Merge returns a Subscription that merges the item streams from subs.
// Closing the merged subscription closes subs.
func Merge(subs ...Subscription) Subscription {

	m := &merge{
		subs:    subs,
		updates: make(chan *DataPoint),
		quit:    make(chan struct{}),
		errs:    make(chan error),
	}

	for _, sub := range subs {
		go func(s Subscription) {
			for {
				var it *DataPoint
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

func (m *merge) Updates() <-chan *DataPoint {
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

// CollectorFactory returns a collector
func CollectorFactory(name string) Collector {
	return NewCollector(name)
}

type collector struct {
	name     string
	items    MuliDataPoint
	interval time.Duration

	// this is the function actually collector data
	collect func() error
}

// NewCollector returns a Collector for uri.
func NewCollector(name string) Collector {
	f := &collector{
		name: name,
	}

	f.collect = func() error {
		for i := 0; i < 1; i++ {
			f.items = append(f.items, &DataPoint{
				Metric:    fmt.Sprintf("metric.%s", f.name),
				Timestamp: time.Now(),
				Value:     1,
				Tags:      nil,
				Meta:      nil,
			})
		}
		return nil
	}

	// f.interval = time.Duration(1) * time.Second
	f.interval = time.Duration(1000) * time.Millisecond
	return f
}

func (f *collector) CollectOnce() (items MuliDataPoint, next time.Time, err error) {
	if err = f.collect(); err != nil {
		return
	}
	items = f.items
	f.items = nil

	next = time.Now().Add(f.interval)
	return
}

func toHuman(bytes uint64) string {
	return fmt.Sprintf("%0.0fk", float64(bytes)/1024)
}

func main() {

	// Subscribe to some feeds, and create a merged update stream.
	merged := Merge(
		Subscribe(CollectorFactory("c1")),
		Subscribe(CollectorFactory("c2")),
		Subscribe(CollectorFactory("c3")))

	// Close the subscriptions after some time.
	time.AfterFunc(600*time.Second, func() {
		// fmt.Println("closed:", merged.Close())
		merged.Close()
	})

	var a = 0

	tick := time.Tick(time.Duration(1) * time.Second)
	var mem runtime.MemStats

	fmt.Println("PeakRSS(k), CurrentRSS(k), Alloc(k), Sys(k), HeapSys(k), HeapAlloc(k), HeapInuse(k), HeapIdle(k), HeapReleased(k), HeapObjects")

	var dp *DataPoint
	var channel_closed bool

	for {
		select {
		case dp, channel_closed = <-merged.Updates():
			if dp == nil && channel_closed == false {
				fmt.Println("merged closed")
				return
			}
			// fmt.Println(dp)
			a += 1
		case <-tick:
			runtime.ReadMemStats(&mem)
			// fmt.Printf("peakRSS: %dk, curRSS:  %dk, Alloc: %s, Sys: %s, HeapSys: %s, HeapAlloc: %s, HeapInuse: %s, HeapIdle: %s, HeapObjects: %d, HeapReleased: %s \n",
			// 	C.getPeakRSS()/1024, C.getCurrentRSS()/1024, toHuman(mem.Alloc), toHuman(mem.Sys), toHuman(mem.HeapSys), toHuman(mem.HeapAlloc), toHuman(mem.HeapInuse), toHuman(mem.HeapIdle), mem.HeapObjects, toHuman(mem.HeapReleased))
			fmt.Printf("%d, %d, %d, %d, %d, %d, %d, %d, %d, %d \n",
				C.getPeakRSS()/1024, C.getCurrentRSS()/1024, mem.Alloc/1024, mem.Sys/1024, mem.HeapSys/1024, mem.HeapAlloc/1024, mem.HeapInuse/1024, mem.HeapIdle/1024, mem.HeapReleased/1024, mem.HeapObjects)
		}
	}

	// // Print the stream.
	// for dp := range merged.Updates() {
	// 	// send data to backend writer.
	// 	fmt.Println(dp)
	// 	a += 1

	// 	select {
	// 	case <-tick:
	// 		runtime.ReadMemStats(&mem)
	// 		// fmt.Printf("peakRSS: %dk, curRSS:  %dk, Alloc: %s, Sys: %s, HeapSys: %s, HeapAlloc: %s, HeapInuse: %s, HeapIdle: %s, HeapObjects: %d, HeapReleased: %s \n",
	// 		// 	C.getPeakRSS()/1024, C.getCurrentRSS()/1024, toHuman(mem.Alloc), toHuman(mem.Sys), toHuman(mem.HeapSys), toHuman(mem.HeapAlloc), toHuman(mem.HeapInuse), toHuman(mem.HeapIdle), mem.HeapObjects, toHuman(mem.HeapReleased))
	// 		fmt.Printf("%d, %d, %d, %d, %d, %d, %d, %d, %d, %d \n",
	// 			C.getPeakRSS()/1024, C.getCurrentRSS()/1024, mem.Alloc/1024, mem.Sys/1024, mem.HeapSys/1024, mem.HeapAlloc/1024, mem.HeapInuse/1024, mem.HeapIdle/1024, mem.HeapReleased/1024, mem.HeapObjects)
	// 	}
	// }

	// panic("show me the stacks")

	// On macbook, run run 300 seconds,  private memory is stabilized at 836k
}
