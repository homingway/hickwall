package config

import (
	"fmt"
	// "github.com/spf13/viper"
	"github.com/oliveagle/viper"
	// "log"
	"os"
	"path"
	"path/filepath"
	"reflect"

	log "github.com/oliveagle/seelog"

	// "encoding/json"
	// "io/ioutil"
	"time"
)

const (
	VERSION = "v0.0.1"
)

var (
	config_path        []string
	LOG_DIR            = ""
	SHARED_DIR         = ""
	CONF_FILEPATH      = ""
	CORE_CONF_FILEPATH = ""

	LOG_FILE     = "hickwall.log"
	LOG_FILEPATH = ""

	CoreConf   CoreConfig // only file
	core_viper = viper.New()

	rconf Config // can retrived from file or etcd

	RuntimeConfChan = make(chan *Config, 1)
)

type CoreConfig struct {
	// Log_console_level string   // always 'info'
	Log_level         string `json:"log_level"`
	Log_file_maxsize  int    `json:"log_file_maxsize"`
	Log_file_maxrolls int    `json:"log_file_maxrolls"`

	Etcd_enabled        bool   `json:"etcd_enabled"`
	Etcd_url            string `json:"etcd_url"`
	Etcd_path           string `json:"etcd_path"`
	Etcd_check_interval string `json:"etcd_check_interval"`

	Heart_beat_interval string `json:"heart_beat_interval"`
}

type Config struct {
	Tags map[string]string `json:"tags"`

	Client_metric_enabled  bool   `json:"client_metric_enabled"`
	Client_metric_interval string `json:"client_metric_interval"`

	Transport_flat_metric_key_format string `json:"transport_flat_metric_key_format"`
	Transport_backfill_enabled       bool   `json:"transport_backfill_enabled"`

	Transport_stdout         Transport_stdout `json:"transport_stdout"`
	Transport_file           Transport_file   `json:"transport_file"`
	Transport_graphite_hosts []string         `json:"transport_graphite_hosts"`

	Transport_influxdb []Transport_influxdb `json:"transport_influxdb"`

	Collector_win_pdh []Conf_win_pdh `json:"collector_win_pdh"`
	Collector_win_wmi []Conf_win_wmi `json:"collector_win_wmi"`

	Collector_mysql_query []c_mysql_query `json:"collector_mysql_query"`

	Collector_ping []c_ping `json:"collector_ping"`

	Collector_cmd []Conf_cmd `json:"collector_cmd"`
}

type Conf_win_pdh struct {
	Tags     map[string]string    `json:"tags"`
	Interval string               `json:"interval"`
	Queries  []Conf_win_pdh_query `json:"queries"`
}

type Conf_win_pdh_query struct {
	Query  string            `json:"query"`
	Metric string            `json:"metric"`
	Tags   map[string]string `json:"tags"`
	// Meta   map[string]string		//TODO: Meta
}

type Conf_win_wmi struct {
	Tags     map[string]string    `json:"tags"`
	Interval string               `json:"interval"`
	Queries  []Conf_win_wmi_query `json:"queries"`
}
type Conf_win_wmi_query struct {
	Query   string                      `json:"query"`
	Tags    map[string]string           `json:"tags"`
	Metrics []Conf_win_wmi_query_metric `json:"metrics"`
}
type Conf_win_wmi_query_metric struct {
	//TODO: Meta
	Value_from string            `json:"value_from"`
	Metric     string            `json:"metric"`
	Tags       map[string]string `json:"tags"`
	Meta       map[string]string `json:"meta"`
	Default    interface{}       `json:"default"`
}

type Conf_cmd struct {
	Cmd      []string          `json:"cmd"`
	Interval string            `json:"interval"`
	Tags     map[string]string `json:"tags"`
}

type Transport_file struct {
	Enabled        bool   `json:"enabled"`
	Flush_Interval string `json:"flush_interval"`
	Path           string `json:"path"`

	// TODO: max_size, max_rotation
	Max_size     int `json:"max_size"`
	Max_rotation int `json:"max_rotation"`
}

type Transport_stdout struct {
	Enabled bool `json:"enabled"`
}

type Transport_influxdb struct {
	Version        string `json:"version"`
	Enabled        bool   `json:"enabled"`
	Interval       string `json:"interval"`
	Max_batch_size int    `json:"max_match_size"`

	// Client Config
	// for v0.8.8
	Host string `json:"host"`

	// for v0.9.0
	URL string `json:"url"`

	Username string `json:"username"`
	Password string `json:"password"`
	Database string `json:"database"`

	// Write Config
	RetentionPolicy string `json:"retentionpolicy"`
	FlatTemplate    string `json:"flattemplate"`

	Backfill_enabled              bool   `json:"backfill_enabled"`
	Backfill_interval             string `json:"backfill_interval"`
	Backfill_handsoff             bool   `json:"backfill_handsoff"`
	Backfill_latency_threshold_ms int    `json:"backfill_latency_threshold_ms"`
	Backfill_cool_down            string `json:"backfill_cool_down"`

	// try best to merge small group of points to no more than max_batch_size
	Merge_Requests bool `json:"merge_requests"`
}

type c_mysql_query struct {
	Metric_key string               `json:"metric_key"`
	Tags       [][]string           `json:"tags"`
	Host       string               `json:"host"`
	Port       int                  `json:"port"`
	Username   string               `json:"username"`
	Password   string               `json:"password"`
	Queries    []c_mysql_query_item `json:"queries"`
}

type c_mysql_query_item struct {
	Metric_key string     `json:"metric_key"`
	Tags       [][]string `json:"tags"`
	Database   string     `json:"database"`
	Desc       string     `json:"desc"`
	Query      string     `json:"query"`
	ValuesFrom string     `json:"valuesfrom"`
	Comment    string     `json:"comment"`
}

type c_ping struct {
	Metric_key string     `json:"metric_key"`
	Tags       [][]string `json:"tags"`

	Hosts    []string `json:"hosts"`
	Interval int      `json:"interval"`
}

func (c *Config) setDefaultByKey(key string, val interface{}) (err error) {

	runtime_conf := GetRuntimeConf()

	if !viper.IsSet(key) {
		// fmt.Println("key is not set", key)
		el := reflect.ValueOf(runtime_conf).Elem().FieldByName(key)
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

func loadCoreConfig() {

	initPathes()

	core_viper.SetConfigName("core_config")
	// core_viper.SetConfigFile(CORE_CONF_FILEPATH)
	core_viper.AddConfigPath(SHARED_DIR) // packaged distribution
	core_viper.AddConfigPath(".")        // for hickwall
	core_viper.AddConfigPath("..")       // for hickwall/misc
	core_viper.AddConfigPath("../..")    // for hickwall/misc/try_xxx

	// err := LoadCoreConfig()
	err := core_viper.ReadInConfig()
	if err != nil {
		log.Errorf("No configuration file loaded. core_config.yml :%v", err)
		log.Flush()
		os.Exit(1)
	}

	fmt.Println("core config file used: ", core_viper.ConfigFileUsed())

	err = core_viper.Marshal(&CoreConf)
	if err != nil {
		log.Errorf("Error: unable to parse Core Configuration: %v\n", err)
		log.Flush()
		os.Exit(1)
	}

	fmt.Println("enable_remote_config: ", CoreConf.Etcd_enabled)

	ConfigLogger()
	if err != nil {
		log.Errorf("LoadCoreConfFile failed: %v", err)
		log.Error("SHARED_DIR: ", SHARED_DIR)
		// log.Error("CORE_CONF_FILEPATH: ", CORE_CONF_FILEPATH)
		log.Flush()
		os.Exit(1)
	} else {
		log.Debug("init config, core config loaded")
		log.Debug("LOG_DIR: ", LOG_DIR)
		log.Debug("LOG_FILEPATH: ", LOG_FILEPATH)
	}
}

type RespConfig struct {
	Config *Config
	Err    error
}

func loadRuntimeConfFromFile() <-chan *RespConfig {
	var (
		out           = make(chan *RespConfig, 1)
		runtime_viper = viper.New()
	)
	// runtime_viper.SetConfigFile(config_file)
	runtime_viper.SetConfigName("config")
	runtime_viper.SetConfigType("yaml")
	runtime_viper.AddConfigPath(SHARED_DIR) // packaged distribution
	runtime_viper.AddConfigPath("../..")    // for hickwall/misc/try_xxx
	runtime_viper.AddConfigPath(".")        // for hickwall
	runtime_viper.AddConfigPath("..")       // for hickwall/misc

	log.Info("loadRuntimeConfFromFile")

	go func() {
		log.Info("----------------------------====")
		var runtime_conf Config

		err := runtime_viper.ReadInConfig()

		log.Debug("Config File Used: ", runtime_viper.ConfigFileUsed())
		// fmt.Println("file used: ", runtime_viper.ConfigFileUsed())
		if err != nil {
			fmt.Println("err: ", err)
			out <- &RespConfig{nil, fmt.Errorf("No configuration file loaded. config.yml: %v", err)}
			return
		}

		// Marshal values
		err = runtime_viper.Marshal(&runtime_conf)
		if err != nil {
			fmt.Println("err: ", err)
			out <- &RespConfig{nil, fmt.Errorf("Error: unable to parse Configuration: %v\n", err)}
			return
		}

		out <- &RespConfig{&runtime_conf, nil}
		return
	}()

	return out
}

func WatchRuntimeConfFromEtcd(stop chan bool) <-chan *RespConfig {
	var (
		runtime_viper = viper.New()
		out           = make(chan *RespConfig, 1)
	)
	runtime_viper.SetConfigType("YAML")
	runtime_viper.AddRemoteProvider("etcd", CoreConf.Etcd_url, CoreConf.Etcd_path)

	if stop == nil {
		stop = make(chan bool)
	}

	go func() {
		for {
			var runtime_conf Config

			select {
			case <-stop:
				log.Debugf("stop watching etcd remote config.")
				break
			default:
				log.Debugf("watching etcd remote config: %s, %s", CoreConf.Etcd_url, CoreConf.Etcd_path)
				err := viper.WatchRemoteConfig()
				if err != nil {
					log.Errorf("unable to read remote config: %v", err)
					// return nil, fmt.Errorf("unable to read remote config: %v", err)
					continue
				}

				err = runtime_viper.Marshal(&runtime_conf)
				if err != nil {
					// return nil, fmt.Errorf("unable to marshal to config: %v", err)
					log.Errorf("unable to marshal to config: %v", err)
					continue
				}

				log.Debugf("a new config is comming")
				out <- &RespConfig{&runtime_conf, nil}

				//TODO: make it configurable
				time.Sleep(time.Second * 5)
			}
		}
	}()
	return out
}

func initPathes() {
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	dir, _ = filepath.Split(dir)

	SHARED_DIR, _ = filepath.Abs(path.Join(dir, "shared"))
	// SHARED_DIR = `D:\Users\rhtang\oledev\gocodez\src\github.com\oliveagle\shared`

	LOG_DIR, _ = filepath.Abs(path.Join(SHARED_DIR, "logs"))
	LOG_FILEPATH, _ = filepath.Abs(path.Join(LOG_DIR, LOG_FILE))

	CONF_FILEPATH, _ = filepath.Abs(path.Join(SHARED_DIR, "config.yml"))
	CORE_CONF_FILEPATH, _ = filepath.Abs(path.Join(SHARED_DIR, "core_config.yml"))

	// fmt.Println("dir: ", dir)
	// fmt.Println("SHARED_DIR: ", SHARED_DIR)
	// fmt.Println("CONF_FILEPATH: ", CONF_FILEPATH)

	Mkdir_p_logdir(LOG_DIR)
}

func WatchConfig() <-chan *RespConfig {
	if CoreConf.Etcd_enabled == true {
		return WatchRuntimeConfFromEtcd(nil)
	} else {
		return loadRuntimeConfFromFile()
	}
}

func init() {
	loadCoreConfig()
	log.Info("CoreConfig Loaded")
}

func UpdateRuntimeConf(conf *Config) {
	rconf = *conf
}

func GetRuntimeConf() *Config {
	return &rconf
}
