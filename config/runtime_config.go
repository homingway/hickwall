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

func ReadRuntimeConfig(r io.Reader) (rc RuntimeConfig, err error) {
	vp := viper.New()
	vp.SetConfigType("yaml")
	err = vp.ReadConfig(r)
	if err != nil {
		return rc, err
	}
	err = vp.Marshal(&rc)
	if err != nil {
		return rc, err
	}

	return rc, nil
}

type RuntimeConfig struct {
	Client ClientConfig            `json:"client"`
	Groups []*CollectorConfigGroup `json:"groups"`
}

type ClientConfig struct {
	HeartBeat_Interval string
	Tags               map[string]string
	Hostname           string

	Transport_dummy    *Transport_dummy             `json:"transport_dummy"` // for testing purpose
	Transport_file     *b_conf.Transport_file       `json:"transport_file"`
	Transport_influxdb []*b_conf.Transport_influxdb `json:"transport_influxdb"`
}

type Transport_dummy struct {
	Name      string
	Jamming   newcore.Interval
	Printting bool
	Detail    bool
}
