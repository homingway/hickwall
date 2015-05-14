package backends

import (
	"fmt"
	"github.com/oliveagle/hickwall/collectorlib"
	"github.com/oliveagle/hickwall/config"
)

type StdoutWriter struct {
	mdCh    chan collectorlib.MultiDataPoint
	name    string
	running bool
	done    chan bool
}

func NewStdoutWriter(name string, conf config.Transport_stdout) *StdoutWriter {
	return &StdoutWriter{
		mdCh: make(chan collectorlib.MultiDataPoint),
		name: name,
		done: make(chan bool),
	}
}

func (w *StdoutWriter) Enabled() bool {
	return config.GetRuntimeConf().Transport_stdout.Enabled
}

func (w *StdoutWriter) Close() {
	if w.running == true {
		w.done <- true
	}
}

func (w *StdoutWriter) Write(md collectorlib.MultiDataPoint) {
	if w.Enabled() == true {
		w.mdCh <- md
	}
}

func (w *StdoutWriter) Name() string {
	return w.name
}

func (w *StdoutWriter) Run() {
	w.running = true
loop:
	for {
		select {
		case md := <-w.mdCh:
			for _, p := range md {
				fmt.Println(" [stdout] point ---> ", p)
			}
		case <-w.done:
			break loop
		}
	}
	w.running = false
	fmt.Println("StdoutBackend Closed")
}
