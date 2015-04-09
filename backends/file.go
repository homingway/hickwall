package backends

import (
	// "fmt"
	"github.com/oliveagle/hickwall/collectorlib"
	"github.com/oliveagle/hickwall/config"
	log "github.com/oliveagle/seelog"
	"os"
	"path/filepath"
	"time"
)

type FileWriter struct {
	tick   <-chan time.Time
	mdCh   chan collectorlib.MultiDataPoint
	conf   *config.Transport_file
	output *os.File
}

func NewFileWriter(conf config.Transport_file) *FileWriter {

	var default_interval = time.Duration(10) * time.Millisecond

	interval, err := collectorlib.ParseInterval(conf.Flush_Interval)
	if err != nil {
		log.Errorf("cannot parse interval of FileWriter: %s - %v", conf.Flush_Interval, err)
		interval = default_interval
	}

	abspath, _ := filepath.Abs(conf.Path)
	dir, _ := filepath.Split(abspath)
	err = os.MkdirAll(dir, 0770)
	if err != nil {
		log.Errorf("cannot create directory: %v", err)
	}

	f, err := os.OpenFile(abspath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0660)
	if err != nil {
		log.Criticalf("FileWriter cannot open file: %v  path:%s", err, conf.Path)
		f = nil
	}

	return &FileWriter{
		tick:   time.Tick(interval),
		conf:   &conf,
		output: f,
		mdCh:   make(chan collectorlib.MultiDataPoint),
	}
}

func (w *FileWriter) Enabled() bool {
	return config.Conf.Transport_file.Enabled
}

func (w *FileWriter) Close() {
	if w.output != nil {
		w.output.Close()
	}

}

func (w *FileWriter) Write(md collectorlib.MultiDataPoint) {
	if w.Enabled() == true {
		w.mdCh <- md
	}
}

func (w *FileWriter) Run() {
	for {
		select {
		case md := <-w.mdCh:
			for _, p := range md {
				json, _ := p.MarshalJSON()
				if w.output != nil {
					w.output.WriteString(string(json) + "\n")
				}

			}
		case <-w.tick:
			if w.output != nil {
				w.output.Sync()
			}
		}
	}
}
