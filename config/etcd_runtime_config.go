//NOTE:  github.com/oliveagle/viper use `hickwall_inuse` branch

package config

import (
	"bytes"
	//	"errors"
	"fmt"
	"github.com/coreos/go-etcd/etcd"
	"github.com/oliveagle/hickwall/logging"
	//	"github.com/oliveagle/viper"
	//	_ "github.com/oliveagle/viper/remote"
	"time"
)

var (
	_ = fmt.Sprint("")
)

//func ConnectEtcd(machines []string) (*etcd.Client, error) {
//	if len(machines) <= 0 {
//		return errors.New("etcd machines is empty")
//	}
//	client = etcd.NewClient(machines)
//}

func getRuntimeConfFromEtcd(client *etcd.Client, etcd_path string) (*RuntimeConfig, error) {
	resp, err := client.Get(etcd_path, false, false)
	if err != nil {
		return nil, fmt.Errorf("failed to get RuntimeConfig:%v", err)
	}
	r := bytes.NewReader([]byte(resp.Node.Value))
	return ReadRuntimeConfig(r)
}

func WatchRuntimeConfFromEtcd(etcd_machines []string, etcd_path string, stop chan error) <-chan RespConfig {
	var (
		out            = make(chan RespConfig, 1)
		sleep_duration = time.Second
		// sleep_duration = time.Second * 5
	)

	if stop == nil {
		panic("stop chan is nil")
	}

	go func() {
		var (
			the_first_time = true
			watching       = false
			chGetConf      <-chan time.Time
			chWaching      <-chan time.Time
		)

		client := etcd.NewClient(etcd_machines)

		cached_conf, _ := LoadRuntimeConfFromPath(CONF_CACHE_PATH)

		watch_stop := make(chan bool, 0)

	loop:
		for {
			if watching == false && chGetConf == nil {
				if the_first_time == true {
					chGetConf = time.After(0)
				} else {
					chGetConf = time.After(sleep_duration)
				}
			}

			if watching == true && chWaching == nil {
				chWaching = time.After(sleep_duration)
			}

			select {
			case <-stop:
				logging.Debugf("stop watching etcd.")
				watch_stop <- true
				logging.Debugf("watching etcd stopped.")
				break loop
			case <-chGetConf:
				the_first_time = false
				chGetConf = nil

				tmp_conf, err := getRuntimeConfFromEtcd(client, etcd_path)
				if err != nil {
					if cached_conf != nil {
						// if failed to get config from etcd but we have a cached copy. then use
						// this cached version first.
						out <- RespConfig{cached_conf, nil}
						cached_conf = nil // cached copy only need to emit once.
					}
				} else {
					out <- RespConfig{tmp_conf, nil}
					watching = true
				}
			case <-chWaching:
				logging.Debugf("watching etcd remote config: %s, %s", etcd_machines, etcd_path)
				resp, err := client.Watch(etcd_path, 0, false, nil, watch_stop)
				if err != nil {
					logging.Errorf("watching etcd error: %v", err)
					break
				}

				r := bytes.NewReader([]byte(resp.Node.Value))
				tmp_conf, err := ReadRuntimeConfig(r)
				if err != nil {
					logging.Errorf("watching etcd. changes detected but faild to parse config: %v", err)
					break
				}

				logging.Debugf("a new config is comming")
				out <- RespConfig{tmp_conf, nil}
			}
		}
	}()
	return out
}
