// write MultiDataPoint into a file.
// no rotation currently

package backends

import (
	"bytes"
	"fmt"
	"github.com/oliveagle/hickwall/backends/config"
	"github.com/oliveagle/hickwall/logging"
	"github.com/oliveagle/hickwall/newcore"
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
	closing chan chan error             // for Close
	updates chan newcore.MultiDataPoint // for receive updates

	// file backend specific attributes
	output *os.File
	conf   *config.Transport_file
}

func NewFileBackend(name string, conf *config.Transport_file) (*fileBackend, error) {
	s := &fileBackend{
		name:    name,
		closing: make(chan chan error),
		updates: make(chan newcore.MultiDataPoint),
		conf:    conf,
	}
	if conf.Path == "" {
		return nil, fmt.Errorf("backend file should not be empty")
	}

	go s.loop()
	return s, nil
}

func (b *fileBackend) loop() {
	var (
		startConsuming     <-chan newcore.MultiDataPoint
		try_open_file_once chan bool
		try_open_file_tick <-chan time.Time
		buf                = bytes.NewBuffer(make([]byte, 0, 1024))
	)
	startConsuming = b.updates
	logging.Debugf("filebackend.loop started")

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
					logging.Errorf("filebackend trying to open file but failed: %s", err)
					try_open_file_once <- false
				}
			}()
		}

		select {
		case md := <-startConsuming:
			for _, p := range md {
				if b.output != nil {
					res, _ := p.MarshalJSON()
					buf.Write(res)
					buf.Write([]byte("\n"))
					b.output.Write(buf.Bytes())
					buf.Reset()
				}
			}

		case opened := <-try_open_file_once:
			try_open_file_once = nil // disable this branch
			if !opened {
				// failed open it the first time,
				// then we try to open file with time interval, until opened successfully.
				logging.Error("open the first time failed, try to open with interval of 1s")
				try_open_file_tick = time.Tick(time.Second * 1)
			} else {
				logging.Debugf("file opened the first time.")
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
				logging.Errorf("filebackend trying to open file but failed: %s", err)
			}
		case errc := <-b.closing:
			logging.Debug("filebackend.loop closing")
			startConsuming = nil // stop comsuming
			errc <- nil
			close(b.updates)
			logging.Debug("filebackend.loop stopped")
			return
		}
	}
}

func (b *fileBackend) Updates() chan<- newcore.MultiDataPoint {
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
