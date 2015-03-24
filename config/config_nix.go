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
	config_path = append(config_path, fmt.Sprintf("/etc/%s/", APP_NAME))
	config_path = append(config_path, ".")

	viper.SetConfigName("config")
	for _, path := range config_path {
		viper.AddConfigPath(path)
	}
}
