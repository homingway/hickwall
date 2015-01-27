// +build linux darwin

package config

import (
	"fmt"
	"github.com/spf13/viper"
)

const (
	VERSION   = "v0.0.1"
	CONF_NAME = "config"
	APP_NAME  = "hickwall"
	APP_DESC  = "monitoring system"
)

type Config struct {
	Name string
}

func LoadConfig() {
	viper.SetConfigName(CONF_NAME)

	viper.AddConfigPath(fmt.Sprintf("/etc/%s/", APP_NAME))
	viper.AddConfigPath(".")

	viper.SetConfigName("config")
	viper.SetConfigType("toml")

	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("No configuration file loaded - using defaults")
	}

	viper.SetDefault("port", ":9977")
}
