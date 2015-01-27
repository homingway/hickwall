// +build linux darwin

package config

import (
	"fmt"
	"github.com/spf13/viper"
)

func addConfigPath() {
	viper.SetConfigName("config")
	viper.AddConfigPath(fmt.Sprintf("/etc/%s/", APP_NAME))
	viper.AddConfigPath(".")
}
