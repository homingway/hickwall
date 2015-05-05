package main

import (
	"encoding/json"
	"fmt"
	"github.com/oliveagle/hickwall/config"
	"github.com/oliveagle/viper"
	// "github.com/spf13/viper"

	"github.com/coreos/go-etcd/etcd"
)

func try_watch() {
	machines := []string{"http://192.168.59.103:4001"}
	client := etcd.NewClient(machines)

	// res, err := client.Get("/config/host/DST54869.yml", false, false)
	// fmt.Println(res, err)
	// fmt.Println("res.Action", res.Action)
	// fmt.Println("res.EtcdIndex", res.EtcdIndex)
	// // fmt.Println("res.Node", res.Node)
	// fmt.Println("res.PrevNode", res.PrevNode)
	// fmt.Println("res.RaftIndex", res.RaftIndex)
	// fmt.Println("res.RaftTerm", res.RaftTerm)

	// fmt.Println("res.Node.ModifiedIndex: ", res.Node.ModifiedIndex)

	// client.Get(key, sort, recursive)

	// client.Watch("/config/host/DST54869.yml", waitIndex, recursive, receiver, stop)

	// resp, err := client.Watch("/config/host/DST54869.yml", res.Node.ModifiedIndex+1, false, nil, nil)

	for {
		fmt.Println("watching")
		resp, err := client.Watch("/config/host/DST54869.yml", 0, false, nil, nil)
		fmt.Println(resp, err)
		fmt.Println(resp.Node.ModifiedIndex, resp.Node.Key)
	}
}

func try_vip_get_remote() {
	viper.AddRemoteProvider("etcd", "http://192.168.59.103:4001", "/config/host/DST54869.yml")
	viper.SetConfigType("YAML") // because there is no file extension in a stream of bytes

	var x config.Config

	err := viper.ReadRemoteConfig()
	if err != nil {
		fmt.Println(err)
	}

	err = viper.Marshal(&x)
	if err != nil {
		fmt.Println(err)
	}

	j, _ := json.Marshal(x)

	fmt.Println(string(j))
}

func try_viper_watch_remote() {
	// within my fork of viper, I implemented WatchRemoteConfig function, but
	// github.com/xordataexchange/crypt/backend/etcd/etcd.go also need to modifiy a little bit.
	// use this line instead of c.waitIndex+1 :
	//   resp, err = c.client.Watch(key, 0, false, nil, stop)

	viper.AddRemoteProvider("etcd", "http://192.168.59.103:4001", "/config/host/DST54869.yml")
	viper.SetConfigType("YAML") // because there is no file extension in a stream of bytes

	var x config.Config

	for {
		fmt.Println("watching")
		err := viper.WatchRemoteConfig()
		if err != nil {
			fmt.Println(err)
		}

		err = viper.Marshal(&x)
		if err != nil {
			fmt.Println(err)
		}

		j, _ := json.Marshal(x)

		fmt.Println(string(j))
	}
}

func main() {
	// try_watch()
	// try_vip_get_remote()

	try_viper_watch_remote()

}
