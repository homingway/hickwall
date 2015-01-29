package config

import (
	"fmt"
	// "github.com/kr/pretty"
	"github.com/spf13/viper"
	"os"
	"reflect"
)

const (
	VERSION   = "v0.0.1"
	CONF_NAME = "config"
	APP_NAME  = "hickwall"
	APP_DESC  = "monitoring system"
)

// Design Principle: Flat is better than nasted!
type Config struct {
	Port                string
	Logfile             string
	Log_colored_console bool
	Log_console_level   string
	Log_console_format  string
	Log_file_level      string
	Log_file_filepath   string
	Log_file_format     string
	Log_file_maxsize    int
	Log_file_maxrolls   int
}

var Conf Config

func (c *Config) setDefaultByKey(key string, val interface{}) (err error) {

	if !viper.IsSet(key) {
		// fmt.Println("key is not set", key)
		el := reflect.ValueOf(&Conf).Elem().FieldByName(key)
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
			el.SetInt(int64(val.(int)))
		default:
			err = fmt.Errorf("unexpected type %T, key: %s, value: %v", kind, key, val)
		}
		// } else {
		// 	fmt.Println("key set: ", key, viper.Get(key))
	}
	return
}

func init() {
	fmt.Println("Initializing Configuration")

	viper.SetConfigType("toml")

	// read config file
	addConfigPath()
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("No configuration file loaded - using defaults")
	}

	// Marshal values
	err = viper.Marshal(&Conf)
	if err != nil {
		fmt.Printf("Error: unable to parse Configuration: %v", err)
		os.Exit(1)
	}

	// fmt.Println("-------- after marshal --------------")
	// pretty.Println(&Conf)

	// place all setDefault here -----------------
	// First we have to find out which config item is not been set in config.toml
	// then we only set default values to these missing items.
	//TODO: remove port :9977
	Conf.setDefaultByKey("port", ":9977")
	Conf.setDefaultByKey("Logfile", "/var/log/hickwall/hickwall.log")

	// fmt.Println("-------- after setdefault --------------")
	// pretty.Println(Conf)
	ConfigLogger()
	// logger, err := log.LoggerFromConfigAsFile("seelog.xml")
	// if err != nil {
	// 	fmt.Println("Error: cannot load log config file: ", err)
	// 	return
	// }
	// log.ReplaceLogger(logger)

}
