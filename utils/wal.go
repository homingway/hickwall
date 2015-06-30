package utils

import (
	"fmt"
	// "io"
	"os"
	// "path/filepath"
	// "bytes"
	"io/ioutil"
	"path"
	"strings"
	"sync"
	"time"
)

type RotateWriter struct {
	lock        sync.Mutex
	filename    string // should be set to the actual filename
	Fp          *os.File
	max_size_kb int64     // max size in kb
	stop        chan bool // stop chan
}

func Split(path string) (dir, file string) {
	i := strings.LastIndex(path, "\\")
	return path[:i+1], path[i+1:]
}

func Join(elem ...string) string {
	for i, e := range elem {
		if e != "" {
			return path.Clean(strings.Join(elem[i:], ""))
		}
	}
	return ""
}

// Make a new RotateWriter. Return nil if error occurs during setup.
func NewRotateWriter(filename string, max_size_kb int64) (*RotateWriter, error) {
	var lock sync.Mutex
	w := &RotateWriter{filename: filename,
		lock:        lock,
		max_size_kb: max_size_kb,
		stop:        make(chan bool)}
	err := w.create_output(filename)
	if err != nil {
		return nil, err
	}

	go w.watching_myself(w.stop)
	return w, nil
}

func (w *RotateWriter) Close() error {
	w.stop <- true
	return nil
}

func (w *RotateWriter) ListArchives() []string {
	archives := []string{}
	dir, _ := Split(w.filename)
	// dir := w.filename[:strings.LastIndex(w.filename, "\\")+1]
	files, _ := ioutil.ReadDir(dir)
	for _, file := range files {
		if file.IsDir() {
			continue
		} else {
			archives = append(archives, Join(dir, file.Name()))
		}
	}
	return archives
}

// Write satisfies the io.Writer interface.
func (w *RotateWriter) Write(output []byte) (int, error) {
	w.lock.Lock()
	defer w.lock.Unlock()
	return w.Fp.Write(output)
}

func (w *RotateWriter) watching_myself(stop chan bool) {
	if stop == nil {
		panic("stop chan is nil")
	}

	tick := time.Tick(time.Second)
	for {
		select {
		case <-tick:
			currentSize, err := w.GetSizeKb()
			if err != nil {
				return
			}
			if currentSize > w.GetLimitKb() {
				// rotate.
				err := w.rotate()
				if err != nil {
					fmt.Println(err)
				}
			}
		case <-stop:
			return
		}
	}
}

func (w *RotateWriter) GetSizeKb() (int64, error) {
	fi, err := os.Stat(w.filename)
	if err != nil {
		return 0, err
	}
	return fi.Size() / 1024, nil
}

func (w *RotateWriter) GetLimitKb() int64 {
	return w.max_size_kb
}

func (w *RotateWriter) create_output(log_filepath string) (err error) {
	w.lock.Lock()
	defer w.lock.Unlock()

	if w.Fp != nil {
		w.Fp.Close()
		w.Fp = nil
	}

	output, err := os.OpenFile(log_filepath[:], os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	w.Fp = output
	if err != nil {
		return err
	}
	// writer = io.MultiWriter(os.Stdout, output)
	return
}

// Perform the actual act of rotating and reopening file.
func (w *RotateWriter) rotate() (err error) {
	w.lock.Lock()
	defer w.lock.Unlock()

	// Close existing file if open
	if w.Fp != nil {
		err = w.Fp.Close()
		w.Fp = nil
		if err != nil {
			return
		}
	}
	// Rename dest file if it already exists
	_, err = os.Stat(w.filename)
	if err == nil {
		err = os.Rename(w.filename, w.filename+"."+time.Now().Format("2006-01-02T15.04.05Z07.00"))
		if err != nil {
			return
		}
	}

	//remove over file
	files := w.ListArchives()
	if len(files) > 5 {
		for _, file := range files[:len(files)-5] {

			os.Remove(file)
		}
	}
	// Create a file.
	w.Fp, err = os.Create(w.filename)
	return
}
