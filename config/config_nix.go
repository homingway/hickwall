// +build linux darwin

package config

import (
	"fmt"
	"github.com/spf13/viper"
)

const (
	ALLOWED_COLOR_LOG = true
)

func addConfigPath() {
	viper.SetConfigName("config")
	viper.AddConfigPath(fmt.Sprintf("/etc/%s/", APP_NAME))
	viper.AddConfigPath(".")
}
