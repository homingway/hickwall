package backends

import (
	"bufio"
	"fmt"
	"github.com/oliveagle/hickwall/backends/config"
	//	"github.com/oliveagle/hickwall/collectors"
	"github.com/oliveagle/hickwall/newcore"
	"os"
	"testing"
	"time"
)

const (
	test_file_path = "test.txt"
)

var (
	_ = fmt.Sprintf("")
	_ = time.Now()
)

func fileLines(path string) (int, error) {
	// check how many data we collected.
	f, err := os.Open(path)
	if err != nil {
		return 0, err
	}

	fileScanner := bufio.NewScanner(f)
	lineCount := 0
	for fileScanner.Scan() {
		lineCount++
	}
	return lineCount, nil
}

func TestFileBackend(t *testing.T) {
	os.Remove(test_file_path)

	conf := &config.Transport_file{
		Enabled: true,
		Path:    test_file_path,
	}

	merge := newcore.Merge(
		newcore.Subscribe(newcore.NewDummyCollector("c1", time.Millisecond*100, 1), nil),
		newcore.Subscribe(newcore.NewDummyCollector("c2", time.Millisecond*100, 1), nil),
	)

	b1, _ := NewFileBackend("b1", conf)

	hook := newcore.NewHookBackend()
	bks := []newcore.Publication{b1, hook}

	fset := newcore.FanOut(merge, bks...)

	fset_closed_chan := make(chan error)

	time.AfterFunc(time.Second*time.Duration(1), func() {
		// merge will be closed within FanOut
		fset_closed_chan <- fset.Close()
	})

	timeout := time.After(time.Second * time.Duration(3))

	expected := 0
main_loop:
	for {
		select {
		case md, ok := <-hook.Hook():
			if ok != false {
				expected += len(md)
			} else {
				break main_loop
			}

		case <-fset_closed_chan:
			fmt.Println("fset closed")
			break main_loop
		case <-timeout:
			t.Error("timed out! something is blocking")
			break main_loop
		}
	}

	// check how many data we collected.
	lines, err := fileLines(test_file_path)
	if err != nil {
		t.Error("failed counting lines", err)
		return
	}
	// on windows this may fail!
	os.Remove(test_file_path)

	t.Logf("expected: %d, got: %d", expected, lines)
	// 1s / 100 ms = 10 batch x 1 for each x 2 collectors = 20
	//	expected := 20
	if lines != expected {
		t.Error("lines mismatch: lines: %d, expected: %d", lines, expected)
	}

}

func TestFileBackend_Failed(t *testing.T) {
	conf := &config.Transport_file{
		Enabled: true,
		Path:    "",
	}
	b1, err := NewFileBackend("b1", conf)
	if err == nil || b1 != nil {
		t.Error("create file backend should fail but not")
	}
}

//func TestFileBackendLongLoop(t *testing.T) {
//	os.Remove(test_file_path)
//
//	conf := &config.Transport_file{
//		Enabled:        true,
//		Flush_Interval: "100ms",
//		Path:           test_file_path,
//	}
//
//	merge := newcore.Merge(
//		newcore.Subscribe(collectors.NewDummyCollector("c1", time.Millisecond*100, 1), nil),
//		newcore.Subscribe(collectors.NewDummyCollector("c2", time.Millisecond*100, 1), nil),
//	)
//
//	b1, _ := NewFileBackend("b1", conf)
//	fset := newcore.FanOut(merge, b1)
//
//	fset_closed_chan := make(chan error)
//
//	time.AfterFunc(time.Second*time.Duration(100), func() {
//		// merge will be closed within FanOut
//		fset_closed_chan <- fset.Close()
//	})
//
//	timeout := time.After(time.Second * time.Duration(101))
//
//main_loop:
//	for {
//		select {
//		case <-fset_closed_chan:
//			fmt.Println("fset closed")
//			break main_loop
//		case <-timeout:
//			t.Error("timed out! something is blocking")
//			break main_loop
//		}
//	}
//
//	os.Remove(test_file_path)
//}
