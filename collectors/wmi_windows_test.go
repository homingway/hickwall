package collectors

import (
	"bytes"
	"fmt"
	"github.com/oliveagle/hickwall/collectors/config"
	"github.com/oliveagle/hickwall/newcore"
	"github.com/oliveagle/viper"
	"strings"
	"testing"
	"time"
)

func TestWinWmiCollector(t *testing.T) {
	opts := config.Config_win_wmi{
		Interval: "1s",
		Queries: []config.Config_win_wmi_query{
			{Query: "select Name, FileSystem, FreeSpace, Size from Win32_LogicalDisk where MediaType=11 or mediatype=12",
				Metrics: []config.Config_win_wmi_query_metric{
					{Value_from: "Size",
						Metric: "win.wmi.fs.size.bytes",
						Tags: newcore.TagSet{
							"mount": "{{.Name}}",
						},
					},
				}},
		},
	}

	sub := newcore.Subscribe(NewWinWmiCollector("c1", "prefix", opts), nil)

	time.AfterFunc(time.Second*1, func() {
		sub.Close()
	})

	timeout := time.After(time.Second * time.Duration(2))

main_loop:
	for {
		select {
		case md, openning := <-sub.Updates():
			if openning {
				if md == nil {
					fmt.Println("md is nil")
				} else {
					is_mount_c_exists := false
					for _, dp := range *md {
						fmt.Println("dp: ---> ", dp)
						if _, ok := dp.Tags["host"]; ok == false {
							t.Error("host is not in tags")
							return
						}
						if !strings.HasPrefix(dp.Metric.Clean(), "win.wmi.") {
							t.Error("metric wrong")
							return
						}

						m, ok := dp.Tags["mount"]
						if ok && strings.ToLower(m) == "c" {
							is_mount_c_exists = true
						}
					}
					if is_mount_c_exists == false {
						t.Error("mount c is not exists")
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

func TestWinWmiCollectorViper1(t *testing.T) {

	opts_str := []byte(`
interval: 1s
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
`)
	vp := viper.New()
	var opts config.Config_win_wmi

	vp.SetConfigType("yaml")
	vp.ReadConfig(bytes.NewBuffer(opts_str))
	vp.Marshal(&opts)

	fmt.Printf("opts loaded from viper: %+v \n", opts)

	sub := newcore.Subscribe(NewWinWmiCollector("c1", "prefix", opts), nil)

	time.AfterFunc(time.Second*1, func() {
		sub.Close()
	})

	timeout := time.After(time.Second * time.Duration(2))

main_loop:
	for {
		select {
		case md, openning := <-sub.Updates():
			if openning {
				if md == nil {
					fmt.Println("md is nil")
				} else {
					is_mount_c_exists := false

					for _, dp := range *md {
						fmt.Println("dp: ---> ", dp)
						if _, ok := dp.Tags["host"]; ok == false {
							t.Error("host is not in tags")
							return
						}
						if !strings.HasPrefix(dp.Metric.Clean(), "win.wmi.") {
							t.Error("metric wrong")
							return
						}
						m, ok := dp.Tags["mount"]
						if ok && strings.ToLower(m) == "c" {
							is_mount_c_exists = true
						}
					}

					if is_mount_c_exists == false {
						t.Error("mount c is not exists")
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
