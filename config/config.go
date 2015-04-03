package config

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
	"reflect"
)

const (
	VERSION = "v0.0.1"
)

var config_path []string

type Config struct {
	Tags map[string]string

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

	Client_metric_enabled  bool
	Client_metric_interval string

	Transport_flat_metric_key_format string
	Transport_backfill_enabled       bool
	Transport_graphite_hosts         []string

	Transport_influxdb []Transport_influxdb

	Collector_win_pdh []Conf_win_pdh
	Collector_win_wmi []Conf_win_wmi

	Collector_mysql_query []c_mysql_query

	Collector_ping []c_ping

	Collector_cmd []Conf_cmd
}

type Conf_win_pdh struct {
	Tags     map[string]string
	Interval int
	Queries  []Conf_win_pdh_query
}

type Conf_win_pdh_query struct {
	Query  string
	Metric string
	Tags   map[string]string
	// Meta   map[string]string		//TODO: Meta
}

type Conf_win_wmi struct {
	Tags     map[string]string
	Interval string
	Queries  []Conf_win_wmi_query
}
type Conf_win_wmi_query struct {
	Query   string
	Tags    map[string]string
	Metrics []Conf_win_wmi_query_metric
}
type Conf_win_wmi_query_metric struct {
	Value_from string
	Metric     string
	Tags       map[string]string
	Meta       map[string]string //TODO: Meta
	Default    interface{}
}

type Conf_cmd struct {
	Cmd      []string
	Interval string
	Tags     map[string]string
}

type Transport_influxdb struct {
	Version        string
	Enabled        bool
	Interval       string
	Max_batch_size int

	// Client Config
	Host     string // for v0.8.8
	URL      string // for v0.9.0
	Username string
	Password string
	Database string

	// Write Config
	RetentionPolicy string
	FlatTemplate    string

	Backfill_enabled              bool
	Backfill_interval             string
	Backfill_handsoff             bool
	Backfill_latency_threshold_ms int
	Backfill_cool_down_s          int

	Merge_Requests bool // try best to merge small group of points to no more than max_batch_size
}

type c_mysql_query struct {
	Metric_key string
	Tags       [][]string
	Host       string
	Port       int
	Username   string
	Password   string
	Queries    []c_mysql_query_item
}

type c_mysql_query_item struct {
	Metric_key string
	Tags       [][]string
	Database   string
	Desc       string
	Query      string
	ValuesFrom string
	Comment    string
}

type c_ping struct {
	Metric_key string
	Tags       [][]string

	Hosts    []string
	Interval int
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

func LoadConfigFile() error {
	err := viper.ReadInConfig()
	if err != nil {
		// log.Println("No configuration file loaded - using defaults")
		return fmt.Errorf("No configuration file loaded. config.yml")
	}

	// Marshal values
	err = viper.Marshal(&Conf)
	if err != nil {
		log.Fatalln("Error: unable to parse Configuration: %v", err)
	}

	Conf.setDefaultByKey("port", ":9977")
	Conf.setDefaultByKey("Logfile", "/var/log/hickwall/hickwall.log")

	// // fmt.Println("-------- after setdefault --------------")
	ConfigLogger()

	return nil
}

func init() {
	// viper.SetConfigType("toml")
	//viper.SetConfigType("yml")

	// read config file
	addConfigPath()

	LoadConfigFile()

	// err := viper.ReadInConfig()
	// if err != nil {
	// 	// log.Println("No configuration file loaded - using defaults")
	// 	return fmt.Errorf("No configuration file loaded. config.yml")
	// }

	// Marshal values
	// err = viper.Marshal(&Conf)
	// if err != nil {
	// 	log.Fatalln("Error: unable to parse Configuration: %v", err)
	// }

	// place all setDefault here -----------------
	// First we have to find out which config item is not been set in config.toml
	// then we only set default values to these missing items.

	//TODO: remove port :9977
	// Conf.setDefaultByKey("port", ":9977")
	// Conf.setDefaultByKey("Logfile", "/var/log/hickwall/hickwall.log")

	// // fmt.Println("-------- after setdefault --------------")
	// ConfigLogger()
}
