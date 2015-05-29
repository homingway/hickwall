package hickwall

import (
	"bytes"
	"fmt"
	"github.com/oliveagle/hickwall/config"
	"testing"
	"time"
)

var configs = map[string][]byte{
	"file": []byte(`
client:
    transport_file:
        path: "/var/lib/hickwall/fileoutput.txt"`),
}

func Test_CreateRunningCore(t *testing.T) {
	for key, data := range configs {
		rconf, err := config.ReadRuntimeConfig(bytes.NewBuffer(data))
		if err != nil {
			t.Errorf("%s, err %v", key, err)
			return
		}

		core, err := CreateRunningCore(rconf)
		if err != nil {
			t.Errorf("%s, err %v", key, err)
		}
		if core == nil {
			t.Errorf("%s, err %v", key, err)
		}
		t.Logf("%s - %+v", key, core)
	}
}

func Test_CreateRunningCore_Nil(t *testing.T) {
	core, err := CreateRunningCore(nil)
	if err == nil || core != nil {
		t.Errorf("%s should fail but not", err)
	}
	t.Logf("%+v", core)
}

// make sure heartbeat is always created
func Test_CreateRunningCore_Alwasy_Heartbeat(t *testing.T) {
	data := []byte(``)

	rconf, err := config.ReadRuntimeConfig(bytes.NewBuffer(data))
	if err != nil {
		t.Errorf("err %v", err)
		return
	}

	core, hook, err := create_running_core_hooked(rconf, true)
	if err != nil {
		t.Errorf("err %v", err)
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
			for _, dp := range *md {
				fmt.Println("--> dp", string(dp.Json()))
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
