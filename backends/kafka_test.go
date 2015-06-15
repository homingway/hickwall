package backends

import (
	"fmt"
	"github.com/oliveagle/hickwall/backends/config"
	// gconfig "github.com/oliveagle/hickwall/config"
	"github.com/oliveagle/hickwall/logging"
	"github.com/oliveagle/hickwall/newcore"
	"testing"
	"time"
)

var (
	_ = fmt.Sprintf("")
	_ = time.Now()
)

func TestKafkaBackend(t *testing.T) {
	logging.SetLevel("debug")

	conf := &config.Transport_kafka{
		Broker_list: []string{
			"oleubuntu:9092",
		},

		Topic_id:           "test",
		Ack_timeout_ms:     100,
		Flush_frequency_ms: 100,
	}

	merge := newcore.Merge(
		newcore.Subscribe(newcore.NewDummyCollector("c1", time.Millisecond*100, 1), nil),
		newcore.Subscribe(newcore.NewDummyCollector("c2", time.Millisecond*100, 1), nil),
	)

	b1 := MustNewKafkaBackend("b1", conf)
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
