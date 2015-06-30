package backends

import (
	"fmt"
	"github.com/oliveagle/hickwall/backends/config"
	// gconfig "github.com/oliveagle/hickwall/config"
	"github.com/Shopify/sarama/mocks"
	"github.com/oliveagle/hickwall/logging"
	"github.com/oliveagle/hickwall/newcore"
	"strings"
	"testing"
	"time"
)

var (
	_ = fmt.Sprintf("")
	_ = time.Now()
)

//func Test_kafka_KafkaBackend(t *testing.T) {
//	logging.SetLevel("debug")
//
//	conf := &config.Transport_kafka{
//		Broker_list: []string{
//			"opsdevhdp02.qa.nt.ctripcorp.com:9092",
//		},
//
//		Topic_id:           "test",
//		Ack_timeout_ms:     100,
//		Flush_frequency_ms: 100,
//	}
//
//	merge := newcore.Merge(
//		newcore.Subscribe(newcore.NewDummyCollector("c1", time.Millisecond*100, 1), nil),
//		newcore.Subscribe(newcore.NewDummyCollector("c2", time.Millisecond*100, 1), nil),
//	)
//
//	b1 := MustNewKafkaBackend("b1", conf)
//
//	fset := newcore.FanOut(merge, b1)
//
//	fset_closed_chan := make(chan error)
//
//	time.AfterFunc(time.Second*time.Duration(2), func() {
//		// merge will be closed within FanOut
//		fset_closed_chan <- fset.Close()
//	})
//
//	timeout := time.After(time.Second * time.Duration(3))
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
//}

func Test_kafka_KafkaBackend_mock1(t *testing.T) {
	logging.SetLevel("error")

	conf := &config.Transport_kafka{
		Broker_list: []string{
			"localhost:9092",
		},

		Topic_id:           "test",
		Ack_timeout_ms:     100,
		Flush_frequency_ms: 100,
	}

	b1 := MustNewKafkaBackend("b1", conf)
	b1.kconf.Producer.Return.Successes = true
	mp := mocks.NewAsyncProducer(t, b1.kconf)

	mp.ExpectInputAndSucceed()
	mp.ExpectInputAndSucceed()

	b1.producer = mp

	b1.Updates() <- newcore.MultiDataPoint{
		newcore.NewDataPoint("metric1", 1, time.Now(), nil, "", "", ""),
		newcore.NewDataPoint("metric2", 2, time.Now(), nil, "", "", ""),
	}

	msg1 := <-mp.Successes()
	if msg1.Topic != "test" {
		t.Errorf("topic is not 'test'")
	}

	if key, err := msg1.Key.Encode(); string(key) != "metric1" || err != nil {
		t.Errorf("metric key")
	}

	b1.Close()
	t.Logf("%+v", msg1)
	d, err := msg1.Value.Encode()
	if err != nil {
		t.Errorf("failed to encode data")
	}
	if strings.Index(string(d), "host") < 0 {
		t.Errorf("host tag is not set.")
	}
	t.Logf("data: %s", string(d))

}

// TODO: mock failure
