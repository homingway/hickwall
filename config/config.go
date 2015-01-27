package config

import (
	"fmt"
	"github.com/kr/pretty"
	"github.com/spf13/viper"
	"os"
	"reflect"
	"strings"
)

const (
	VERSION   = "v0.0.1"
	CONF_NAME = "config"
	APP_NAME  = "hickwall"
	APP_DESC  = "monitoring system"
)

type Config struct {
	Port       string
	Logfile    string
	Consolelog ConsoleLog
	Filelog    FileLog
}

type ConsoleLog struct {
	Level  string
	Format string
}

type FileLog struct {
	Level    string
	Logfile  string
	Format   string
	Rotate   bool
	Maxsize  string
	Maxlines string
	Daily    bool
}

var Conf Config

func (c *Config) setDefault(key string, val interface{}) (err error) {
	el := reflect.ValueOf(&Conf).Elem()
	key_slice := strings.Split(key, ".")
	for idx, k := range key_slice {
		// fmt.Println(idx, k)
		el = el.FieldByName(k)

		if idx == len(key_slice)-1 {
			// fmt.Println(idx, k, el, el.Type(), el.Kind(), el.Type())
			kind := el.Kind()
			switch kind {
			case reflect.Bool:
				el.SetBool(val.(bool))
			case reflect.Float32:
				el.SetFloat(val.(float64))
			case reflect.Float64:
				el.SetFloat(val.(float64))
			case reflect.String:
				el.SetString(val.(string))
			case reflect.Int:
				el.SetInt(val.(int64))
			default:
				err = fmt.Errorf("unexpected type %T, key: %s, value: %v", kind, key, val)
			}
		}
	}
	return
}

func init() {
	viper.SetConfigType("toml")

	// read config file
	addConfigPath()
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("No configuration file loaded - using defaults")
	}

	//TODO: remove port :9977
	viper.SetDefault("port", ":9977")

	//
	err = viper.Marshal(&Conf)
	if err != nil {
		fmt.Printf("Error: unable to parse Configuration: %v", err)
		os.Exit(1)
	}

	pretty.Println(Conf)

	parseLoggerConfig()
}
