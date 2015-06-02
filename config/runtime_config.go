package config

import (
	//	"fmt"
	//	"github.com/oliveagle/viper"
	//	"reflect"
	//	"time"
	b_conf "github.com/oliveagle/hickwall/backends/config"
	"github.com/oliveagle/hickwall/newcore"
	"github.com/oliveagle/viper"
	"io"
)

// public variables
var (
	RuntimeConfChan = make(chan *RuntimeConfig, 1)
)

// private vairables
var (
	rconf *RuntimeConfig
)

func GetRuntimeConfig() *RuntimeConfig {
	return rconf
}

func ReadRuntimeConfig(r io.Reader) (*RuntimeConfig, error) {
	var rc RuntimeConfig

	vp := viper.New()
	vp.SetConfigType("yaml")
	err := vp.ReadConfig(r)
	if err != nil {
		return nil, err
	}
	err = vp.Marshal(&rc)
	if err != nil {
		return nil, err
	}

	return &rc, nil
}

type RuntimeConfig struct {
	Client ClientConfig            `json:"client"`
	Groups []*CollectorConfigGroup `json:"groups"`
}

type ClientConfig struct {
	HeartBeat_Interval string
	Tags               map[string]string
	Hostname           string

	Transport_dummy *Transport_dummy // for testing purpose

	Transport_file *b_conf.Transport_file

	//	Transport_influxdb []Transport_influxdb `json:"transport_influxdb"`
	//	Transport_influxdb string
}

type Transport_dummy struct {
	Name      string
	Jamming   newcore.Interval
	Printting bool
	Detail    bool
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
