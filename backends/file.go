// write MultiDataPoint into a file.
// no rotation currently

package backends

import (
	"fmt"
	"github.com/oliveagle/hickwall/collectorlib"
	"github.com/oliveagle/hickwall/config"
	"github.com/oliveagle/hickwall/newcore"
	"log"
	"os"
	"time"
)

var (
	_ = time.Now()
	_ = fmt.Sprintf("")
)

type fileBackend struct {
	name    string
	closing chan chan error              // for Close
	updates chan *newcore.MultiDataPoint // for receive updates
	// path    string

	output *os.File
	// once       sync.Once
	flush_tick <-chan time.Time
	conf       config.Transport_file
}

func NewFileBackend(name string, conf config.Transport_file) newcore.Publication {

	var default_interval = time.Duration(10) * time.Millisecond

	interval, err := collectorlib.ParseInterval(conf.Flush_Interval)
	if err != nil {
		log.Errorf("cannot parse interval of FileWriter: %s - %v", conf.Flush_Interval, err)
		interval = default_interval
	}

	s := &fileBackend{
		name:       name,
		closing:    make(chan chan error),
		updates:    make(chan *newcore.MultiDataPoint),
		flush_tick: time.Tick(interval),
		conf:       conf,
	}

	// open output laziness
	// abspath, _ := filepath.Abs(conf.Path)
	// f, err := s.openFile(abspath)
	// if err != nil {
	// 	log.Println("CRITICAL: fileBackend cannot open file: %s, err: %v", abspath, err)
	// }
	// s.output = f

	go s.loop()
	return s
}

func (b *fileBackend) loop() {
	var (
		startConsuming     <-chan *newcore.MultiDataPoint
		try_open_file_once chan bool
		try_open_file_tick <-chan time.Time
	)

	startConsuming = b.updates

	for {
		if b.output == nil && try_open_file_once == nil && try_open_file_tick == nil {
			startConsuming = nil
			try_open_file_once = make(chan bool)
			try_open_file_tick = nil

			go func() {
				err := b.openFile()
				if b.output != nil && err == nil {
					try_open_file_once <- true
				} else {
					log.Printf("CRITICAL: filebackend trying to open file but failed: %s", err)
					try_open_file_once <- false
				}
			}()
		}

		select {
		case md := <-startConsuming:
			if b.output != nil {
				fmt.Printf("fileBackend.loop name:%s, consuming md: 0x%X \n", b.name, &md)
			}
		case created := <-try_open_file_once:
			try_open_file_once = nil // disable this branch
			if !created {
				// try open it with time interval, until opened successfully.
				try_open_file_tick = time.Tick(time.Second * 1)
			}
		case <-try_open_file_tick:
			err := b.openFile()
			if b.output != nil && err == nil {
				try_open_file_tick = nil
			} else {
				log.Printf("CRITICAL: filebackend trying to open file but failed: %s", err)
			}
		case errc := <-b.closing:
			// fmt.Println("errc <- b.closing")
			startConsuming = nil // stop comsuming
			errc <- nil
			close(b.updates)
			return
		}
	}
}

func (b *fileBackend) Updates() chan<- *newcore.MultiDataPoint {
	return b.updates
}

func (b *fileBackend) Close() error {
	// fmt.Println("bk.Close() start")
	errc := make(chan error)
	b.closing <- errc

	if b.output != nil {
		b.output.Close()
	}

	// fmt.Println("bk.Closed() finished")
	return <-errc
}

func (b *fileBackend) Name() string {
	return b.name
}

func (b *fileBackend) openFile() error {
	abspath, err := filepath.Abs(b.conf.Path)
	if err != nil {
		return fmt.Errorf("failed to get abs path: %s, err: %v", b.conf.Path, err)
	}

	f, err := os.OpenFile(abspath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0660)
	if err != nil {
		return fmt.Errorf("failed openFile: %v", err)
	}

	b.output = f
	return nil
}
