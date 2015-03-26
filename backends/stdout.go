package backends

import (
	"fmt"
	"github.com/oliveagle/boltq"
	"github.com/oliveagle/hickwall/collectorlib"
	"log"
	"sync"
	"time"
)

type StdoutWriter struct {
	tick    <-chan time.Time
	tickBkf <-chan time.Time
	mdCh    chan collectorlib.MultiDataPoint
	buf     collectorlib.MultiDataPoint
	lock    sync.RWMutex
	conf    StdoutWriterConf
	q       *boltq.BoltQ
}

type StdoutWriterConf struct {
	Enabled                       bool
	Interval                      time.Duration
	Backfill_interval             time.Duration
	Max_batch_size                int
	Backfill_enabled              bool
	Backfill_handsoff             bool
	Backfill_latency_threshold_ms int
	Backfill_cool_down_second     int
}

func NewStdoutWriter(conf StdoutWriterConf) *StdoutWriter {
	q, err := boltq.NewBoltQ("backend_stdout.queue", MAX_QUEUE_SIZE, boltq.POP_ON_FULL)
	if err != nil {
		log.Panicf("cannot open backend_stdout.queue: %v", err)
	}
	return &StdoutWriter{
		conf:    conf,
		tick:    time.Tick(conf.Interval),
		tickBkf: time.Tick(conf.Backfill_interval),
		mdCh:    make(chan collectorlib.MultiDataPoint),
		buf:     collectorlib.MultiDataPoint{},
		q:       q,
	}
}

func (w *StdoutWriter) Enabled() bool {
	return w.conf.Enabled
}

func (w *StdoutWriter) Close() {
	w.flushToQueue()
}

func (w *StdoutWriter) Write(md collectorlib.MultiDataPoint) {
	if w.Enabled() == true {
		w.mdCh <- md
	}
}

func (w *StdoutWriter) Run() {
	for {
		select {
		case md := <-w.mdCh:
			w.addMD2Buf(md)
		case <-w.tick:
			w.flushToQueue()
			w.consume()
		case <-w.tickBkf:
			w.backfill()
		}
	}
}

// --------------------------------------------------------------------
//  Internal Helpers
// --------------------------------------------------------------------

// internal buf, which holds metrics for a very short period. such as 1 second.
func (w *StdoutWriter) addMD2Buf(md collectorlib.MultiDataPoint) {
	if w.Enabled() == false {
		return
	}

	w.lock.Lock()
	defer w.lock.Unlock()

	for _, p := range md {
		w.buf = append(w.buf, p)
		// fmt.Println("len(w.buf): ", len(w.buf))
		if len(w.buf) >= w.conf.Max_batch_size {
			// fmt.Println("make it a batch")
			md1 := collectorlib.MultiDataPoint(w.buf[:len(w.buf)])
			MdPush(w.q, md1)

			w.buf = nil
			// fmt.Println("w.buf = nil")
		}
	}
}

func (w *StdoutWriter) flushToQueue() {
	if w.Enabled() == false {
		return
	}

	w.lock.Lock()
	defer w.lock.Unlock()
	if len(w.buf) > 0 {
		md := collectorlib.MultiDataPoint(w.buf[:len(w.buf)])
		MdPush(w.q, md)
		w.buf = nil
	}
}

func (w *StdoutWriter) consume() {
	if w.Enabled() == false {
		return
	}

	md, err := MdPop(w.q)
	if err != nil && err.Error() != "Queue is empty" {
		log.Println(err)
		return
	}

	// Do something
	fmt.Printf(" * [stdout]consuming: batch size: %d, boltQ len: %d , mdCh len: %d, buf size: %d\n", len(md), w.q.Size(), len(w.mdCh), len(w.buf))

	// when error happened during consume, md will be push back to queue again
	if err != nil {
		MdPush(w.q, md)
	}
}

func (w *StdoutWriter) backfill() {
	if w.Enabled() == false {
		return
	}

	if w.conf.Backfill_enabled == true && w.q.Size() > 0 {

		// backfill when boltq is not empty
		md, err := MdPopBottom(w.q)
		if err != nil {
			fmt.Println(err)
			return
		}

		// do something with backfilling
		fmt.Printf(" - [stdout]backfilling:  batch size: %d, boltQ len: %d\n", len(md), w.q.Size())
		// push back to queue.

		if err != nil {
			// push back to queue.
			MdPush(w.q, md)
		}
	}
}
