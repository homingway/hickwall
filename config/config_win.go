// +build windows

package config

import (
	// "fmt"
	"github.com/spf13/viper"
)

const (
	ALLOWED_COLOR_LOG = false
)

func addConfigPath() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AddConfigPath("..")
	viper.AddConfigPath("../..")
}
