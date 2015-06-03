package config

import (
	b_conf "github.com/oliveagle/hickwall/backends/config"
	"github.com/oliveagle/hickwall/newcore"
)

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
