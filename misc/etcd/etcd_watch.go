package main

import (
	"encoding/json"
	"fmt"
	"github.com/oliveagle/hickwall/config"
	"github.com/spf13/viper"
)

func main() {
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
