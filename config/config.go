package config

import (
	"fmt"
	"github.com/spf13/viper"
)

const (
	APPNAME   = "oleservice"
	CONF_NAME = "config"
)

type Config struct {
	Name string
}

func SetDefault() {
	viper.SetConfigName(CONF_NAME)
	// viper.AddConfigPath(fmt.Sprintf("/etc/%s/", APPNAME))

	viper.SetConfigName("config")
	viper.SetConfigType("toml")

	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("No configuration file loaded - using defaults")
	}

	viper.SetDefault("msg", "hello")

}
