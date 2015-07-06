package utils

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	// "path"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Wal struct {
	lock        sync.Mutex
	filename    string // should be set to the actual filename
	fp          *os.File
	max_size_kb int64 // max size in kb
	max_rolls   int
	stop        chan bool // stop chan
	is_index    bool
}

var (
	islast        bool
	indexfilename string
	datafile      string
	offset        int64
)

// func Split(path string) (dir, file string) {
// 	i := strings.LastIndex(path, "\\")
// 	return path[:i+1], path[i+1:]
// }

// func Join(elem ...string) string {
// 	for i, e := range elem {
// 		if e != "" {
// 			return path.Clean(strings.Join(elem[i:], ""))
// 		}
// 	}
// 	return ""
// }

// fileExists return flag whether a given file exists
// and operation error if an unclassified failure occurs.
func fileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

//the position of string in slice
func pos(value string, slice []string) int {
	for p, v := range slice {
		if v == value {
			return p
		}
	}
	return -1
}

// Make a new wal. Return nil if error occurs during setup.
func NewWal(filename string, max_size_kb int64, max_rolls int, is_index bool) (*Wal, error) {
	var lock sync.Mutex
	w := &Wal{filename: filename,
		lock:        lock,
		max_size_kb: max_size_kb,
		max_rolls:   max_rolls,
		stop:        make(chan bool),
		is_index:    is_index,
	}
	err := w.create_output(filename)
	if err != nil {
		return nil, err
	}

	if is_index {
		indexfilename = filename + ".index"
	}
	go w.watching_myself(w.stop)
	return w, nil
}

func (w *Wal) Close() error {
	w.stop <- true
	return nil
}

func (w *Wal) ListArchives() []string {
	archives := []string{}
	dir, _ := filepath.Split(w.filename)
	// dir := w.filename[:strings.LastIndex(w.filename, "\\")+1]
	files, _ := ioutil.ReadDir(dir)
	for _, file := range files {
		if file.IsDir() {
			continue
		} else {
			archives = append(archives, filepath.Join(dir, file.Name()))
		}
	}
	sort.Sort(sort.StringSlice(archives))
	return archives
}

// Write satisfies the io.Writer interface.
func (w *Wal) Write(output []byte) (int, error) {
	w.lock.Lock()
	defer w.lock.Unlock()
	return w.fp.Write(output)
}

func (w *Wal) WriteLine(data string) (int, error) {
	return w.Write([]byte(data + "\n"))
}

func (w *Wal) watching_myself(stop chan bool) {
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

func (w *Wal) GetSizeKb() (int64, error) {
	fi, err := os.Stat(w.filename)
	if err != nil {
		return 0, err
	}
	return fi.Size() / 1024, nil
}

func (w *Wal) GetLimitKb() int64 {
	return w.max_size_kb
}

func (w *Wal) create_output(log_filepath string) (err error) {
	w.lock.Lock()
	defer w.lock.Unlock()

	if w.fp != nil {
		w.fp.Close()
		w.fp = nil
	}

	output, err := os.OpenFile(log_filepath[:], os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	w.fp = output
	if err != nil {
		return err
	}
	return
}

// Perform the actual act of rotating and reopening file.
func (w *Wal) rotate() (err error) {
	w.lock.Lock()
	defer w.lock.Unlock()

	// Close existing file if open
	if w.fp != nil {
		err = w.fp.Close()
		w.fp = nil
		if err != nil {
			return
		}
	}
	// Rename dest file if it already exists
	var newname string
	_, err = os.Stat(w.filename)
	if err == nil {
		newname = w.filename + "." + time.Now().Format("20060102_150405")
		err = os.Rename(w.filename, newname)
		if err != nil {
			return
		}
	}
	//remove over file
	files := w.ListArchives()
	if len(files) > w.max_rolls {
		for _, file := range files[:len(files)-w.max_rolls] {

			os.Remove(file)
		}
	}
	// Create a file.
	// w.fp, err = os.Create(w.filename)
	w.fp, err = os.OpenFile(w.filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)

	if w.is_index {
		datafile, _, err = w.GetIndex()
		if err == nil {
			if datafile == "" || datafile == w.filename {
				datafile = newname
			}
			err := w.SetIndex()
			if err != nil {
				fmt.Println(err)
			}
		}
	}

	return
}

// commit index.
func (w *Wal) Commit() (err error) {

	if islast == true {
		err = w.DeleteFinishedArchives()
		if err == nil {
			islast = false
		}

	}

	if exist, _ := fileExists(datafile); exist == true {
		w.SetIndex()
	}
	return
}

func (w *Wal) DeleteFinishedArchives() (err error) {
	w.lock.Lock()
	defer w.lock.Unlock()
	err = os.Remove(datafile)
	if err != nil {
		fmt.Println(err)
	}
	files := w.ListArchives()
	if len(files) >= 2 {
		datafile = files[len(files)-2]
	}
	offset = 0
	return
}

//read line refer to index
func (w *Wal) ReadLine() (line string, err error) {
	var (
		file   *os.File
		part   []byte
		prefix bool
	)

	datafile, offset, err = w.GetIndex()
	exist, err := fileExists(datafile)
	if datafile == "" || !exist {
		files := w.ListArchives()
		if len(files) >= 1 {
			datafile = files[len(files)-1]
		} else {
			datafile = w.filename
		}

	}
	if file, err = os.Open(datafile); err != nil {
		return
	}

	defer file.Close()
	file.Seek(offset, os.SEEK_CUR)
	reader := bufio.NewReader(file)
	buffer := bytes.NewBuffer(make([]byte, 1024))

	part, prefix, err = reader.ReadLine()

	buffer.Write(part)
	if !prefix {
		line = buffer.String()
		buffer.Reset()
	}
	if err == io.EOF {
		islast = true
		return
	}
	offset += int64(len(line) - 1024 + len("\n"))
	return
}

//set index
func (w *Wal) SetIndex() error {
	newindex := datafile + "|" + strconv.FormatInt(offset, 10)
	ioutil.WriteFile(indexfilename, []byte(newindex), 0)
	return nil
}

//get index
func (w *Wal) GetIndex() (string, int64, error) {
	if exist, _ := fileExists(indexfilename); exist != true {
		return "", 0, nil

	}
	buf, err := ioutil.ReadFile(indexfilename)
	if err != nil {
		return "", 0, err
	}
	index := string(buf)
	ss := strings.Split(index, "|")
	off, err := strconv.ParseInt(ss[1], 10, 64)
	file := ss[0]
	return file, off, nil
}
