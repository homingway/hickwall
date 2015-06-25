package hickwall

import (
	"bytes"
	"fmt"
	"github.com/oliveagle/hickwall/config"
	"path/filepath"
	"testing"
	"time"
)

var configs = map[string][]byte{
	"file": []byte(`
client:
    transport_file:
        path: "/var/lib/hickwall/fileoutput.txt"`),
}

func Test_running_core_CreateRunningCore(t *testing.T) {
	for key, data := range configs {
		rconf, err := config.ReadRuntimeConfig(bytes.NewBuffer(data))
		if err != nil {
			t.Errorf("failed read runtimeconfig: %s, err %v", key, err)
			return
		}

		core, err := UpdateRunningCore(rconf)
		if err != nil {
			t.Errorf("failed create running core: %s, err %v", key, err)
		}
		if core == nil {
			t.Errorf("core is nil: %s, err %v", key, err)
		}
		//		t.Logf("%s - %+v", key, core)
	}
}

func Test_running_core_CreateRunningCore_Nil(t *testing.T) {
	core, err := UpdateRunningCore(nil)
	if err == nil || core != nil {
		t.Errorf("should fail but not. core: %+v", core)
		return
	}
}

// make sure heartbeat is always created
func Test_running_core_CreateRunningCore_Alwasy_Heartbeat(t *testing.T) {
	data := []byte(`
client:
    transport_dummy:
        name: "dummy"
`)

	rconf, err := config.ReadRuntimeConfig(bytes.NewBuffer(data))
	if err != nil {
		t.Errorf("failed to read runtime config. err %v", err)
		return
	}

	core, hook, err := create_running_core_hooked(rconf, true)
	if err != nil {
		t.Errorf("err %v", err)
		return
	}
	if core == nil {
		t.Errorf("err %v", err)
	}

	// closed_chan := make(chan int)
	time.AfterFunc(time.Second*1, func() {
		core.Close()
	})
	timeout := time.After(time.Second * 2)

	// t.Logf("core: %+v", core)
	// t.Logf("hook: %+v", hook)
	heartbeat_exists := false
main_loop:
	for {
		select {
		// case <-closed_chan:
		// break main_loop
		case md, opening := <-hook.Hook():
			if opening == false {
				t.Log("HookBackend closed")
				break main_loop
			}
			for _, dp := range md {
				res, _ := dp.MarshalJSON()
				fmt.Println("--> dp", string(res))
				if dp.Metric == "hickwall.client.alive" {
					heartbeat_exists = true
					break main_loop
				}
			}
		case <-timeout:
			t.Error("something is blocking")

			break main_loop
		}
	}
	if heartbeat_exists == false {
		t.Error("heartbeat didn't show up")
	}
}

func Test_running_core_Start_From_File(t *testing.T) {
	config.CORE_CONF_FILEPATH, _ = filepath.Abs("./test/core_config.yml")
	config.CONF_FILEPATH, _ = filepath.Abs("./test/config.yml")
	Stop() // stop if already exists while test all cases
	err := Start()
	if err != nil {
		t.Errorf("failed to Start() from file: %v", err)
	} else {
		Stop()
	}
}

// multiple core can run side by side.
// TODO: possible data lose while replace 2 cores if we support counters.
// counters is something works like this way:
//  c.Add(1)  c.Decr(1)
// if 2 cores are running. these counters will have different value.
// who to preserve those counters ??? or we just don't support counters
// internally.
func Test_running_core_MultipleCore(t *testing.T) {
	rconf, err := config.ReadRuntimeConfig(bytes.NewBuffer(configs["file"]))
	if err != nil {
		t.Errorf("err %v", err)
		return
	}

	core1, err := UpdateRunningCore(rconf)
	if err != nil {
		t.Errorf("err %v", err)
	}
	if core1 == nil {
		t.Errorf("err %v", err)
	}

	core2, err := UpdateRunningCore(rconf)
	if err != nil {
		t.Errorf("err %v", err)
	}
	if core2 == nil {
		t.Errorf("err %v", err)
	}

	t.Logf("%+v", core1)
	t.Logf("%+v", core2)
}

func Test_running_core_kafka_producer(t *testing.T) {
	config.CORE_CONF_FILEPATH, _ = filepath.Abs("./test/core_config.yml")
	config.CONF_FILEPATH, _ = filepath.Abs("./test/config_kafka_producer.yml")
	Stop() // stop if already exists while test all cases
	err := Start()
	if err != nil {
		t.Errorf("failed to Start() from file: %s", err)
	}
	done := time.After(time.Second * 1)

	for {
		select {
		case <-done:
			Stop()
			return
		}
	}
}
