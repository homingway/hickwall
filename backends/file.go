// write MultiDataPoint into a file.
// no rotation currently

package backends

import (
	"bytes"
	"fmt"
	"github.com/oliveagle/hickwall/config"
	"github.com/oliveagle/hickwall/newcore"
	"log"
	"os"
	"path/filepath"
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

	// file backend specific attributes
	output *os.File
	conf   *config.Transport_file
}

func NewFileBackend(name string, conf *config.Transport_file) newcore.Publication {
	s := &fileBackend{
		name:    name,
		closing: make(chan chan error),
		updates: make(chan *newcore.MultiDataPoint),
		conf:    conf,
	}

	go s.loop()
	return s
}

func (b *fileBackend) loop() {
	var (
		startConsuming     <-chan *newcore.MultiDataPoint
		try_open_file_once chan bool
		try_open_file_tick <-chan time.Time
		buf                = bytes.NewBuffer(make([]byte, 0, 1024))
	)
	startConsuming = b.updates
	log.Println("filebackend.loop started")

	for {
		if b.output == nil && try_open_file_once == nil && try_open_file_tick == nil {
			startConsuming = nil // disable consuming
			try_open_file_once = make(chan bool)
			// log.Println("try to open file the first time.")

			// try to open file the first time async.
			go func() {
				err := b.openFile()

				if b.output != nil && err == nil {
					// log.Println("openFile first time OK", b.output)
					try_open_file_once <- true
				} else {
					log.Printf("CRITICAL: filebackend trying to open file but failed: %s", err)
					try_open_file_once <- false
				}
			}()
		}

		select {
		case md := <-startConsuming:
			// fmt.Println("start consuming ")
			for _, p := range *md {
				if b.output != nil {
					// fmt.Printf("fileBackend.loop name:%s, consuming md: 0x%X \n", b.name, &md)
					// fmt.Println(p.Metric)
					fmt.Fprintf(buf, "%+v\n", p)
					b.output.Write(buf.Bytes())
					buf.Reset()
				}
			}

		case opened := <-try_open_file_once:
			try_open_file_once = nil // disable this branch
			if !opened {
				// failed open it the first time,
				// then we try to open file with time interval, until opened successfully.
				log.Println("open the first time failed, try to open with interval of 1s")
				try_open_file_tick = time.Tick(time.Second * 1)
			} else {
				// log.Println("file opened the first time.", b.output, try_open_file_once, try_open_file_tick)
				startConsuming = b.updates
			}
		case <-try_open_file_tick:
			// try to open with interval
			err := b.openFile()
			if b.output != nil && err == nil {
				// finally opened.
				try_open_file_tick = nil
				startConsuming = b.updates
			} else {
				log.Printf("CRITICAL: filebackend trying to open file but failed: %s", err)
			}
		case errc := <-b.closing:
			// fmt.Println("errc <- b.closing")
			log.Println("filebackend.loop closing")
			startConsuming = nil // stop comsuming
			errc <- nil
			close(b.updates)
			log.Println("filebackend.loop stopped")
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
