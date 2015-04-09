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
	config_path = append(config_path, "/etc/hickwall/")
	config_path = append(config_path, "/opt/hickwall/current/")
	config_path = append(config_path, ".")
	config_path = append(config_path, "..")
	config_path = append(config_path, "../..")

	viper.SetConfigName("config")
	for _, path := range config_path {
		viper.AddConfigPath(path)
	}
}
