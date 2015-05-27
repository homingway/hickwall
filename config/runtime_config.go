package config

import (
//	"fmt"
//	log "github.com/oliveagle/seelog"
//	"github.com/oliveagle/viper"
//	"reflect"
//	"time"
)

// public variables
var (
	RuntimeConfChan = make(chan *RuntimeConfig, 1)
)

// private vairables
var (
//	rconf *RuntimeConfig
)

type RuntimeConfig struct {
	Client ClientConfig
	Groups []*CollectorConfigGroup
}

type ClientConfig struct {

	//	Transport_influxdb []Transport_influxdb `json:"transport_influxdb"`
	Transport_influxdb string
	Transport_stdout   string
	Transport_file     string
}

type CollectorConfigGroup struct {
	Collector_ping string
	//	Collector_win_sys Conf_win_sys `json:"collector_win_sys"`
	//
	//	Collector_win_pdh []Conf_win_pdh `json:"collector_win_pdh"`
	//	Collector_win_wmi []Conf_win_wmi `json:"collector_win_wmi"`

	//	Collector_mysql_query []c_mysql_query `json:"collector_mysql_query"`
	//
	//	Collector_ping []Conf_ping `json:"collector_ping"`
	//	Collector_cmd []Conf_cmd `json:"collector_cmd"`
}

//func (c *RuntimeConfig) setDefaultByKey(key string, val interface{}) (err error) {
//
//	runtime_conf := GetRuntimeConf()
//
//	if !viper.IsSet(key) {
//		// fmt.Println("key is not set", key)
//		el := reflect.ValueOf(runtime_conf).Elem().FieldByName(key)
//		kind := el.Kind()
//		switch kind {
//		case reflect.Bool:
//			el.SetBool(val.(bool))
//		case reflect.Float32:
//			el.SetFloat(val.(float64))
//		case reflect.Float64:
//			el.SetFloat(val.(float64))
//		case reflect.String:
//			el.SetString(val.(string))
//		case reflect.Int:
//			el.SetInt(int64(val.(int)))
//		default:
//			err = fmt.Errorf("unexpected type %T, key: %s, value: %v", kind, key, val)
//		}
//		// } else {
//		//  fmt.Println("key set: ", key, viper.Get(key))
//	}
//	return
//}

//func UpdateRuntimeConf(conf *RuntimeConfig) {
//	rconf = conf
//}
//
//func GetRuntimeConf() *RuntimeConfig {
//	return rconf
//}
