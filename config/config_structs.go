package config

//import (
////	"fmt"
//
//// "github.com/spf13/viper"
////	"github.com/oliveagle/hickwall/newcore"
////	"github.com/oliveagle/viper"
////	"reflect"
//
//// log "github.com/oliveagle/seelog"
//// "time"
//)

//type CoreConfig struct {
//	// Log_console_level string   // always 'info'
//	Log_level         string `json:"log_level"`
//	Log_file_maxsize  int    `json:"log_file_maxsize"`
//	Log_file_maxrolls int    `json:"log_file_maxrolls"`
//
//	Etcd_enabled bool   `json:"etcd_enabled"`
//	Etcd_url     string `json:"etcd_url"`
//	Etcd_path    string `json:"etcd_path"`
//
//	Heart_beat_interval string `json:"heart_beat_interval"`
//}

//type RuntimeConfig struct {
//}

//type RuntimeConfig struct {
//	Client Conf_client `json:"client"`
//
//	Transport_stdout         Transport_stdout `json:"transport_stdout"`
//	Transport_file           Transport_file   `json:"transport_file"`
//	Transport_graphite_hosts []string         `json:"transport_graphite_hosts"`
//
//	Transport_influxdb []Transport_influxdb `json:"transport_influxdb"`
//
//	Collector_win_sys Conf_win_sys `json:"collector_win_sys"`
//
//	Collector_win_pdh []Conf_win_pdh `json:"collector_win_pdh"`
//	Collector_win_wmi []Conf_win_wmi `json:"collector_win_wmi"`
//
//	Collector_mysql_query []c_mysql_query `json:"collector_mysql_query"`
//
//	Collector_ping []Conf_ping `json:"collector_ping"`
//
//	Collector_cmd []Conf_cmd `json:"collector_cmd"`
//}

//type Conf_client struct {
//	Hostname           string
//	Heartbeat_interval string         `json:"heartbeat_interval"`
//	Metric_enabled     bool           `json:"metric_enabled"`
//	Metric_interval    string         `json:"metric_interval"`
//	Tags               newcore.TagSet `json:"tags"`
//}
//
//type Conf_win_sys struct {
//	Interval string `json:"interval"`
//}
//
//type Conf_win_pdh struct {
//	Tags     map[string]string    `json:"tags"`
//	Interval string               `json:"interval"`
//	Queries  []Conf_win_pdh_query `json:"queries"`
//}
//
//type Conf_win_pdh_query struct {
//	Query  string            `json:"query"`
//	Metric string            `json:"metric"`
//	Tags   map[string]string `json:"tags"`
//	// Meta   map[string]string     //TODO: Meta
//}
//
//type Conf_win_wmi struct {
//	Tags     map[string]string    `json:"tags"`
//	Interval string               `json:"interval"`
//	Queries  []Conf_win_wmi_query `json:"queries"`
//}
//type Conf_win_wmi_query struct {
//	Query   string                      `json:"query"`
//	Tags    map[string]string           `json:"tags"`
//	Metrics []Conf_win_wmi_query_metric `json:"metrics"`
//}
//type Conf_win_wmi_query_metric struct {
//	//TODO: Meta
//	Value_from string            `json:"value_from"`
//	Metric     string            `json:"metric"`
//	Tags       map[string]string `json:"tags"`
//	Meta       map[string]string `json:"meta"`
//	Default    interface{}       `json:"default"`
//}
//
//type Conf_cmd struct {
//	Cmd      []string          `json:"cmd"`
//	Interval string            `json:"interval"`
//	Tags     map[string]string `json:"tags"`
//}
//
//type Transport_file struct {
//	Enabled        bool   `json:"enabled"`
//	Flush_Interval string `json:"flush_interval"`
//	Path           string `json:"path"`
//
//	// TODO: max_size, max_rotation
//	Max_size     int `json:"max_size"`
//	Max_rotation int `json:"max_rotation"`
//}
//
//type Transport_stdout struct {
//	Enabled bool `json:"enabled"`
//}
//
//type Transport_influxdb struct {
//	Version        string `json:"version"`
//	Enabled        bool   `json:"enabled"`
//	Interval       string `json:"interval"`
//	Max_batch_size int    `json:"max_match_size"`
//
//	Max_queue_size int64 `json:"max_queue_size"`
//
//	// Client Config
//	// for v0.8.8
//	Host string `json:"host"`
//
//	// for v0.9.0
//	URL string `json:"url"`
//
//	Username string `json:"username"`
//	Password string `json:"password"`
//	Database string `json:"database"`
//
//	// Write Config
//	RetentionPolicy string `json:"retentionpolicy"`
//	FlatTemplate    string `json:"flattemplate"`
//
//	Backfill_enabled              bool   `json:"backfill_enabled"`
//	Backfill_interval             string `json:"backfill_interval"`
//	Backfill_handsoff             bool   `json:"backfill_handsoff"`
//	Backfill_latency_threshold_ms int    `json:"backfill_latency_threshold_ms"`
//	Backfill_cool_down            string `json:"backfill_cool_down"`
//
//	// try best to merge small group of points to no more than max_batch_size
//	Merge_Requests bool `json:"merge_requests"`
//}
