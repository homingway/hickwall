package backends

import (
	"fmt"
	"github.com/oliveagle/hickwall/collectorlib"
	"github.com/oliveagle/hickwall/config"
)

type StdoutWriter struct {
	mdCh chan collectorlib.MultiDataPoint
	name string
}

func NewStdoutWriter(name string, conf config.Transport_stdout) *StdoutWriter {
	return &StdoutWriter{
		mdCh: make(chan collectorlib.MultiDataPoint),
		name: name,
	}
}

func (w *StdoutWriter) Enabled() bool {
	return config.GetRuntimeConf().Transport_stdout.Enabled
}

func (w *StdoutWriter) Close() {}

func (w *StdoutWriter) Write(md collectorlib.MultiDataPoint) {
	if w.Enabled() == true {
		w.mdCh <- md
	}
}

func (w *StdoutWriter) Name() string {
	return w.name
}

func (w *StdoutWriter) Run() {
	for {
		select {
		case md := <-w.mdCh:
			for _, p := range md {
				fmt.Println(" [stdout] point ---> ", p)
			}
		}
	}
}
