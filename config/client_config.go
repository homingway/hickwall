package config

import (
	b_conf "github.com/oliveagle/hickwall/backends/config"
	c_conf "github.com/oliveagle/hickwall/collectors/config"
	"github.com/oliveagle/hickwall/newcore"
)

type ClientConfig struct {
	HeartBeat_Interval string            `json:"heartbeat_interval"`
	Tags               map[string]string `json:"tags"`
	Metric_Enabled     bool              `json:"metric_enabled"`
	Metric_Interval    string            `json:"metric_interval"`

	Transport_dummy         *Transport_dummy                   `json:"transport_dummy"` // for testing purpose
	Transport_file          *b_conf.Transport_file             `json:"transport_file"`
	Transport_influxdb      []*b_conf.Transport_influxdb       `json:"transport_influxdb"`
	Transport_kafka         *b_conf.Transport_kafka            `json:"transport_kafka"`
	Subscribe_kafka         []*c_conf.Config_KafkaSubscription `json:"subscribe_kafka"` // subscribe kafka
	Transport_elasticsearch *b_conf.Transport_elasticsearch    `json:"transport_elasticsearch"`
}

type Transport_dummy struct {
	Name      string
	Jamming   newcore.Interval
	Printting bool
	Detail    bool
}
