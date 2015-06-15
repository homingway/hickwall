package collectors

import (
	"bytes"
	"github.com/oliveagle/hickwall/config"
	"testing"
)

var configs = map[string][]byte{
	"collector_ping": []byte(`
groups:
    -
        prefix: "prefix"
        collector_ping:
            -
                # in seconds
                interval: 10s
                metric_key: "ping"
                tags: {
                    "some": "test2"
                }
                targets:
                    - "www.baidu.com"
                    - "www.12306.com"
                timeout: 50ms
                packets: 5`),
	// cmd disabled for safety
	//	"collector_cmd": []byte(`
	//groups:
	//    -
	//        prefix: "prefix"
	//        collector_cmd:
	//            -
	//                cmd:
	//                    - 'c:\python27\python.exe'
	//                    - 'D:\Users\rhtang\oledev\gocodez\src\github.com\oliveagle\hickwall\misc\collector_cmd.py'
	//                interval: 1s`),
	"collector_pdh_win": []byte(`
groups:
    -
        prefix: "prefix"
        collector_win_pdh:
            -
                interval: 2s
                tags: {
                    "bu": "train"
                }
                queries:
                    -
                        query: "\\System\\Processes"
                        metric: "win.pdh.process_cnt"
                        # metric: "win.processes.count"     duplicated metric key: win.processes.count
                    -
                        query: "\\Memory\\Available Bytes"
                        metric: "win.pdh.memory.available_bytes"

            -
                interval: 2s
                tags: {
                    "bu": "train"
                }
                queries:
                    -
                        query: "\\System\\Processes"
                        metric: "win.pdh.process_cnt_1"
                        tags: {
                            "mount": "C",
                            "prodution": "中文",
                        }
                        #TODO: support meta
                        # meta: {
                        #     "unit": "bytes"
                        # }
                    -
                        query: "\\Memory\\Available Bytes"
                        metric: "win.pdh.memory.available_bytes_1"
                        tags: {
                            "mount": "C"
                        }`),
	"collector_wmi_win": []byte(`
groups:
    -
        prefix: "prefix"
        collector_win_wmi:
            -
                interval: 2s
                tags: {
                    "bu": "train",
                    "prodution": "短周期"
                }

                queries:
                    -
                        query: "select Name, FileSystem, FreeSpace, Size from Win32_LogicalDisk where MediaType=11 or mediatype=12"
                        metrics:
                            -
                                value_from: "Size"
                                metric: "win.wmi.fs.size.bytes"
                                tags: {
                                    "mount": "{{.Name}}",
                                }
                            -
                                value_from: "FreeSpace"
                                metric: "win.wmi.fs.freespace.bytes"
                                tags: {
                                    "mount": "{{.Name}}",
                                }

            -
                interval: 5s
                tags: {
                    "bu": "train",
                    "prodution": "长周期"
                }
                queries:
                    -
                        query: "select Name, NumberOfCores from Win32_Processor"
                        metrics:
                            -
                                value_from: "Name"
                                metric: "win.wmi.cpu.name"
                            -
                                value_from: "NumberOfCores"
                                metric: "win.wmi.cpu.numberofcores"`),
}

func Test_UseConfigCreateCollectors_1(t *testing.T) {
	for key, data := range configs {
		rconf, err := config.ReadRuntimeConfig(bytes.NewBuffer(data))
		if err != nil {
			t.Errorf("%s, err %v", key, err)
			return
		}

		clrs, err := UseConfigCreateCollectors(rconf)
		if err != nil {
			t.Errorf("%s, err %v", key, err)
		}
		if len(clrs) <= 0 {
			t.Errorf("%s, nothing created", key)
		}
		t.Logf("%s - %v", key, clrs)
	}
}

var fail_configs = map[string][]byte{
	"no_prefix": []byte(`
groups:
    -
        prefix: ""
        collector_cmd:
            -
                cmd:
                    - 'c:\python27\python.exe'
                    - 'D:\Users\rhtang\oledev\gocodez\src\github.com\oliveagle\hickwall\misc\collector_cmd.py'
                interval: 1s`),
	"duplicated_prefix": []byte(`
groups:
    -
        prefix: "prefix"
        collector_cmd:
            -
                cmd:
                    - 'c:\python27\python.exe'
                    - 'D:\Users\rhtang\oledev\gocodez\src\github.com\oliveagle\hickwall\misc\collector_cmd.py'
                interval: 1s
    -
        prefix: "prefix"
        collector_cmd:
            -
                cmd:
                    - 'c:\python27\python.exe'
                    - 'D:\Users\rhtang\oledev\gocodez\src\github.com\oliveagle\hickwall\misc\collector_cmd.py'
                interval: 1s`),
}

func Test_UseConfigCreateCollectors_Fails(t *testing.T) {
	for key, data := range fail_configs {
		rconf, err := config.ReadRuntimeConfig(bytes.NewBuffer(data))
		if err != nil {
			t.Errorf("%s, err %v", key, err)
			return
		}

		clrs, err := UseConfigCreateCollectors(rconf)
		if err == nil {
			t.Errorf("%s should fail but not", key)
		}
		if len(clrs) > 0 {
			t.Errorf("%s should fail but not", key)
		}
		t.Logf("%s - %v", key, clrs)
	}
}

var sub_configs = map[string][]byte{
	"1": []byte(`
client:
    subscribe_kafka:
        - 
            Name: "kafka1"
            Topic: "test"
            Broker_list:
                - "10.1.1.1:9092"
`),
}

func Test_UseConfigCreateSubscription(t *testing.T) {
	for key, data := range sub_configs {
		rconf, err := config.ReadRuntimeConfig(bytes.NewBuffer(data))
		if err != nil {
			t.Errorf("%s, err %v", key, err)
			return
		}

		subs, err := UseConfigCreateSubscription(rconf)
		if err != nil {
			t.Errorf("%s should not fail: %v", key, err)
		}
		if len(subs) != 1 {
			t.Errorf("%s should not fail", key)
		}
		t.Logf("%s - %v", key, subs)
	}
}

var sub_fails_configs = map[string][]byte{
	"empty_name": []byte(`
client:
    subscribe_kafka:
        - 
            Name: ""
            Topic: "test"
            Broker_list:
                - "10.1.1.1:9092"
`),
	"duplicated_name": []byte(`
client:
    subscribe_kafka:
        - 
            Name: "A"
            Topic: "test"
            Broker_list:
                - "10.1.1.1:9092"
        - 
            Name: "A"
            Topic: "test1"
            Broker_list:
                - "10.1.1.1:9092"
`),
}

func Test_UseConfigCreateSubscription_fails(t *testing.T) {
	for key, data := range sub_fails_configs {
		rconf, err := config.ReadRuntimeConfig(bytes.NewBuffer(data))
		if err != nil {
			t.Errorf("%s, err %v", key, err)
			return
		}

		subs, err := UseConfigCreateSubscription(rconf)
		if err == nil {
			t.Errorf("%s should fail but not", key)
		}
		if len(subs) > 0 {
			t.Errorf("%s should fail but not", key)
		}
		t.Logf("%s - %v", key, subs)
	}
}
