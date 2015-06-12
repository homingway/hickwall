package collectors

import (
	"fmt"
	"github.com/oliveagle/hickwall/collectors/config"
	"github.com/oliveagle/hickwall/logging"
	// "github.com/oliveagle/hickwall/newcore"
	// "strings"
	"testing"
	"time"
)

func TestKafkaSubscription(t *testing.T) {
	logging.SetLevel("debug")

	opts := config.KafkaSubscription{
		Broker_list: []string{
			"oleubuntu:9092",
		},
		Topic_id:       "test",
		Max_batch_size: 10,
		Flush_interval: "100ms",
	}
	sub := newKafkaSubscription("kafka", opts)

	time.AfterFunc(time.Second*10, func() {
		sub.Close()
	})

	timeout := time.After(time.Second * time.Duration(11))

main_loop:
	for {
		select {
		case md, openning := <-sub.Updates():
			if openning {
				if md == nil {
					fmt.Println("md is nil")
				} else {
					fmt.Printf("count: %d\n", len(md))
					// for _, dp := range md {
					// 	fmt.Println("dp: ---> ", dp)
					// }
				}
			} else {
				break main_loop
			}
		case <-timeout:
			t.Error("timed out! something is blocking")
			break main_loop
		}
	}
}
