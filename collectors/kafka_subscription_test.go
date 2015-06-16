package collectors

import (
	"fmt"
	. "github.com/oliveagle/hickwall/collectors/config"
	"github.com/oliveagle/hickwall/config"
	"github.com/oliveagle/hickwall/logging"
	// "github.com/oliveagle/hickwall/newcore"
	// "strings"
	// "encoding/gob"
	// "encoding/json"
	// "io/ioutil"
	"testing"
	"time"
)

func Test_get_kafka_sub_state_offset(t *testing.T) {
	state := newKafkaSubState("test.state")
	offset := state.Offset("test", 0)
	if offset != -1 {
		t.Error("empty record should return -1")
	}
}

func Test_update_kafka_sub_state(t *testing.T) {
	state := newKafkaSubState("test.state")
	state.Update("test", 0, 121)
	if state.Length() != 1 {
		t.Error("no record saved")
	}

	offset := state.Offset("test", 0)
	if offset != 121 {
		t.Error("update or get failed")
	}
	state.Update("test", 0, 122)
	if state.Length() != 1 {
		t.Error("?? %+v", state)
	}

	offset = state.Offset("test", 0)
	if offset != 122 {
		t.Error("update or get failed")
	}

	state.Update("test", 1, 111)
	if state.Length() != 2 {
		t.Error("?? %+v", state)
	}

	offset = state.Offset("test", 1)
	if offset != 111 {
		t.Error("update or get failed")
	}
}

func Test_update_kafka_sub_state_no_change(t *testing.T) {
	state := newKafkaSubState("test.state")
	state.Update("test", 0, 122)
	if state.Changed() != true {
		t.Error("should == true")
	}
	err := state.Save()
	if err != nil {
		t.Errorf("failed save: %v", err)
		return
	}

	state.Update("test", 0, 122)
	if state.Changed() != false {
		t.Error("should == false")
	}
}

func Test_get_kafka_sub_state_partitions(t *testing.T) {
	state := newKafkaSubState("test.state")
	parts := state.Partitions("test")
	if len(parts) > 0 {
		t.Error("no partition failed")
	}

	state.Update("test", 0, 121)
	parts = state.Partitions("test")
	if len(parts) != 1 {
		t.Error("1 partition failed")
	}

	state.Update("test", 0, 1)
	parts = state.Partitions("test")
	if len(parts) != 1 {
		t.Error("1 partition failed")
	}

	state.Update("test", 1, 123)
	parts = state.Partitions("test")
	if len(parts) != 2 {
		t.Error("2 partition failed: %+v", state)
	}
	t.Log("parts: %+v", parts)
}

func Test_save_kafka_sub_state(t *testing.T) {
	var err error
	state := newKafkaSubState("kafka_sub.state")

	state.Update("test", 1, 123)
	err = state.Save()
	if err != nil {
		t.Error("failed to save")
	}
	state.Clear()
	t.Log("kafka_sub_state: %+v", state.State())
	err = state.Load()
	if err != nil {
		t.Error("failed to load")
	}

	t.Log("kafka_sub_state: %+v", state.State())
	offset := state.Offset("test", 1)
	if offset != 123 {
		t.Error("failed to get offset")
	}
}

func TestKafkaSubscription(t *testing.T) {
	logging.SetLevel("debug")
	config.SHARED_DIR = "."

	opts := Config_KafkaSubscription{
		Name: "kafka_sub",
		Broker_list: []string{
			"opsdevhdp02.qa.nt.ctripcorp.com:9092",
			// "oleubuntu:9092",
		},
		Topic:          "test",
		Max_batch_size: 10,
		Flush_interval: "100ms",
	}
	sub, err := NewKafkaSubscription(opts)
	if err != nil {
		t.Errorf("failed to create kafka subscription: %v", err)
		return
	}

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
					for _, dp := range md {
						fmt.Println("dp: ---> ", dp)
					}
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
