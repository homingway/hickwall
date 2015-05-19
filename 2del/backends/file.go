package backends

import (
	"bytes"
	"fmt"
	"github.com/oliveagle/hickwall/collectorlib"
	"github.com/oliveagle/hickwall/config"
	log "github.com/oliveagle/seelog"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type FileWriter struct {
	tick    <-chan time.Time
	mdCh    chan collectorlib.MultiDataPoint
	conf    *config.Transport_file
	output  *os.File
	name    string
	path    string
	once    sync.Once
	done    chan bool
	running bool
}

func NewFileWriter(name string, conf config.Transport_file) *FileWriter {

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
		tick:    time.Tick(interval),
		conf:    &conf,
		output:  f,
		mdCh:    make(chan collectorlib.MultiDataPoint),
		name:    name,
		path:    abspath,
		once:    sync.Once{},
		done:    make(chan bool),
		running: false,
	}
}

func (w *FileWriter) Enabled() bool {
	return config.GetRuntimeConf().Transport_file.Enabled
}

func (w *FileWriter) Close() {
	if w.output != nil {
		w.output.Close()
	}
	if w.running == true {
		w.done <- true
	}
}

func (w *FileWriter) Name() string {
	return w.name
}

func (w *FileWriter) Write(md collectorlib.MultiDataPoint) {
	if w.Enabled() == true {
		w.mdCh <- md
	}
}

func (w *FileWriter) createFile() {
	w.once.Do(func() {
		if w.output == nil {
			log.Debug("FileWriter is creating file.")
			f, err := os.OpenFile(w.path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0660)
			if err != nil {
				log.Criticalf("FileWriter cannot open file: %v  path:%s", err, w.path)
				f = nil
			}
			w.output = f
		}
	})
}

func (w *FileWriter) Run() {
	// var tmp string
	w.running = true
	buf := bytes.NewBuffer(make([]byte, 0, 1024))

loop:
	for {
		select {
		case md := <-w.mdCh:
			for _, p := range md {
				//json, _ := p.MarshalJSON()
				//if w.output != nil {
				//	w.output.WriteString(string(json) + "\n")
				//}
				// fmt.Println("Run Write Point: ", p)
				// log.Debugf("datapoing: %v", p)
				// json, _ := p.MarshalJSON2String()
				// log.Debugf("datapoint: %s", json)
				if w.output != nil {
					// w.output.WriteString(fmt.Sprintf("%s\n", json))
					// w.output.WriteString(json)

					fmt.Fprintf(buf, "%+v\n", p)
					w.output.Write(buf.Bytes())
					buf.Reset()

				} else {
					w.createFile()
				}
			}
		case <-w.tick:
			if w.output != nil {
				w.output.Sync()
			}
		case <-w.done:
			break loop
		}
	}
	w.running = false
	log.Info("FileBackend Closed")
}
