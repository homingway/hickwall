package backends

import (
	"fmt"
	"github.com/oliveagle/hickwall/backends/config"
	"github.com/oliveagle/hickwall/logging"
	"github.com/oliveagle/hickwall/newcore"
	"testing"
	"time"
)

var (
	_ = fmt.Sprintf("")
	_ = time.Now()
)

func TestInfluxdbBackend(t *testing.T) {
	logging.SetLevel("debug")

	conf := &config.Transport_influxdb{
		Version:         "0.9.0",
		URL:             "http://10.3.6.225:8086/",
		Username:        "root",
		Password:        "root",
		Database:        "mydb",
		RetentionPolicy: "test",
	}

	merge := newcore.Merge(
		newcore.Subscribe(newcore.NewDummyCollector("c1", time.Millisecond*100, 1), nil),
		newcore.Subscribe(newcore.NewDummyCollector("c2", time.Millisecond*100, 1), nil),
	)

	b1, _ := NewInfluxdbBackend("b1", conf)
	fset := newcore.FanOut(merge, b1)

	fset_closed_chan := make(chan error)

	time.AfterFunc(time.Second*time.Duration(2), func() {
		// merge will be closed within FanOut
		fset_closed_chan <- fset.Close()
	})

	timeout := time.After(time.Second * time.Duration(3))

main_loop:
	for {
		select {
		case <-fset_closed_chan:
			fmt.Println("fset closed")
			break main_loop
		case <-timeout:
			t.Error("timed out! something is blocking")
			break main_loop
		}
	}

}
