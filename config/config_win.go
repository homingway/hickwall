// +build windows

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

	viper.AddConfigPath(fmt.Sprintf("c:\\hickwall\\%s", APP_NAME))
	viper.AddConfigPath(".")

	viper.SetConfigName("config_win")
	viper.SetConfigType("toml")

	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("No configuration file loaded - using defaults")
	}

	viper.SetDefault("port", ":9977")
}
